package main

import (
	"fmt"
	"github.com/identicalaffiliation/task-8/config"
	"log"
)

func main() {
	cfg, err := config.Load(config.ConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s %s", cfg.Environment, cfg.LogLevel)
}
