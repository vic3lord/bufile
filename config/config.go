package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Modules []string `json:"modules"`
}

func Parse(path string) (Config, error) {
	var cfg Config
	f, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("reading config file: %w", err)
	}

	err = json.Unmarshal(f, &cfg)
	if err != nil {
		return cfg, fmt.Errorf("parsing config file: %w", err)
	}
	return cfg, nil
}