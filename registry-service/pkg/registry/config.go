package registry

type Config struct {
	DatabaseBackend string
}

func (c *Config) Validate() error {
	return nil
}
