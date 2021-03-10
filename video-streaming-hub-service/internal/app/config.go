package app

import "fmt"

type Config struct {
	DatabaseAddr     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string

	ApiAddr string
}

func (c Config) Validate() error {
	if c.ApiAddr == "" {
		return fmt.Errorf("field 'Addr' cannot be empty")
	}

	return nil
}
