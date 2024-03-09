package config

import (
	"github.com/BurntSushi/toml"
	"log"
)

type Config struct {
	Server Server
	Mongo  MongoConfig
	Client Client
}

type Server struct {
	Port int
}

type MongoConfig struct {
	ConnectionString string
	DbName           string
	Collections      Collections
}

type Collections struct {
	Posts    string
	Comments string
	Users    string
}

type Client struct {
	Username string
	Password string
}

const configFileName = "config.toml"

func GetConfig() *Config {

	var c Config
	if _, err := toml.DecodeFile(configFileName, &c); err != nil {
		log.Fatal(err)
	}

	return &c
}
