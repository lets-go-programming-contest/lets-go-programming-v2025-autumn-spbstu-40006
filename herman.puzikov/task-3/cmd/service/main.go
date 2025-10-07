package main

import (
	"flag"
	"log"
	"slices"

	"github.com/Segfault-chan/task-3/internal/utils"
)

func main() {
	configPath := flag.String("config", "", "path to YAML config file")
	flag.Parse()

	if *configPath == "" {
		log.Panic("no config file provided")
	}

	configFile, err := utils.ParseYAML(*configPath)
	if err != nil {
		log.Panic(err)
	}

	exchangeRates, err := utils.ParseXML(configFile.InputFile)
	if err != nil {
		log.Panic(err)
	}

	slices.SortFunc(exchangeRates.Currencies, utils.DescendingComparatorCurrency)

	if err := utils.ParseJSON(exchangeRates.Currencies, configFile.OutputFile); err != nil {
		log.Panic(err)
	}
}
