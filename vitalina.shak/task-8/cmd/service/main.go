package main

import (
	"fmt"

	"github.com/vitsh1/task-8/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("failed to load config: %v", err)

		return
	}

	fmt.Printf("%s %s", cfg.Environment, cfg.LogLevel)
}
