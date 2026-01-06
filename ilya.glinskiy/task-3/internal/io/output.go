package io

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func WriteOutput(path string, output []Valute) error {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return fmt.Errorf("couldn't create directory for output file: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("couldn't create output file: %w", err)
	}

	err = json.NewEncoder(file).Encode(output)
	if err != nil {
		return fmt.Errorf("couldn't encode json into output file: %w", err)
	}

	return nil
}
