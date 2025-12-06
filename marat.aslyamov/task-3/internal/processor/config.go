package processor

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

var (
	ErrConfigNotFound  = errors.New("config file not found")
	ErrInputFileEmpty  = errors.New("input-file cannot be empty")
	ErrOutputFileEmpty = errors.New("output-file cannot be empty")
	ErrReadConfig      = errors.New("failed to read config file")
	ErrParseYAML       = errors.New("failed to parse YAML config")
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func (cp *CurrencyProcessor) LoadConfig() error {
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	if _, err := os.Stat(*configPath); os.IsNotExist(err) {
		return fmt.Errorf("%w: %s", ErrConfigNotFound, *configPath)
	}

	data, err := os.ReadFile(*configPath)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrReadConfig, err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("%w: %w", ErrParseYAML, err)
	}

	if config.InputFile == "" {
		return ErrInputFileEmpty
	}

	if config.OutputFile == "" {
		return ErrOutputFileEmpty
	}

	cp.config = &config

	return nil
}
