package config



import (
	"encoding/json"
	"os"
)

type Config struct {
	Login         int    `json:"Login"`
	Password      string `json:"Password"`
	Server        string `json:"Server"`
	DefaultSymbol string `json:"DefaultSymbol"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
