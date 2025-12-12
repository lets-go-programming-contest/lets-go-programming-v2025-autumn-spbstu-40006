package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"github.com/abdelrhmanbaha/task-3/internal/config"
	"github.com/abdelrhmanbaha/task-3/internal/parser"
)

func main() {
	configPath := flag.String("config", "", "path to YAML config file")
	flag.Parse()

	if *configPath == "" {
		fmt.Fprintln(os.Stderr, "missing required flag: --config")
		os.Exit(1)
	}

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	valutes, err := parser.ParseXMLFile(cfg.InputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse XML: %v\n", err)
		os.Exit(1)
	}

	// Sorting logic: Value (desc), then CharCode (asc)
	sort.Slice(valutes, func(i, j int) bool {
		if valutes[i].Value == valutes[j].Value {
			return valutes[i].CharCode < valutes[j].CharCode
		}
		return valutes[i].Value > valutes[j].Value
	})

	if err := parser.SaveToJSON(cfg.OutputFile, valutes); err != nil {
		fmt.Fprintf(os.Stderr, "failed to save JSON: %v\n", err)
		os.Exit(1)
	}
}