package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Segfault-chan/task-3/internal/rates"
)

const (
	dirPerm  = 0o755
	filePerm = 0o644
)

func WriteJSON(list []rates.Currency, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), dirPerm); err != nil {
		return fmt.Errorf("couldn't create a directory: %w", err)
	}

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, filePerm)
	if err != nil {
		return fmt.Errorf("couldn't open/create a file: %w", err)
	}

	encoder := json.NewEncoder(file)

	if err := encoder.Encode(list); err != nil {
		return fmt.Errorf("problem while writing json: %w", err)
	}

	return nil
}
