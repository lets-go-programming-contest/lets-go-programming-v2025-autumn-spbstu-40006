package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type FinalRecord struct {
	NumCode  int     `json:"NumCode"`
	CharCode string  `json:"CharCode"`
	Value    float64 `json:"Value"`
}

func SaveAsJSON(items []Record, path string) error {
	finalRecords := make([]FinalRecord, len(items))
	for i, item := range items {
		finalRecords[i] = FinalRecord{
			NumCode:  item.ID,
			CharCode: item.Name,
			Value:    item.Value,
		}
	}

	sort.Slice(finalRecords, func(i, j int) bool {
		return finalRecords[i].Value < finalRecords[j].Value
	})

	dir := filepath.Dir(path)
	err := os.MkdirAll(dir, 0755)

	if err != nil {
		return fmt.Errorf("create dir %s: %w", dir, err)
	}

	jsonData, err := json.MarshalIndent(finalRecords, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json %w", err)
	}

	return os.WriteFile(path, jsonData, 0600)
}
