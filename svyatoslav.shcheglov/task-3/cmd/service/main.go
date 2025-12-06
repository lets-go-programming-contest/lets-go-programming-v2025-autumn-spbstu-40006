package main

import (
	parser "github.com/SpeaarIt/task-3/internal/parcer"
)

func main() {
	config, err := parser.LoadApplicationSettings()
	if err != nil {
		panic(err)
	}

	currencyData, err := parser.LoadCurrencyData(config.SourceFilePath)
	if err != nil {
		panic(err)
	}

	err = parser.ExportToJSON(currencyData, config.ResultFilePath)
	if err != nil {
		panic(err)
	}
}
