package main

import (
	"github.com/HorekProgrammer/task-3/internal/app"
)

func main() {
	configPath := app.GetConfigPath()

	var (
		cfg  app.Config
		curs app.ValCurs
	)

	err := app.ReadDataFromConfig(&cfg, configPath)

	if err != nil {
		return
	}

	err = app.ReadDataFileNCanGetCurs(&curs, cfg.InputFile)

	if err != nil {
		return
	}

	rates := app.FillNSortRates(curs)

	err = app.FillOutputFile(rates, cfg)

	if err != nil {
		return
	}
}
