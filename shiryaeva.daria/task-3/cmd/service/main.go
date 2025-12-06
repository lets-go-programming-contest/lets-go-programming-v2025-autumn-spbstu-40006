package main

import (
	"flag"
	"log"

	"github.com/Dora-shi/task-3/internal/config"
	"github.com/Dora-shi/task-3/internal/converter"
)

func main() {
	configPath := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	if *configPath == "" {
		log.Panic("config path is required")
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Panic(err)
	}

	if err := converter.Process(cfg); err != nil {
		log.Panicf("processing failed: %v", err)
	}
}
