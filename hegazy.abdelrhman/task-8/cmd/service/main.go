package main

import (
	"fmt"
	"log"

	"github.com/abdelrhmanbaha/task-8/package/config"
)

func main() {
	cfg, err := config.Load(config.ConfigFile)
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	fmt.Printf("%s %s", cfg.Environment, cfg.LogLevel)
}