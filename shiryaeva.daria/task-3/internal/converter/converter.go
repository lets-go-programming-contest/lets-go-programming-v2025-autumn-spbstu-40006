package converter

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/Dora-shi/task-3/internal/config"
	"github.com/Dora-shi/task-3/internal/currency"
	"golang.org/x/text/encoding/charmap"
)

const dirPerm = 0o755

func decodeXML(data []byte) (*currency.ValCurs, error) {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return input, nil
		}
	}

	var valCurs currency.ValCurs
	if err := decoder.Decode(&valCurs); err != nil {
		return nil, fmt.Errorf("failed to decode XML: %w", err)
	}

	return &valCurs, nil
}

func saveJSON(data []currency.JSONCurrency, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), dirPerm); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create output file '%s': %w", path, err)
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fmt.Printf("warning: failed to close file: %v\n", closeErr)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

func Process(cfg *config.Config) error {
	data, err := os.ReadFile(cfg.InputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file '%s': %w", cfg.InputFile, err)
	}

	valCurs, err := decodeXML(data)
	if err != nil {
		return fmt.Errorf("failed to parse XML: %w", err)
	}

	jsonCurrencies := make([]currency.JSONCurrency, 0, len(valCurs.Currencies))
	for _, curr := range valCurs.Currencies {
		jsonCurrencies = append(jsonCurrencies, curr.ToJSON())
	}

	sort.Slice(jsonCurrencies, func(i, j int) bool {
		return jsonCurrencies[i].Value > jsonCurrencies[j].Value
	})

	if err := saveJSON(jsonCurrencies, cfg.OutputFile); err != nil {
		return err
	}

	return nil
}
