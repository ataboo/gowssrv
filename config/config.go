package config

import (
"log"

"github.com/BurntSushi/toml"
)

var Config tomlConfig

func init() {
	Config.Read()
}

// Represents database server and credentials
type tomlConfig struct {
	Mongo mongo
	Api api
}

type mongo struct {
	DB   string
	Server string
}

type api struct {
	JwtSecret string
}

// Read and parse the configuration file
func (c *tomlConfig) Read() {
	if _, err := toml.DecodeFile("config.toml", &c); err != nil {
		log.Println("Error loading config.toml, check config.toml.example")
		log.Fatal(err)
	}
}
