package main

import (
	"fmt"

	"github.com/abdelrhmanbaha/task-8/pkg/config"
)

func main() {
	cfg := config.Get()
	fmt.Printf("Environment: %s, Log Level: %s\n", cfg.Environment, cfg.LogLevel)
}
