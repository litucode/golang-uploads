package imagemanager

import (
	"gopkg.in/mgo.v2"
	"log"
)

type Image struct {
	Url	string `json:"url"`
	Adult 	bool	`json:"adult"`
	Etag	string	`json:"etag"`
}

func SaveImage(url string, adult bool, etag string, host string) (*Image){
	session_mgo, err := mgo.Dial(host)
	if err != nil {
		panic(err)
	}
	defer session_mgo.Close()

	c := session_mgo.DB("test").C("image")

	image := &Image{url, adult, etag}
	err = c.Insert(image)
	if err != nil {
		log.Fatal(err)
	}

	return image
}