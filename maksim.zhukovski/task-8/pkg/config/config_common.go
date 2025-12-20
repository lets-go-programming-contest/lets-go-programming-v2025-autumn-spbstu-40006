package config

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Environment string `json:"environment" yaml:"environment"`
	LogLevel    string `json:"log_level"   yaml:"log_level"`
}

func Load(data []byte) (Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		if err2 := json.Unmarshal(data, &cfg); err2 != nil {
			return Config{}, fmt.Errorf("failed to parse config: %w", err)
		}
	}

	return cfg, nil
}
