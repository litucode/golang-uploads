package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/olebedev/config"

	ih "gitlab.com/viajeros-modernos/vm-uploads-media/controllers/image"
)

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

			r.HandleFunc("/images", ih.PostImagesHandler).Methods("POST")
			/* r.HandleFunc("/video", VideoImagesHandler).Methods("POST")*/

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
