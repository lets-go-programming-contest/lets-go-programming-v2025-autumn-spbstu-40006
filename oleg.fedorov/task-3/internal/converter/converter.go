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

var (
	ErrUnknownCharset = errors.New("unknown charset")
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func decode(data []byte, value interface{}) error {
	decoder := xml.NewDecoder(bytes.NewReader(data))
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return nil, fmt.Errorf("%w: %s", ErrUnknownCharset, charset)
		}
	}

	return decoder.Decode(value)
}

func Process(cfg *config.Config) {
	data, err := os.ReadFile(cfg.InputFile)
	panicIfErr(fmt.Errorf("failed to read input file '%s': %w", cfg.InputFile, err))

	var xmlCurrencies currency.XMLCurrencies
	err = decode(data, &xmlCurrencies)
	panicIfErr(fmt.Errorf("failed to parse XML data: %w", err))

	jsonCurrencies := make([]currency.JSONCurrency, 0, len(xmlCurrencies.Currencies))

	for _, xmlCurrency := range xmlCurrencies.Currencies {
		jsonCurrency, err := xmlCurrency.ToJSONCurrency()
		panicIfErr(fmt.Errorf("failed to convert currency: %w", err))

		jsonCurrencies = append(jsonCurrencies, jsonCurrency)
	}

	sort.Slice(jsonCurrencies, func(i, j int) bool {
		return jsonCurrencies[i].Value > jsonCurrencies[j].Value
	})

	err = os.MkdirAll(filepath.Dir(cfg.OutputFile), 0755)

	panicIfErr(fmt.Errorf("failed to create output directory: %w", err))

	file, err := os.Create(cfg.OutputFile)
	panicIfErr(fmt.Errorf("failed to create output file '%s': %w", cfg.OutputFile, err))

	defer func() {
		if err := file.Close(); err != nil {
			panic(fmt.Errorf("failed to close output file: %w", err))
		}
	}()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")
	if err := encoder.Encode(jsonCurrencies); err != nil {
		panicIfErr(fmt.Errorf("failed to encode JSON: %w", err))
	}
}
