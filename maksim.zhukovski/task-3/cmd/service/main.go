package main

import (
	"github.com/sp3c7r/task-3/internal/exporter"
	"github.com/sp3c7r/task-3/internal/parser/config"
	"github.com/sp3c7r/task-3/internal/parser/xml"
)

func main() {
	configPath := config.ReadConfigPath()

	cnf := *(config.ParseConfig(configPath))

	vc := xml.ParseXML(cnf.InputFile)

	exporter.WriteToJSON(vc.Valute, cnf.OutputFile)
}
