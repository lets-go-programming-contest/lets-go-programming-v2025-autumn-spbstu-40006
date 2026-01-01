package main

import (
	"flag"
	"sort"

	"github.com/vitsh1/task-3/internal/config"
	"github.com/vitsh1/task-3/internal/parser"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	cfg, err := config.ReadConfig(*configPath)

	if err != nil {
		panic(err)
	}

	valCurs, err := parser.ParseXML(cfg.InputFile)

	if err != nil {
		panic(err)
	}

	sort.Slice(valCurs, func(i, k int) bool {
		return valCurs[i].Value > valCurs[k].Value
	})

	if err := parser.WriteJSON(cfg.OutputFile, valCurs); err != nil {
		panic(err)
	}
}
