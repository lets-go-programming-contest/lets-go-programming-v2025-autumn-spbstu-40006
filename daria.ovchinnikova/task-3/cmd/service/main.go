package main

import (
	"flag"

	"github.com/monka6/task-3/internal/processing"
)

func main() {
	configPath := flag.String("config", "", "path to config file")
	flag.Parse()

	if *configPath == "" {
		panic("Config file path is required")
	}

	cfg, err := processing.LoadConfig(*configPath)
	if err != nil {
		panic(err)
	}

	valCurs, err := processing.LoadXML(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	processing.SortCurrencies(valCurs)

	if err := processing.SaveJSON(cfg.OutputFile, valCurs); err != nil {
		panic(err)
	}
}
