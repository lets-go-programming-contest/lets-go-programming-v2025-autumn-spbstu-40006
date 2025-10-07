package main

import (
	"log"
	"os"
	"slices"

	"github.com/Segfault-chan/task-3/internal/utils"
)

func main() {
	if len(os.Args) != 3 {
		log.Panicf("Invalid number of args: %v\nThe correct usage is: -config <path-to-config>\n", os.Args)
		return
	}

	if os.Args[1] != "-config" {
		log.Panicf("Invalid operand: %v\nThe correct usage is: -config <path-to-config>\n", os.Args)
		return
	}

	configFile, err := utils.ParseYAML(os.Args[2])
	if err != nil {
		log.Panic(err)
	}

	exchangeRates, err := utils.ParseXML(configFile.InputFile)
	if err != nil {
		log.Panic(err)
	}

	slices.SortFunc(exchangeRates.Currencies, utils.DescendingComparatorCurrency)
	if err := utils.ParseJSON(*exchangeRates, configFile.OutputFile); err != nil {
		log.Panic(err)
	}
}
