package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

const (
	directoryPermissions = 0o755
	filePermissions      = 0o644
)

func ExportToJSON(currencyEntries []CurrencyData, outputPath string) error {
	sort.Slice(currencyEntries, func(firstIndex, secondIndex int) bool {
		return currencyEntries[firstIndex].Rate > currencyEntries[secondIndex].Rate
	})

	outputDirectory := filepath.Dir(outputPath)

	err := os.MkdirAll(outputDirectory, directoryPermissions)
	if err != nil {
		return fmt.Errorf("create output directory %s: %w", outputDirectory, err)
	}

	formattedJSON, err := json.MarshalIndent(currencyEntries, "", "  ")
	if err != nil {
		return fmt.Errorf("format JSON data: %w", err)
	}

	err = os.WriteFile(outputPath, formattedJSON, filePermissions)
	if err != nil {
		return fmt.Errorf("write JSON file %s: %w", outputPath, err)
	}

	return nil
}
