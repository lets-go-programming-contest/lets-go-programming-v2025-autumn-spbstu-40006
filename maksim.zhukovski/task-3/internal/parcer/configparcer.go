package parcer

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

func ParseConfig() (Config, error) {
	configPath := flag.String("config", "config.yaml", "Path to the configuration file")
	flag.Parse()

	data, err := os.ReadFile(*configPath)
	if err != nil {
		return Config{}, fmt.Errorf("read yanl: %w", err)
	}

	var cfg Config

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return Config{}, fmt.Errorf("decode yaml: %w", err)
	}

	return cfg, nil
}
