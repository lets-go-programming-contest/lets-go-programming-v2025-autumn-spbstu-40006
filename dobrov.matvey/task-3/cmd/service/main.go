package main

import (
	"github.com/HorekProgrammer/task-3/internal/config"
	"github.com/HorekProgrammer/task-3/internal/currency"
)

func main() {
	configPath := config.GetConfigPath()

	var (
		cfg  config.Config
		curs currency.ValCurs
	)

	err := config.Read(&cfg, configPath)

	if err != nil {
		panic(err.Error())
	}

	err = currency.ReadDataFileNCanGetCurs(&curs, cfg.InputFile)

	if err != nil {
		panic(err.Error())
	}

	rates := currency.FillNSortRates(&curs)

	err = currency.FillOutputFile(rates, cfg.OutputFile)

	if err != nil {
		panic(err.Error())
	}
}
