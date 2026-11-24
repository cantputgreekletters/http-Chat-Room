package main

import (
	server "ChatRoom/incl"
	"encoding/json"
	"os"
)

type Config struct {
	URL string `json:"url"`
}

func LoadConfig() Config {
	data, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		panic(err)
	}

	return config
}

func main() {
	var conf Config = LoadConfig()
	server.Setup(conf.URL)
}
