package main

import (
	"fmt"
	"log"

	"github.com/HorekProgrammer/task-8/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(cfg.Environment, " ", cfg.LogLevel)
}
