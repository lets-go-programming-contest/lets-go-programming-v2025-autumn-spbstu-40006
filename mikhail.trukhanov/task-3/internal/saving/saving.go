package saving

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Mishaa105/task-3/internal/decoding"
)

func SaveToJSON(path string, data []decoding.Valute) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("cannot create output directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot create output file: %w", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			fmt.Printf("failed to close file: %v\n", err)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("cannot write JSON: %w", err)
	}

	return nil
}
