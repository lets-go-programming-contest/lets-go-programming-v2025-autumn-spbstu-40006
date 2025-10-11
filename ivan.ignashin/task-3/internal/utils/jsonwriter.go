package utils

import (
	"encoding/json"
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

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	jsonData, err := json.MarshalIndent(finalRecords, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, jsonData, 0600)
}
