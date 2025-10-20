package main

import (
	"flag"
	"os"

	"github.com/IvanIgnashin7D/task-3/internal/XMLJSONParce"
	"gopkg.in/yaml.v3"
)

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func main() {
	configPath := flag.String("config", "config.yaml", "Path to the configuration file")
	flag.Parse()

	data, err := os.ReadFile(*configPath)
	if err != nil {
		panic(err)
	}

	var cfg Config

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		panic(err)
	}

	records, err := XMLJSONParce.ParseXML(cfg.InputFile)
	if err != nil {
		panic(err)
	}

	err = XMLJSONParce.SaveAsJSON(records, cfg.OutputFile)
	if err != nil {
		panic(err)
	}
}
