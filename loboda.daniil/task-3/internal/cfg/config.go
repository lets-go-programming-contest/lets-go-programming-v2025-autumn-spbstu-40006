package cfg

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

var ErrConfigPathsEmpty = errors.New("both input-file and output-file must be set in config")

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func Load(path string) (Config, error) {
	body, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config file: %w", err)
	}

	var conf Config
	if err := yaml.Unmarshal(body, &conf); err != nil {
		return Config{}, fmt.Errorf("yaml decode: %w", err)
	}

	conf.InputFile = strings.TrimSpace(conf.InputFile)
	conf.OutputFile = strings.TrimSpace(conf.OutputFile)

	if conf.InputFile == "" || conf.OutputFile == "" {
		return Config{}, ErrConfigPathsEmpty
	}

	return conf, nil
}
