package main

import (
	"fmt"
	"log"

	"loboda.daniil/task-8/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s %s\n", cfg.Environment, cfg.LogLevel)
}
