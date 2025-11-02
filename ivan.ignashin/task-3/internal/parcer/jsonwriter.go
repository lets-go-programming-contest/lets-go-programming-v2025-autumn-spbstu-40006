package parcer

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

const (
	dirPerm  = 0o755
	filePerm = 0o644
)

func SaveAsJSON(items []Record, path string) error {
	sort.Slice(items, func(i, j int) bool {
		return items[i].Value > items[j].Value
	})

	dir := filepath.Dir(path)

	err := os.MkdirAll(dir, dirPerm)
	if err != nil {
		return fmt.Errorf("create dir %s: %w", dir, err)
	}

	jsonData, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json %w", err)
	}

	err = os.WriteFile(path, jsonData, filePerm)
	if err != nil {
		return fmt.Errorf("write file %s: %w", path, err)
	}

	return nil
}
