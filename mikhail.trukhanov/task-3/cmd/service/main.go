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

	valCurs, err := decoding.Decoding(*configPath)
	if err != nil {
		fmt.Println("Error while decoding:", err)
		return
	}

	saving.SaveToJSON(outputPath, valCurs.Valutes)
}
