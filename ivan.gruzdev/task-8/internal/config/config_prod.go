//go:build !dev

package config

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:embed prod.yaml
var prodConfig []byte

func GetConfig() (*Config, error) {
	var cfg Config
	err := yaml.Unmarshal(prodConfig, &cfg)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения конфигурации: %w", err)
	}
	return &cfg, nil
}
