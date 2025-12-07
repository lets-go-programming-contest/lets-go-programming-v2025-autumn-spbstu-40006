package config

import (
	yaml "github.com/goccy/go-yaml"
)

type Config struct {
	Environment string
	LogLevel    string `yaml:"log_level"`
}

func New() *Config {
	config := &Config{} //nolint:exhaustruct
	if err := yaml.Unmarshal([]byte(configContents), config); err != nil {
		panic(err)
	}

	return config
}