package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Segfault-chan/task-3/internal/utils"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Invalid number of args: %v\nThe correct usage is: -config <path-to-config>\n", os.Args)

		return
	}

	if os.Args[1] != "-config" {
		fmt.Fprint(os.Stderr, "Invalid first operand.\nThe correct usage is: -config <path-to-config>\n", os.Args)

		return
	}

	configFile, err := utils.ParseYAML(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	exchangeRates, err := utils.ParseXML(configFile.InputFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(exchangeRates)
}
