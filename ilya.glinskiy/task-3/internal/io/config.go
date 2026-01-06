package io

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func ReadConfig(path string, config *Config) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("couldn't open config file: %w", err)
	}

	err = yaml.Unmarshal(content, config)
	if err != nil {
		return fmt.Errorf("couldn't unmarshal config file: %w", err)
	}

	return nil
}
