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

func ReadConfig(path string, config interface{}) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Couldn't open config file")
	}

	err = yaml.Unmarshal(content, config)
	if err != nil {
		return fmt.Errorf("Couldn't read config file")
	}

	return nil
}
