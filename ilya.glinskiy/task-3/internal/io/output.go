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
		return fmt.Errorf("Couldn't create directory for OutputFile")
	}

	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Errorf("Couldn't marshal valutes into json")
	}

	err = os.WriteFile(path, jsonData, 0600)
	if err != nil {
		return fmt.Errorf("Couldn't write an output file")
	}

	return nil
}
