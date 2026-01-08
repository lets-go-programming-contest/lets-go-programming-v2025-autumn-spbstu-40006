package main

import (
	"flag"
	"sort"

	"github.com/MargotBush/task-3/internal/config"
	"github.com/MargotBush/task-3/internal/parser"
)

func main() {
	configPath := flag.String("config", "", "path to YAML config file")
	flag.Parse()

	if *configPath == "" {
		panic("missing required flag: --config")
	}

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		panic("failed to load config: " + err.Error())
	}

	currencies, err := parser.ParseXMLFile(cfg.InputFile)
	if err != nil {
		panic("failed to parse XML: " + err.Error())
	}

	sort.Slice(currencies, func(i, j int) bool {
		return currencies[i].Value > currencies[j].Value
	})

	if err := parser.SaveToJSON(cfg.OutputFile, currencies); err != nil {
		panic("failed to save JSON: " + err.Error())
	}
}
