package main

import (
	"flag"
	"sort"

	"github.com/filon6/task-3/internal/config"
	"github.com/filon6/task-3/internal/parser"
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

	valutes, err := parser.ParseXMLFile(cfg.InputFile)
	if err != nil {
		panic("failed to parse XML: " + err.Error())
	}

	sort.Slice(valutes, func(i, j int) bool {
		return valutes[i].Value > valutes[j].Value
	})

	if err := parser.SaveToJSON(cfg.OutputFile, valutes); err != nil {
		panic("failed to save JSON: " + err.Error())
	}
}
