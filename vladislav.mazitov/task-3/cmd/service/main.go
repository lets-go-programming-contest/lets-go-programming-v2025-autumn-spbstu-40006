package main

import (
	"flag"

	"github.com/identicalaffiliation/task-3/internal/config"
	"github.com/identicalaffiliation/task-3/internal/parser"
)

func main() {
	confgiPath := flag.String("config", "config.yaml", "path to config")
	flag.Parse()

	cfg, err := config.Load(*confgiPath)
	if err != nil {
		panic(err)
	}

	err = parser.Process(cfg)
	if err != nil {
		panic(err)
	}
}
