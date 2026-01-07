package main

import (
	"fmt"

	"github.com/vitsh1/task-8/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Println("failed to load config: %w", err)

		return
	}

	fmt.Println(cfg.Environment, cfg.LogLevel)
}
