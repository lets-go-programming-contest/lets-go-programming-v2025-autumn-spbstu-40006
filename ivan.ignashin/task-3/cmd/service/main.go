package main

import (
	"github.com/IvanIgnashin7D/task-3/internal/parcer"
)

func main() {
	cfg, err := parcer.ParseConfig()
	if err != nil {
		panic(err)
	}

	records, err := parcer.ParseXML(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	err = parcer.SaveAsJSON(records, cfg.OutputFile)
	if err != nil {
		panic(err)
	}
}
