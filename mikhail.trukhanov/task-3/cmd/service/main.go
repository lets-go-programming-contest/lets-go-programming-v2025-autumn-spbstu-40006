package main

import (
	"flag"
	"fmt"

	"github.com/Mishaa105/task-3/internal/decoding"
	"github.com/Mishaa105/task-3/internal/saving"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config")
	outputPath := flag.String("output", ".output/result.json", "path to output file")
	flag.Parse()

	valCurs := decoding.Decoding(*configPath)
	saving.SaveToJSON(*outputPath, valCurs.Valutes)

	for _, val := range valCurs.Valutes {
		fmt.Printf("NumCode: %d, CharCode: %s, Value: %.2f\n", val.NumCode, val.CharCode, val.Value)
	}
}
