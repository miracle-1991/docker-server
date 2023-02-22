package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
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

func Init(path string) error {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		return err
	}
	return nil
}

func GetConfig() *Config {
	return &config
}
