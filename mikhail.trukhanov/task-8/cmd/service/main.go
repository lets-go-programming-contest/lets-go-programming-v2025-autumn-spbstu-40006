package main

import (
	"fmt"
	"log"

	"github.com/Mishaa105/task-8/config"
)

func main() {
	cfg, err := config.Load(config.ConfigFile)
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	fmt.Printf("%s %s", cfg.Environment, cfg.LogLevel)
}
