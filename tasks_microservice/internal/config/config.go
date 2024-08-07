package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string     `yaml:"env" env-default:"local"`
	Postgres   Postgres   `yaml:"postgres"`
	RabbitMQ RabbitMQ `yaml:"http_server"`
}

type Postgres struct {
	Host     string `yaml:"host" env-default:"localhost"`
	Port     int    `yaml:"port" env-default:"5432"`
	User     string `yaml:"user" env-default:"postgres"`
	Password string `yaml:"password" env-default:"root"`
	DBName   string `yaml:"dbname" env-default:""`
	SSLMode  string `yaml:"sslmode" env-default:"disable"`
}

type RabbitMQ struct {
	Host          string `yaml:"host" env-default:"localhost"`
	RabbitMQPort  int    `yaml:"rabbitMQport" env-default:"5672"`
	User          string `yaml:"user" env-default:"guest"`
	Password      string `yaml:"password" env-default:"guest"`
}


func MustLoad() *Config {
	configPath := "../../config/local.yaml"
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _,err := os.Stat(configPath); os.IsNotExist(err){
		log.Fatalf("config file does not exist: %s", configPath)
	}
	var cfg Config

	if err:= cleanenv.ReadConfig(configPath, &cfg); err != nil{
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
