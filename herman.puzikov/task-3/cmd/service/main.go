package main

import (
	"flag"
	"log"
	"slices"

	"github.com/Segfault-chan/task-3/internal/parser"
	"github.com/Segfault-chan/task-3/internal/rates"
)

func main() {
	configPath := flag.String("config", "", "path to YAML config file")
	outputFormat := flag.String("output-format", "json", "output file format (either json, xml or yaml)")
	flag.Parse()

	if *configPath == "" {
		log.Panic("no config file provided")
	}

	configFile, err := parser.ReadYAML(*configPath)
	if err != nil {
		log.Panic(err)
	}

	exchangeRates, err := parser.ReadXML(configFile.InputFile)
	if err != nil {
		log.Panic(err)
	}

	slices.SortFunc(exchangeRates.Currencies, rates.CompareByValueDesc)

	switch *outputFormat {
	case "json":
		if err := parser.WriteJSON(exchangeRates.Currencies, configFile.OutputFile); err != nil {
			log.Panic(err)
		}
	case "xml":
		if err := parser.WriteXML(exchangeRates.Currencies, configFile.OutputFile); err != nil {
			log.Panic(err)
		}
	case "yaml":
		if err := parser.WriteYAML(exchangeRates.Currencies, configFile.OutputFile); err != nil {
			log.Panic(err)
		}
	}
}
