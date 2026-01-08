package main

import (
	"fmt"
	"log"

	"github.com/Daniil-drom/task-8/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("%s %s", cfg.Environment, cfg.LogLevel)
}
