package processing

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func SaveJSON(path string, currencies []Currency) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	if err := json.NewEncoder(file).Encode(currencies); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}
