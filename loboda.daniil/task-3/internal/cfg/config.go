package cfg

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func Load(path string) (Config, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config file: %w", err)
	}

	var c Config
	if err := yaml.Unmarshal(body, &c); err != nil {
		return Config{}, fmt.Errorf("yaml decode: %w", err)
	}

	c.InputFile = strings.TrimSpace(c.InputFile)
	c.OutputFile = strings.TrimSpace(c.OutputFile)

	if c.InputFile == "" || c.OutputFile == "" {
		return Config{}, errors.New("both input-file and output-file must be set in config")
	}
	return c, nil
}
