package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"task-3/internal/config"
	"task-3/internal/currencies"
)

func main() {
	configPath := flag.String("config", "", "Path to config")
	flag.Parse()

	if *configPath == "" {
		fmt.Fprintln(os.Stderr, "Error: config path is required")
		os.Exit(1)
	}

	cfg, err := config.New(*configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to load config: %v\n", err)
		os.Exit(1)
	}

	curr, err := currencies.New(cfg.InputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to load currencies: %v\n", err)
		os.Exit(1)
	}

	sort.Slice(curr.Currencies, func(i, j int) bool {
		return curr.Currencies[i].Value > curr.Currencies[j].Value
	})

	if err := curr.SaveToOutputFile(cfg.OutputFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to save to output file: %v\n", err)
		os.Exit(1)
	}
}
