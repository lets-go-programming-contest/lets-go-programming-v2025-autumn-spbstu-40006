package config

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

func GetConfigPath() string {
	configPathPtr := flag.String("config", "", "config.yaml")
	flag.Parse()

	return *configPathPtr
}

func Read(cfg *Config, path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(fmt.Errorf("open %q: %w", path, err))
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		panic(err)
	}
}
