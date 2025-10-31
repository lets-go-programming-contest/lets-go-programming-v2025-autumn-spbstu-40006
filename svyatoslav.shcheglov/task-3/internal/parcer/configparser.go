package parser

import (
	"flag"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type ApplicationSettings struct {
	SourceFilePath string `yaml:"input-file"`
	ResultFilePath string `yaml:"output-file"`
}

func LoadApplicationSettings() (ApplicationSettings, error) {
	settingsFile := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	fileData, err := os.ReadFile(*settingsFile)
	if err != nil {
		return ApplicationSettings{}, fmt.Errorf("read configuration: %w", err)
	}

	var appConfig ApplicationSettings

	err = yaml.Unmarshal(fileData, &appConfig)
	if err != nil {
		return ApplicationSettings{}, fmt.Errorf("process configuration: %w", err)
	}

	return appConfig, nil
}
