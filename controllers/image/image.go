package controllers

import (
	"log"
	"encoding/json"

	"io/ioutil"
	"os"
	"bytes"
	"net/http"

	im "gitlab.com/viajeros-modernos/vm-uploads-media/imagemanager"

	nude "github.com/koyachi/go-nude"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/olebedev/config"
	"strings"
)

type Image struct {
	Url	string `json:"url"`
	Adult 	bool	`json:"adult"`
	Etag	string	`json:"etag"`
}

func PostImagesHandler(w http.ResponseWriter, r *http.Request) {

	// Config
	cfg, err := config.ParseYamlFile("./config.yaml")
	host, err := cfg.String("development.database.host")
	if err != nil {
		log.Println(err)
	} else {
		// Aws Credentials
		aws_access_key_id, err := cfg.String("services.AWSAccessKeyId")
		aws_secret_access_key, err := cfg.String("services.AWSSecretKey")
		aws_region, err := cfg.String("services.AWSRegion")
		aws_bucket, err := cfg.String("services.AWSBucket")
		token := ""

		// AWS cred
		creds := credentials.NewStaticCredentials(aws_access_key_id, aws_secret_access_key, token)
		_, err_creds := creds.Get()
		if err_creds != nil {
			log.Printf("bad credentials: %s", err)
		}
		cfg_aws := aws.NewConfig().WithRegion(aws_region).WithCredentials(creds)
		svc := s3.New(session.New(), cfg_aws)

		/*///////////////////////////////////////////*/
		// Temp file upload
		fileForm, handle, err := r.FormFile("myFile")
		if err != nil {
			log.Printf("Error al cargar el archivo %v", err)
			return
		}
		defer fileForm.Close()

		data, err := ioutil.ReadAll(fileForm)
		if err != nil {
			log.Printf("Error al leer el archivo %v", err)
			return
		}

		err = ioutil.WriteFile("./files/" + handle.Filename, data, 777)
		if err != nil {
			log.Printf("Error al escribir el archivo %v", err)
		}
		log.Printf("Cargado exitosamente")

		/*///////////////////////////////////////////*/
		// open
		file, err := os.Open("files/" + handle.Filename)
		if err != nil {
			log.Printf("err opening file: %s", err)
		}
		defer file.Close()
		fileInfo, _ := file.Stat()
		size := fileInfo.Size()
		buffer := make([]byte, size) // read file content to buffer

		//read
		file.Read(buffer)
		fileBytes := bytes.NewReader(buffer)
		fileType := http.DetectContentType(buffer)

		userID := "asdasd"
		filep := strings.Split(file.Name(), "/")
		path := "/media/" + userID + "/" + filep[1]

		params := &s3.PutObjectInput{
			Bucket: aws.String(aws_bucket),
			Key: aws.String(path),
			Body: fileBytes,
			ContentLength: aws.Int64(size),
			ContentType: aws.String(fileType),
		}

		resp, err := svc.PutObject(params)
		if err != nil {
			log.Printf("bad response: %s", err)
			w.WriteHeader(http.StatusNotFound)
		} else {
			// Detect nude in image
			isNude, err := nude.IsNude("./files/" + handle.Filename)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("isNude = %v\n", isNude)
			log.Printf("response %s", awsutil.StringValue(resp))

			image_return := im.SaveImage(path, isNude, awsutil.StringValue(resp), host)

			/* /////////////////////////////////////////// */

			j, err := json.Marshal(&image_return)
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusOK)
			w.Write(j)
		}
	}
}

