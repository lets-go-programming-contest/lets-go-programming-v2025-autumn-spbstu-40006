package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

const (
	dirPerm  = 0o755
	filePerm = 0o644
)

type FinalRecord struct {
	NumCode  int     `json:"num_code"`
	CharCode string  `json:"char_code"`
	Value    float64 `json:"value"`
}

func SaveAsJSON(items []Record, path string) error {
	finalRecords := make([]FinalRecord, len(items))

	for index, item := range items {
		valueFloat, err := strconv.ParseFloat(item.Value, 64)
		if err != nil {
			return fmt.Errorf("parse float %s: %w", item.Value, err)
		}

		finalRecords[index] = FinalRecord{
			NumCode:  item.ID,
			CharCode: item.Name,
			Value:    valueFloat,
		}
	}

	sort.Slice(finalRecords, func(i, j int) bool {
		return finalRecords[i].Value > finalRecords[j].Value
	})

	dir := filepath.Dir(path)

	err := os.MkdirAll(dir, dirPerm)
	if err != nil {
		return fmt.Errorf("create dir %s: %w", dir, err)
	}

	jsonData, err := json.MarshalIndent(finalRecords, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json %w", err)
	}

	err = os.WriteFile(path, jsonData, filePerm)
	if err != nil {
		return fmt.Errorf("write file %s: %w", path, err)
	}

	return nil
}
