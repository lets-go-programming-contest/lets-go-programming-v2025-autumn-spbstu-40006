package config

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `yaml:"environment"`
	LogLevel    string `yaml:"log_level"`
}

type Loader interface {
	Load() ([]byte, error)
}

func Load(data []byte) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal YAML: %w", err)
	}

	if cfg.Environment == "" {
		return nil, fmt.Errorf("environment field is required")
	}
	if cfg.LogLevel == "" {
		return nil, fmt.Errorf("log_level field is required")
	}

	return &cfg, nil
}

func GetConfig(loader Loader) (*Config, error) {
	data, err := loader.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return Load(data)
}
