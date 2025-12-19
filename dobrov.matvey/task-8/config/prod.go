//go:build !dev

package config

import (
	_ "embed"
	"errors"

	"gopkg.in/yaml.v3"
)

//go:embed prod.yaml
var prodConfig []byte

var ErrUnmarshal = errors.New("cant unmarshal yaml")

func Load() (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(prodConfig, &cfg); err != nil {
		return nil, ErrUnmarshal
	}

	return &cfg, nil
}
