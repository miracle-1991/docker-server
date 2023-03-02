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

type Config struct {
	AirLabConf AirLabs `yaml:"airlabs"`
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

func init() {
	initFromLocal()
}

func GetAirLabConfig() *AirLabs {
	return &config.AirLabConf
}
