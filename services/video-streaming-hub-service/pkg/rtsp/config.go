package rtsp

type Config struct {
	Id   int64
	Host string
}

func (c Config) Validate() error {
	return nil
}
