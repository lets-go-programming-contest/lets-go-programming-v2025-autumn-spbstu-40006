package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Mishaa105/task-3/internal/decoding"
	"github.com/Mishaa105/task-3/internal/saving"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config")
	outputPath := flag.String("output", "output/result.json", "path to output file")
	flag.Parse()

	absOutputPath, err := filepath.Abs(*outputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving absolute output path: %v\n", err)
		os.Exit(1)
	}

	valCurs := decoding.Decoding(*configPath)
	saving.SaveToJSON(absOutputPath, valCurs.Valutes)

	for _, val := range valCurs.Valutes {
		fmt.Printf("NumCode: %d, CharCode: %s, Value: %.2f\n", val.NumCode, val.CharCode, val.Value)
	}
}
