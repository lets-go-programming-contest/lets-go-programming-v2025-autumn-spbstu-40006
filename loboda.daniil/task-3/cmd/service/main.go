package main

import (
	"flag"
	"fmt"

	"github.com/Daniil-drom/task-3/internal/cfg"
	"github.com/Daniil-drom/task-3/internal/rates"
)

func parseFlags() string {
	var cfgPath string

	flag.StringVar(&cfgPath, "config", "", "path to YAML config")
	flag.Parse()

	return cfgPath
}

func main() {
	cfgPath := parseFlags()

	if cfgPath == "" {
		panic("config path is required (use -config)")
	}

	conf, err := cfg.Load(cfgPath)
	if err != nil {
		panic(fmt.Errorf("config: %w", err))
	}

	cur, err := rates.Load(conf.InputFile)
	if err != nil {
		panic(fmt.Errorf("xml: %w", err))
	}

	if err := rates.SaveJSON(conf.OutputFile, cur); err != nil {
		panic(fmt.Errorf("json: %w", err))
	}

	fmt.Printf("OK: %d items -> %s\n", len(cur), conf.OutputFile)
}
