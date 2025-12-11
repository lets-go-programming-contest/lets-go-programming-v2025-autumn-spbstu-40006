package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/abdelrhmanbaha/task-3/internal/config"
	"github.com/abdelrhmanbaha/task-3/internal/currencies"
)

func main() {
	configPath := flag.String("config", "", "Path to config file")
	flag.Parse()

	if *configPath == "" {
		log.Fatalf("failed to load config: path not provided")
	}

	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	data, err := os.ReadFile(cfg.InputFile)
	if err != nil {
		log.Fatalf("failed to read xml %v", err)
	}

	svc := currencies.NewCurrencyService()

	list, err := svc.ParseXML(data)
	if err != nil {
		log.Fatalf("failed to parse xml %v", err)
	}

	svc.SortByValue(list)

	if err := svc.SaveToJSON(cfg.OutputFile, list); err != nil {
		log.Fatalf("failed to save json %v", err)
	}

	fmt.Printf("done output %s\n", cfg.OutputFile)
}