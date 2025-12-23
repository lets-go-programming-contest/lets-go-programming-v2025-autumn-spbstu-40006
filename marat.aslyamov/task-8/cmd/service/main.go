package main

import (
	"fmt"
	"log"

	"github.com/tuesdayy1/task-8/config"
)

func main() {
	cfg, err := config.Load(config.ConfigFile)
	if err != nil {
		log.Fatalf("Config download error: %v", err)
	}

	fmt.Printf("%s %s", cfg.Environment, cfg.LogLevel)
}