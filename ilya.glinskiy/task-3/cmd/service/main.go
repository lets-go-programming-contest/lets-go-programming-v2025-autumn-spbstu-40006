package main

import (
	"flag"
	"sort"
	"strconv"
	"strings"

	"github.com/bloomkicks/task-3/internal/io"
)

func main() {
	var (
		input  io.Input
		config io.Config
		err    error
	)

	configPath := flag.String("config", "", "config path")
	flag.Parse()

	if *configPath == "" {
		panic("Config must be provided in arguments")
	}

	err = io.ReadConfig(*configPath, &config)
	if err != nil {
		panic(err)
	}

	err = io.ReadInput(config.InputFile, &input)
	if err != nil {
		panic(err)
	}

	formattedValutes := make([]io.JSONValute, len(input.Valutes))
	for i, valute := range input.Valutes {
		value, err := strconv.ParseFloat(strings.ReplaceAll(valute.Value, ",", "."), 64)
		if err != nil {
			panic("Couldn't read Valute property: Value")
		}

		formattedValutes[i] = io.JSONValute{
			Value:    value,
			NumCode:  valute.NumCode,
			CharCode: valute.CharCode,
		}
	}

	sort.Slice(formattedValutes, func(i int, j int) bool {
		return formattedValutes[i].Value > formattedValutes[j].Value
	})

	err = io.WriteOutput(config.OutputFile, formattedValutes)
	if err != nil {
		panic(err)
	}
}
