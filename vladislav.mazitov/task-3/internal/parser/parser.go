package parser

import (
	"fmt"
	"os"
	"sort"

	"github.com/identicalaffiliation/task-3/internal/config"
	"github.com/identicalaffiliation/task-3/internal/currency"
	"github.com/identicalaffiliation/task-3/internal/decoder"
	"github.com/identicalaffiliation/task-3/internal/encoder"
)

func Process(cfg *config.Config) error {
	data, err := os.ReadFile(cfg.InputFile)
	if err != nil {
		return fmt.Errorf("read xml file: %w", err)
	}

	var valCurs currency.ValCurs

	err = decoder.Decode(data, &valCurs)
	if err != nil {
		return fmt.Errorf("parsing xml file: %w", err)
	}

	valutes := valCurs.Valutes
	sort.Slice(valutes, func(i, j int) bool {
		return float64(valutes[i].Value) > float64(valutes[j].Value)
	})

	err = encoder.SaveToJSON(cfg.OutputFile, valutes)
	if err != nil {
		return fmt.Errorf("saving to json file: %w", err)
	}

	return nil
}
