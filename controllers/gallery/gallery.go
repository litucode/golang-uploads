package galleryhandler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"log"

	"github.com/olebedev/config"
	"time"
	"strconv"
	"gitlab.com/viajeros-modernos/vm-uploads-media/tools"
)

func PostGalleriesHandler(w http.ResponseWriter, r *http.Request) {

	// Config
	_, err := config.ParseYamlFile("./config.yaml")
	if err != nil {
		log.Println(err)
	} else {

		//////// File upload
		err_parse := r.ParseMultipartForm(200000) // grab the multipart form
		if err_parse != nil {
			log.Printf("Error al parsear", err)
		}

		userID := r.FormValue("user_id")
		formdata := r.MultipartForm // ok, no problem so far, read the Form data

		//get the *fileheaders
		files := formdata.File["files"] // grab the filenames

		for i, _ := range files { // loop through the files one by one

			// Open file[i]
			file, err := files[i].Open()
			if err != nil {
				log.Printf("Error open file %s", err)
			}
			defer file.Close()

			// Write file[i]
			file_name := tools.SetName(files[i].Filename, userID)
			ext := ".jpg"
			out, err := os.Create("./files/" + file_name + ext)

			if err != nil {
				log.Printf("Error write file %s", err)
			}
			defer out.Close()

			_, err = io.Copy(out, file) // file not files[i] !

			if err != nil {
				fmt.Fprintln(w, err)
				return
			}
			fmt.Fprintln(w, file_name)
			log.Println(time.Now())
			log.Println(strconv.FormatInt(time.Now().Unix(), 10))
		}
	}
}
