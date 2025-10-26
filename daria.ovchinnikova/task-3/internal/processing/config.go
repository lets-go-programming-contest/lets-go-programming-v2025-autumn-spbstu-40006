package processing

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func LoadConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var cnfg Config
	if err := yaml.Unmarshal(data, &cnfg); err != nil {
		return Config{}, fmt.Errorf("failed to parse config file: %w", err)
	}

	return cnfg, nil
}
