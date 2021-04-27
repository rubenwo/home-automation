package registry

type Config struct {
	DatabaseBackend  string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
}

func (c *Config) Validate() error {
	return nil
}
