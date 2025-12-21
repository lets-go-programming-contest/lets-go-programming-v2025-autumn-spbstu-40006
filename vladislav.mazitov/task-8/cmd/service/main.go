package main

import (
	"github.com/identicalaffiliation/task-8/config"
	"fmt"
	"log"
)

func main() {
	cfg, err := config.Load(config.ConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s %s", cfg.Environment, cfg.LogLevel)
}
