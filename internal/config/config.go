package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
)

type Config struct {
	Telegram TGConfig `yaml:"telegram"`
}

type TGConfig struct {
	Token string `yaml:"token"`
}

func NewConfig() *Config {
	data, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("failed to unmarshal config file: %v", err)
	}

	return &config
}
