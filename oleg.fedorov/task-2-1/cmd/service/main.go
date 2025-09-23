package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type TemperatureRange struct {
	min int
	max int
}

func NewTemperatureRange() *TemperatureRange {
	return &TemperatureRange{
		min: 15,
		max: 30,
	}
}

func (tr *TemperatureRange) ApplyOperation(operation string, temperature int) string {
	switch operation {
	case ">=":
		if temperature > tr.min {
			tr.min = temperature
		}
	case "<=":
		if temperature < tr.max {
			tr.max = temperature
		}
	}

	if tr.min <= tr.max {
		return fmt.Sprintf("%d", tr.min)
	}
	return "-1"
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	depCount, _ := strconv.Atoi(scanner.Text())

	for i := 0; i < depCount; i++ {
		scanner.Scan()
		workersNumber, _ := strconv.Atoi(scanner.Text())

		tempRange := NewTemperatureRange()

		for j := 0; j < workersNumber; j++ {
			scanner.Scan()
			line := scanner.Text()
			parts := strings.Fields(line)
			operation := parts[0]
			temp, _ := strconv.Atoi(parts[1])

			result := tempRange.ApplyOperation(operation, temp)
			fmt.Println(result)
		}
	}
}
