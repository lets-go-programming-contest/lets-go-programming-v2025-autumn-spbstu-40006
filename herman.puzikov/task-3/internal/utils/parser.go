package utils

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func ParseYAML(dst *map[string]string, filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open YAML file: %w", err)
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(dst); err != nil {
		return fmt.Errorf("failed to decode YAML: %w", err)
	}

	return nil
}

// func parseXML()
