package main

import (
	"flag"
	"log"

	"github.com/dizey5k/task-3/internal/config"
	"github.com/dizey5k/task-3/internal/converter"
)

func main() {
	configPath := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	if *configPath == "" {
		log.Panic("Config path is required")
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Panicf("Failed to load config: %v", err)
	}

	if err := converter.Process(cfg); err != nil {
		log.Panicf("Processing failed: %v", err)
	}
}
