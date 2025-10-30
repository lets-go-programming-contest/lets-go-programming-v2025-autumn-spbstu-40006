package main

import (
	"github.com/HorekProgrammer/task-3/internal/config"
	"github.com/HorekProgrammer/task-3/internal/currency"
)

func main() {
	configPath := config.GetConfigPath()

	var curs currency.ValCurs

	cfg, err := config.Read(configPath)
	if err != nil {
		panic(err)
	}

	err = currency.ReadDataFileNCanGetCurs(&curs, cfg.InputFile)
	if err != nil {
		panic(err)
	}

	rates := currency.FillNSortRates(&curs)

	err = currency.FillOutputFile(rates, cfg.OutputFile)
	if err != nil {
		panic(err)
	}
}
