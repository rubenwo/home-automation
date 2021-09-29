package plugins

type Plugin interface {
	Run(cfg Config) error
}
