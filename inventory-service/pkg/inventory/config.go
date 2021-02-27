package inventory

type Config struct {
	DatabaseAddr     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
}

func (c *Config) Validate() error {
	return nil
}
