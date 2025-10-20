package converter

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"path/filepath"
	"sort"

	"github.com/dizey5k/task-3/internal/config"
	"github.com/dizey5k/task-3/internal/currency"
)

func Process(cfg *config.Config) error {
	data, err := os.ReadFile(cfg.InputFile)
	if err != nil {
		return err
	}

	var xmlCurrencies currency.XMLCurrencies
	if err := xml.Unmarshal(data, &xmlCurrencies); err != nil {
		return err
	}

	jsonCurrencies := make([]currency.JSONCurrency, 0, len(xmlCurrencies.Currencies))
	for _, xmlCurrency := range xmlCurrencies.Currencies {
		jsonCurrency, err := xmlCurrency.ToJSONCurrency()
		if err != nil {
			return err
		}
		jsonCurrencies = append(jsonCurrencies, jsonCurrency)
	}

	sort.Slice(jsonCurrencies, func(i, j int) bool {
		return jsonCurrencies[i].Value > jsonCurrencies[j].Value
	})

	if err := os.MkdirAll(filepath.Dir(cfg.OutputFile), 0755); err != nil {
		return err
	}

	file, err := os.Create(cfg.OutputFile)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "\t")

	return encoder.Encode(jsonCurrencies)
}
