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

	cfg := processing.LoadConfig(*configPath)

	valCurs := processing.LoadXML(cfg.InputFile)

	processing.SortCurrencies(valCurs)

	processing.SaveJSON(cfg.OutputFile, valCurs)
}
