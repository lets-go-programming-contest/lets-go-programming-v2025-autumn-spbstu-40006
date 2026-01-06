package io

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type JSONValute struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func WriteOutput(path string, output []JSONValute) error {
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
