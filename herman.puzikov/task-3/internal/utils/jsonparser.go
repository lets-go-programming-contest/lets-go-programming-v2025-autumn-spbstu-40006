package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func ParseJSON(source ExchangeRate, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("couldn't create a directory: %w", err)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("couldn't open/create a file: %w", err)
	}

	encoder := json.NewEncoder(file)

	if err := encoder.Encode(source); err != nil {
		return fmt.Errorf("problem while writing json: %w", err)
	}

	return nil
}
