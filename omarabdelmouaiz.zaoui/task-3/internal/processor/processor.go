package processor

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"golang.org/x/net/html/charset"
	"omarabdelmouaiz.zaoui/task-3/internal/config"
	"omarabdelmouaiz.zaoui/task-3/internal/models"
)

const (
	dirPermissions  = 0o755
	filePermissions = 0o600
)

func Run(cfg *config.Config) error {
	xmlData, err := readXMLData(cfg.InputFile)
	if err != nil {
		return fmt.Errorf("error reading xml file: %w", err)
	}

	valutes := transformAndSort(xmlData)

	err = writeJSONData(cfg.OutputFile, valutes)
	if err != nil {
		return fmt.Errorf("error writing json file: %w", err)
	}

	return nil
}

func readXMLData(path string) (*models.ValCurs, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Printf("warning: failed to close file: %v\n", closeErr)
		}
	}()

	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel

	var valCurs models.ValCurs
	if err := decoder.Decode(&valCurs); err != nil {
		return nil, fmt.Errorf("failed to decode XML: %w", err)
	}

	return &valCurs, nil
}

func transformAndSort(xmlData *models.ValCurs) []models.Valute {
	valutes := make([]models.Valute, len(xmlData.Valutes))
	copy(valutes, xmlData.Valutes)

	sort.Slice(valutes, func(i, j int) bool {
		return float64(valutes[i].Value) > float64(valutes[j].Value)
	})

	return valutes
}

func writeJSONData(path string, data []models.Valute) error {
	dir := filepath.Dir(path)

	if err := os.MkdirAll(dir, dirPermissions); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(path, jsonData, filePermissions); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
