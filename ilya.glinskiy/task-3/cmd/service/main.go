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

	formattedValutes := make([]io.JSONValute, len(input.Valutes))
	for index, valute := range input.Valutes {
		value, err := strconv.ParseFloat(strings.ReplaceAll(valute.Value, ",", "."), 64)
		if err != nil {
			panic("couldn't read valute property: value")
		}

		formattedValutes[index] = io.JSONValute{
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
