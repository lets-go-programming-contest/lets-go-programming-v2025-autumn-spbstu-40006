package currency

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

const (
	filePerm = 0o644
	dirPerm  = 0o755
)

func SortValues(currencies *ValCurs) {
	sort.Slice(currencies.Currencies, func(i, j int) bool {
		return currencies.Currencies[i].Value > currencies.Currencies[j].Value
	})
}

func SaveToJSON(filePath string, currencies *ValCurs) error {
	dir := filepath.Dir(filePath)

	err := os.MkdirAll(dir, dirPerm)
	if err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := json.MarshalIndent(currencies.Currencies, "", " ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	err = os.WriteFile(filePath, data, filePerm)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
