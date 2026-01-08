package main

import (
	"fmt"
	"log"

	"github.com/dizey5k/task-8/config"
)

func main() {
	cfg, err := config.Load(config.ConfigFile)
	if err != nil {
		log.Fatalf("err loading cfg: %v", err)
	}

	fmt.Printf("%s %s", cfg.Environment, cfg.LogLevel)
}
