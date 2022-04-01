package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Listen string `json:"listen"`
}

func (c *Config) Validate() error {
	if c.Listen == "" {
		return fmt.Errorf("c.Listen can't be empty")
	}
	return nil
}

func LoadConfigFromPath(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var c Config
	if err := json.NewDecoder(f).Decode(&c); err != nil {
		return nil, err
	}

	if err := c.Validate(); err != nil {
		return nil, err
	}

	return &c, err
}
