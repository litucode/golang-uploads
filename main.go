package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	nude "github.com/koyachi/go-nude"
	"github.com/olebedev/config"
	mgo "gopkg.in/mgo.v2"
)

type Image struct {
	Url   string
	Adult bool
}

func saveImage(url string, adult bool, host string) {
	session, err := mgo.Dial(host)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("test").C("image")
	err = c.Insert(&Image{url, adult})
	if err != nil {
		log.Fatal(err)
	}

	// result := Image{}
	// err = c.Find(bson.M{"name": "Ale"}).One(&result)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println("Image:", result.Url)
}

func PostImagesHandler(w http.ResponseWriter, r *http.Request) {
	cfg, err := config.ParseYamlFile("./config.yaml")
	host, err := cfg.String("development.database.host")
	if err != nil {
		log.Println(err)
	} else {
		session, err := mgo.Dial("localhost:27017")
		if err != nil {
			panic(err)
		}
		defer session.Close()

		file, handle, err := r.FormFile("myFile")
		if err != nil {
			log.Printf("Error al cargar el archivo %v", err)
			return
		}
		defer file.Close()

		data, err := ioutil.ReadAll(file)
		if err != nil {
			log.Printf("Error al leer el archivo %v", err)
			return
		}

		err = ioutil.WriteFile("./files/"+handle.Filename, data, 777)
		if err != nil {
			log.Printf("Error al escribir el archivo %v", err)
		} else {
			// Detect nude in image
			isNude, err := nude.IsNude("./files/" + handle.Filename)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("isNude = %v\n", isNude)
			saveImage("./files/"+handle.Filename, isNude, host)
		}

		log.Printf("Cargado exitosamente")
		w.WriteHeader(http.StatusOK)
	}
}

func main() {

	cfg, err := config.ParseYamlFile("./config.yaml")

	if err != nil {
		log.Println(err)
	} else {
		port, err := cfg.String("development.server.port")

		if err != nil {
			log.Println(err)
		} else {
			r := mux.NewRouter().StrictSlash(false)

			r.HandleFunc("/images", PostImagesHandler).Methods("POST")

			server := &http.Server{
				Addr:           port,
				Handler:        r,
				ReadTimeout:    10 * time.Second,
				WriteTimeout:   10 * time.Second,
				MaxHeaderBytes: 1 << 20,
			}

			log.Println("Listin in http://localhost" + port)

			log.Fatal(server.ListenAndServe())
		}
	}
}
