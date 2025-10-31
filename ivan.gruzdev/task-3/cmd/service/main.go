package main

import (
	"flag"

	"github.com/MoneyprogerISG/task-3/pkg/config"
	"github.com/MoneyprogerISG/task-3/pkg/currency"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")

	flag.Parse()

	cfg, err := config.LoadConfig((*configPath))
	if err != nil {
		panic(err)
	}

	currencies, err := currency.LoadCurrencies(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	currency.SortValues(&currencies)

	err = currency.SaveToJSON(cfg.OutputFile, &currencies)
	if err != nil {
		panic(err)
	}
}
