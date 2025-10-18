package main

import (
	"fmt"
	"strconv"
)

const (
	baseLower = 15
	baseUpper = 30
)

func readConstraint() (string, int, bool) {
	var operatorToken string
	var temperature int

	fieldsRead, err := fmt.Scan(&operatorToken, &temperature)
	if err != nil || fieldsRead != 2 {
		return "", 0, false
	}
	if operatorToken != ">=" && operatorToken != "<=" {
		return "", 0, false
	}
	return operatorToken, temperature, true
}

func readCounts() (int, int, bool, string, int, bool) {
	var firstNumber int
	if read, err := fmt.Scan(&firstNumber); err != nil || read != 1 || firstNumber < 0 {
		return 0, 0, false, "", 0, false
	}

	var secondToken string
	readSecond, _ := fmt.Scan(&secondToken)
	if readSecond != 1 {
		return 1, firstNumber, false, "", 0, true
	}
	if numeric, err := strconv.Atoi(secondToken); err == nil {
		return firstNumber, numeric, false, "", 0, true
	}

	var temperature int
	if read, err := fmt.Scan(&temperature); err != nil || read != 1 {
		return 0, 0, false, "", 0, false
	}
	return 1, firstNumber, true, secondToken, temperature, true
}

func applyConstraint(lower, upper int, operatorToken string, temperature int) (int, int) {
	switch operatorToken {
	case ">=":
		if temperature > lower {
			lower = temperature
		}
	case "<=":
		if temperature < upper {
			upper = temperature
		}
	}
	return lower, upper
}

func printState(lower, upper int) {
	if lower <= upper {
		fmt.Println(lower)

		return
	}
	fmt.Println(-1)
}

func main() {
	departments, limitsPerDepartment, haveFirst, firstOperator, firstTemperature, ok := readCounts()
	if !ok {
		return
	}

	for i := 0; i < departments; i++ {
		lowerBound, upperBound := baseLower, baseUpper

		for j := 0; j < limitsPerDepartment; j++ {
			if haveFirst && j == 0 {
				lowerBound, upperBound = applyConstraint(lowerBound, upperBound, firstOperator, firstTemperature)
				printState(lowerBound, upperBound)
				continue
			}

			operatorToken, temperature, valid := readConstraint()
			if !valid {
				return
			}

			lowerBound, upperBound = applyConstraint(lowerBound, upperBound, operatorToken, temperature)
			printState(lowerBound, upperBound)
		}
		haveFirst = false
	}
}
