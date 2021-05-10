package foodservice

type Config struct {
	DatabaseAddr     string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string

	ApiAddr string
}

func (c *Config) Validate() error {
	return nil
}
