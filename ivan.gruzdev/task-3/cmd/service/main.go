package main

import (
	"flag"

	"github.com/MoneyprogerISG/task-3/pkg/config"
	"github.com/MoneyprogerISG/task-3/pkg/currency"
)

func main() {

	configPath := flag.String("config", "config.yaml", "path to config file")

	flag.Parse()

	manual := config.LoadConfig((*configPath))

	currencies := currency.LoadCurrencies(manual.InputFile)

	currency.SortValues(&currencies)

	currency.SaveToJson(manual.OutputFile, currencies)

}
