package main

import (
	"fmt"

	"github.com/MoneyprogerISG/task-8/internal/config"
)

func main() {
	cfg, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s %s\n", cfg.Environment, cfg.LogLevel)
}
