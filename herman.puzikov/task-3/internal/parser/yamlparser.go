package parser

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Segfault-chan/task-3/internal/config"
	"github.com/Segfault-chan/task-3/internal/rates"
	yaml "gopkg.in/yaml.v3"
)

func ReadYAML(filename string) (*config.Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	decoder := yaml.NewDecoder(file)
	decoder.KnownFields(true)

	var config config.Config
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("error decoding YAML: %w", err)
	}

	return &config, nil
}

func WriteYAML(list []rates.Currency, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), dirPerm); err != nil {
		return fmt.Errorf("couldn't create a directory: %w", err)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, filePerm)
	if err != nil {
		return fmt.Errorf("couldn't open/create a file: %w", err)
	}

	encoder := yaml.NewEncoder(file)

	if err := encoder.Encode(list); err != nil {
		return fmt.Errorf("problem while writing json: %w", err)
	}

	return nil
}
