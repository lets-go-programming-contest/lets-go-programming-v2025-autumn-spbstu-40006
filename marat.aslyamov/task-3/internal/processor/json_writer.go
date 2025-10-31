package processor

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

type Currency struct {
	NumCode  int     `json:"num_code"  xml:"NumCode"`
	CharCode string  `json:"char_code" xml:"CharCode"`
	Value    float64 `json:"value"     xml:"Value"`
}

func (cp *CurrencyProcessor) SaveToJSON(currencies []Currency, outputPath string) error {
	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Value > currencies[j].Value
	})

	dir := filepath.Dir(outputPath)
	if err := os.MkdirAll(dir, dirPerm); err != nil {
		return fmt.Errorf("create dir: %w", err)
	}

	jsonData, err := json.MarshalIndent(currencies, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}

	if err := os.WriteFile(outputPath, jsonData, filePerm); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
