//go:build dev

package config

import (
	_ "embed"
	"fmt"

	"gopkg.in/yaml.v3"
)

//go:embed dev.yaml
var devConfig []byte

func GetConfig() (*Config, error) {
	var cfg Config
	err := yaml.Unmarshal(devConfig, &cfg)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения dev конфигурации: %w", err)
	}
	return &cfg, nil
}
