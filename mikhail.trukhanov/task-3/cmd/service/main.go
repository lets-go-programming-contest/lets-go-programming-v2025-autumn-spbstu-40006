package main

import (
	"flag"
	"fmt"

	"github.com/Mishaa105/task-3/internal/decoding"
	"github.com/Mishaa105/task-3/internal/saving"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path")
	flag.Parse()

	valCurs := decoding.Decoding(*configPath)

	for _, val := range valCurs.Valutes {
		fmt.Printf("NumCode: %d, CharCode: %s, Value: %.2f\n", val.NumCode, val.CharCode, val.Value)
	}

	saving.SaveToJSON("output/result.json", valCurs.Valutes)
}
