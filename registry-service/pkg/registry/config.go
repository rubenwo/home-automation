package registry

type Config struct {
	DatabaseBackend    string
	PgDatabaseBackend  string
	PgDatabaseUser     string
	PgDatabasePassword string
	PgDatabaseName     string
}

func (c *Config) Validate() error {
	return nil
}
