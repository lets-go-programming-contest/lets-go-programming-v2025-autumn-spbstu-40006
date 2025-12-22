package main

import (
	"fmt"
	"log"

	"github.com/dizey5k/task-8/package/config"
)

func main() {
	loader := config.NewLoader()

	cfg, err := config.GetConfig(loader)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	fmt.Printf("Environment: %s, Log Level: %s\n",
		cfg.Environment,
		cfg.LogLevel)

	displayConfigInfo(cfg)
}

func displayConfigInfo(cfg *config.Config) {
	log.Printf("Application running in %s mode", cfg.Environment)

	switch cfg.LogLevel {
	case "debug":
		log.Println("Debug logging enabled - detailed information will be shown")
	case "error":
		log.Println("Only error logging enabled - minimal output")
	default:
		log.Printf("Log level set to: %s", cfg.LogLevel)
	}
}
