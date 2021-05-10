package ingress

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	ApiVersion string `yaml:"apiVersion"`
	Spec       []struct {
		Host   string `yaml:"host"`
		Routes []struct {
			Path                  string   `yaml:"path"`
			Methods               []string `yaml:"methods"`
			Protocol              string   `yaml:"protocol"`
			RequiresAuthorization bool     `yaml:"required_authorization"`
			AllowedRoles          []string `yaml:"allowed_roles"`
		} `yaml:"routes"`
	} `yaml:"spec"`
}

func ParseConfig(path string) (*Config, error) {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("couldn't read file: %s, with error: %w", path, err)
	}
	var c Config
	if err := yaml.Unmarshal(raw, &c); err != nil {
		return nil, fmt.Errorf("error parsing: %w", err)
	}

	return &c, nil
}
