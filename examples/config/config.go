package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Config struct {
	Login         int    `json:"Login"`
	Password      string `json:"Password"`
	Server        string `json:"Server"`
	DefaultSymbol string `json:"DefaultSymbol"`
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("config.json load error: %w", err)
	}
	defer file.Close()

	log.Printf("loading config from: %s", filename)

	decoder := json.NewDecoder(file)
	config := &Config{}
	if err := decoder.Decode(config); err != nil {
		return nil, fmt.Errorf("config.json decode error: %w", err)
	}

	log.Printf("loaded config: login=%d server=%q password=%q", config.Login, config.Server, config.Password)

	return config, nil
}
