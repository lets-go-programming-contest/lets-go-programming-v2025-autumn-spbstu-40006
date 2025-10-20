package converter

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/dizey5k/task-3/internal/config"
	"github.com/dizey5k/task-3/internal/currency"
	"golang.org/x/text/encoding/charmap"
)

var ErrUnknownCharset = errors.New("unknown charset")

const dirPerm = 0o755

func decode(data []byte, out interface{}) error {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return nil, fmt.Errorf("%w: %s", ErrUnknownCharset, charset)
		}
	}

	if err := decoder.Decode(out); err != nil {
		return fmt.Errorf("failed to decode XML: %w", err)
	}

	return nil
}

func encode(path string, data any) error {
	if err := os.MkdirAll(filepath.Dir(path), dirPerm); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create output file '%s': %w", path, err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")

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

	var xmlCurrencies currency.XMLCurrencies
	if err := decode(data, &xmlCurrencies); err != nil {
		return fmt.Errorf("failed to parse XML: %w", err)
	}

	jsonCurrencies := make([]currency.JSONCurrency, 0, len(xmlCurrencies.Currencies))

	for _, xmlCurrency := range xmlCurrencies.Currencies {
		jsonCurrency, err := xmlCurrency.ToJSONCurrency()
		if err != nil {
			return fmt.Errorf("failed to convert currency '%s': %w", xmlCurrency.CharCode, err)
		}

		jsonCurrencies = append(jsonCurrencies, jsonCurrency)
	}

	sort.Slice(jsonCurrencies, func(i, j int) bool {
		return jsonCurrencies[i].Value > jsonCurrencies[j].Value
	})

	if err := encode(cfg.OutputFile, jsonCurrencies); err != nil {
		return err
	}

	return nil
}
