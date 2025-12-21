package main

import (
	"fmt"
	"log"

	"github.com/identicalaffiliation/task-8/config"
)

func main() {
	cfg, err := config.Load(config.ConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s %s", cfg.Environment, cfg.LogLevel)
}
