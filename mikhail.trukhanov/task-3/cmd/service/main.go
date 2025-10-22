package main

import (
	"fmt"

	"github.com/Mishaa105/task-3/internal/decoding"
	"github.com/Mishaa105/task-3/internal/saving"
)

func main() {
	valCurs := decoding.Decoding("config.yaml")

	for _, val := range valCurs.Valutes {
		fmt.Printf("NumCode: %d, CharCode: %s, Value: %.2f\n", val.NumCode, val.CharCode, val.Value)
	}

	saving.SaveToJSON("output/result.json", valCurs.Valutes)
}
