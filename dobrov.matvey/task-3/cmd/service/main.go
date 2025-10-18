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

	app.ReadDataFromConfig(&cfg, configPath)

	app.ReadDataFileNCanGetCurs(&curs, cfg.InputFile)

	rates := app.FillNSortRates(curs)

	app.FillOutputFile(rates, cfg)
}
