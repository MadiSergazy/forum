package config

import (
	"bytes"
	"encoding/json"
	"os"
)

type Config struct {
	Database struct {
		Driver string `json:"driver"`
		DSN    string `json:"dsn"`
	} `json:"database"`
	Server struct {
		Port string `json:"port"`
	} `json:"server"`
}

func NewConfig(path string) (*Config, error) {
	var cfg *Config
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
