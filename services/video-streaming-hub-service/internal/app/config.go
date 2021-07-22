package app

import (
	"fmt"
	"time"
)

type Config struct {
	DatabaseAddr     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string

	ApiAddr        string
	StreamInterval time.Duration
}

func (c Config) Validate() error {
	if c.ApiAddr == "" {
		return fmt.Errorf("field 'ApiAddr' cannot be empty")
	}

	if c.DatabaseAddr == "" {
		return fmt.Errorf("field 'DatabaseAddr' cannot be empty")
	}

	if c.DatabaseName == "" {
		return fmt.Errorf("field 'DatabaseName' cannot be empty")
	}

	return nil
}
