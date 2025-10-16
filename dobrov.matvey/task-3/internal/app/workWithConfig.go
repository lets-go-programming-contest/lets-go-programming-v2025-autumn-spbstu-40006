package main

import "flag"

type Config struct {
	InputFile  string `yaml:"input-file"`
	OutputFile string `yaml:"output-file"`
}

func getConfigPath() string {
	configPathPtr := flag.String("config", "", "config.yaml")
	flag.Parse()

	return *configPathPtr
}
