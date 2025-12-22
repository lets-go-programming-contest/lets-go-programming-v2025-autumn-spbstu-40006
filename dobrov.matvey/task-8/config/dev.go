//go:build dev

package config

import (
	"errors"

	_ "embed"
	"gopkg.in/yaml.v3"
)

//go:embed dev.yaml
var devConfig []byte

var ErrUnmarshal = errors.New("cant unmarshal yaml")

func Load() (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(devConfig, &cfg); err != nil {
		return nil, ErrUnmarshal
	}

	return &cfg, nil
}
