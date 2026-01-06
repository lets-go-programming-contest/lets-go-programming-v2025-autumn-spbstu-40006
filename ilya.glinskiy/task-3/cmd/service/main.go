package main

import (
	"flag"
	"sort"

	"github.com/bloomkicks/task-3/internal/io"
)

func main() {
	var (
		input  io.Input
		config io.Config
		err    error
	)

	configPath := flag.String("config", "", "path to config file")
	flag.Parse()

	if *configPath == "" {
		panic("couldn't find path to config file")
	}

	err = io.ReadConfig(*configPath, &config)
	if err != nil {
		panic(err)
	}

	err = io.ReadInput(config.InputFile, &input)
	if err != nil {
		panic(err)
	}

	sortedValutes := make([]io.JSONValute, len(input.Valutes))
	sort.Slice(sortedValutes, func(i int, j int) bool {
		return input.Valutes[i].Value > input.Valutes[j].Value
	})

	err = io.WriteOutput(config.OutputFile, sortedValutes)
	if err != nil {
		panic(err)
	}
}
