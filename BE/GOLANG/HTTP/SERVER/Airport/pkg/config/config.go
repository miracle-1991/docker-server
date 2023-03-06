package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"

	"log"
	"os"
)

type AirLabs struct {
	ApiKey string `yaml:"api_key"`
}
type DataBase struct {
	Type     string `yaml:"type"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Name     string `yaml:"name"`
}

type Config struct {
	AirLabConf   AirLabs  `yaml:"airlabs"`
	DataBaseConf DataBase `yaml:"database"`
}

var config Config

func initFromLocal() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	configfile := pwd + "/" + "config.yaml"
	file, err := ioutil.ReadFile(configfile)
	if err != nil {
		log.Fatal(err)
	}
	err = yaml.Unmarshal(file, &config)
	if err != nil {
		log.Fatal(err)
	}
}

func changeConfByEnv() {
	dbtype := os.Getenv("DB_TYPE")
	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")
	dbuser := os.Getenv("DB_USERNAME")
	dbpass := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	if dbtype != "" && dbhost != "" && dbport != "" && dbuser != "" && dbpass != "" && dbname != "" {
		config.DataBaseConf.Type = dbtype
		config.DataBaseConf.Host = dbhost + ":" + dbport
		config.DataBaseConf.User = dbuser
		config.DataBaseConf.Password = dbpass
		config.DataBaseConf.Name = dbname
	}
}

func init() {
	initFromLocal()
	changeConfByEnv()
}

func GetAirLabConfig() *AirLabs {
	return &config.AirLabConf
}

func GetDataBaseConfig() *DataBase {
	return &config.DataBaseConf
}
