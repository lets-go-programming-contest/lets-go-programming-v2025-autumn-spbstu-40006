package main

import (
	"flag"

	"github.com/MoneyprogerISG/task-3/pkg/config"
	"github.com/MoneyprogerISG/task-3/pkg/currency"
)

func main() {
	configPath := flag.String("config", "config.yaml", "path to config file")

	flag.Parse()

	cfg := config.LoadConfig((*configPath))

	currencies := currency.LoadCurrencies(cfg.InputFile)

	currency.SortValues(&currencies)

	currency.SaveToJSON(cfg.OutputFile, &currencies)
}
