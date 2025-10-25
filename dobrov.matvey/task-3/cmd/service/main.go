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

	config.Read(&cfg, configPath)

	currency.ReadDataFileNCanGetCurs(&curs, cfg.InputFile)

	rates := currency.FillNSortRates(&curs)

	currency.FillOutputFile(rates, cfg.OutputFile)
}
