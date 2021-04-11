package actions

type Config struct {
	Addr             string
	IntentProcessors map[string]IntentProcessor
}

func (c *Config) Validate() error { return nil }
