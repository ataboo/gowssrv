package config

import (
"log"

"github.com/BurntSushi/toml"
)

var Config ConfigModel

func init() {
	Config.Read()
}

// Represents database server and credentials
type ConfigModel struct {
	MongoServer   string
	MongoDb string
}

// Read and parse the configuration file
func (c *ConfigModel) Read() {
	if _, err := toml.DecodeFile("config.toml", &c); err != nil {
		log.Println("Error loading config.toml, check config.toml.example")
		log.Fatal(err)
	}
}
