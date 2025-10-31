package processor

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func (cp *CurrencyProcessor) LoadConfig() error {
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		return fmt.Errorf("config file not found: %s", *configPath)
	}

	data, err := os.ReadFile(*configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse YAML config: %w", err)
	}

	if config.InputFile == "" {
		return fmt.Errorf("input-file cannot be empty")
	}
	if config.OutputFile == "" {
		return fmt.Errorf("output-file cannot be empty")
	}

	cp.config = &config
	return nil
}
