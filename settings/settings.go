package settings

import "github.com/olebedev/config"

const (
	Domain string = "https://dcysmibx9jx7t.cloudfront.net"
	Folder string = "./files/"
	Ext string = ".jpg"
	Cfg_file string = "./config.yaml"
	S3Folder = "/uploads/"
)

func GetConf() (*config.Config, error ){
	conf, error := config.ParseYamlFile(Cfg_file)
	return conf, error
}

func GetHost() (string){
	conf, _ := GetConf()
	host, _ := conf.String("database.mongodb.host")

	return host
}
