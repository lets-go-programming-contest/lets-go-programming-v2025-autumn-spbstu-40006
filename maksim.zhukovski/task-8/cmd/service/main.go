package main

import (
	"fmt"
	"log"

	"github.com/sp3c7r/task-8/pkg/config"
)

func main() {
	cfg, err := config.Load(config.ConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s %s", cfg.Environment, cfg.LogLevel)
}
