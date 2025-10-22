package main

import (
	"flag"
	"fmt"

	"github.com/Mishaa105/task-3/internal/config"
	"github.com/Mishaa105/task-3/internal/decoding"
	"github.com/Mishaa105/task-3/internal/saving"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config")
	outputFlag := flag.String("output", "", "path to output file (overrides config)")
	flag.Parse()

	cfg, err := config.CheckInput(*configPath)
	if err != nil {
		panic(err)
	}

	outputPath := *outputFlag
	if outputPath == "" {
		outputPath = cfg.OutputFile
	}

	valCurs := decoding.Decoding(*configPath)
	saving.SaveToJSON(outputPath, valCurs.Valutes)

	for _, val := range valCurs.Valutes {
		fmt.Printf("NumCode: %d, CharCode: %s, Value: %.2f\n", val.NumCode, val.CharCode, val.Value)
	}
}
