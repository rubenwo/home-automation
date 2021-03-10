package rtsp

import (
	"fmt"
	"net/url"
)

type Client struct {
	Id   int64
	Host url.URL
}

func NewClient(cfg Config) (*Client, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("error validating config: %w", err)
	}

	return &Client{}, nil
}

func (c *Client) Connect() error {
	return nil
}
