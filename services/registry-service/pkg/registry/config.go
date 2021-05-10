package registry

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DatabaseBackend  string `json:"database_backend"`
	DatabaseUser     string `json:"database_user"`
	DatabasePassword string `json:"database_password"`
	DatabaseName     string `json:"database_name"`
}

func (c Config) Validate() error {
	if c.DatabaseBackend == "" {
		return fmt.Errorf("c.DatabaseBackend can't be empty")
	}
	if c.DatabaseName == "" {
		return fmt.Errorf("c.DatabaseName can't be empty")
	}
	if c.DatabaseUser == "" {
		return fmt.Errorf("c.DatabaseUser can't be empty")
	}
	if c.DatabasePassword == "" {
		return fmt.Errorf("c.DatabasePassword can't be empty")
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
