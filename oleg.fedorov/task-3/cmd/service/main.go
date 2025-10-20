package main

import (
	"flag"

	"github.com/dizey5k/task-3/internal/config"
	"github.com/dizey5k/task-3/internal/converter"
)

func main() {
	configPath := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	if *configPath == "" {
		panic("Config path is required")
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		panic(err)
	}

	converter.Process(cfg)
}
