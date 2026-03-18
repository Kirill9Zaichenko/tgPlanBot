package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App       AppConfig       `yaml:"app"`
	Telegram  TelegramConfig  `yaml:"telegram"`
	HTTP      HTTPConfig      `yaml:"http"`
	Database  DatabaseConfig  `yaml:"database"`
	Dragonfly DragonflyConfig `yaml:"dragonfly"`
}

type AppConfig struct {
	Env      string `yaml:"env"`
	LogLevel string `yaml:"log_level"`
}

type TelegramConfig struct {
	Token string `yaml:"token"`
}

type HTTPConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type DatabaseConfig struct {
	SQLitePath string `yaml:"sqlite_path"`
}

type DragonflyConfig struct {
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

func NewConfig() *Config {
	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatalf("failed to read config file: %v", err)
	}

	var cfg Config

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("failed to parse config file: %v", err)
	}

	validate(&cfg)

	return &cfg
}

func validate(cfg *Config) {
	if cfg.Telegram.Token == "" {
		log.Fatal("telegram.token is empty")
	}

	if cfg.Database.SQLitePath == "" {
		log.Fatal("database.sqlite_path is empty")
	}
}
