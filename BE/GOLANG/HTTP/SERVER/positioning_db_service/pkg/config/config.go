package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
)

// Presto的配置，连接hive数据库时要用
type PrestoConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
	CatLog   string `yaml:"catlog"`
	Schema   string `yaml:"schema"`
}

// googleapi的配置，渲染html的时候要用
type GoogleApiConfig struct {
	Key          string `yaml:"key"`
	NearbyRadius uint   `yaml:"nearbyRadius"`
}

type Config struct {
	Presto    PrestoConfig    `yaml:"presto"`
	GoogleApi GoogleApiConfig `yaml:"googleapi"`
}

var config Config

func initFromEnv() {
	presto_host := os.Getenv("PRESTO_HOST")
	presto_port := os.Getenv("PRESTO_PORT")
	presto_user := os.Getenv("PRESTO_USER")
	presto_pwd := os.Getenv("PRESTO_PWD")
	presto_catlog := os.Getenv("PRESTO_CATLOG")
	presto_schema := os.Getenv("PRESTO_SCHEMA")
	if presto_port == "" || presto_host == "" || presto_user == "" || presto_pwd == "" || presto_catlog == "" || presto_schema == "" {
		return
	}
	log.Printf("config presto from env: host:%v, port:%v, user:%v, pwd:%v, catlog:%v, schema:%v",
		presto_host, presto_port, presto_user, presto_pwd, presto_catlog, presto_schema)
	config.Presto.Host = presto_host
	config.Presto.Port = presto_port
	config.Presto.UserName = presto_user
	config.Presto.Password = presto_pwd
	config.Presto.CatLog = presto_catlog
	config.Presto.Schema = presto_schema
}

func Init(path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return err
	}
	initFromEnv()
	return nil
}

func GetConfig() *Config {
	return &config
}
