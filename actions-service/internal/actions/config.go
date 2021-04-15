package actions

import "github.com/rubenwo/home-automation/actions-service/internal/actions/intentprocessor"

type Config struct {
	Addr             string
	IntentProcessors map[string]intentprocessor.IntentProcessor
}

func (c *Config) Validate() error { return nil }
