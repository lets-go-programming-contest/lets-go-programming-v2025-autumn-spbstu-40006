package main

import "fmt"

const (
	baseLower = 15
	baseUpper = 30
)

func readConstraint() (string, int, bool) {
	var (
		operatorToken string
		temperature   int
	)

	fieldsRead, err := fmt.Scan(&operatorToken, &temperature)
	if err != nil || fieldsRead != 2 {

		return "", 0, false
	}

	if operatorToken != ">=" && operatorToken != "<=" {

		return "", 0, false
	}

	return operatorToken, temperature, true
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
	} else {
		fmt.Println(-1)
	}
}

func main() {
	var limitsCount int
	if read, err := fmt.Scan(&limitsCount); err != nil || read != 1 || limitsCount < 0 {

		return
	}

	lowerBound, upperBound := baseLower, baseUpper

	for range limitsCount {
		operatorToken, temperature, valid := readConstraint()
		if !valid {

			return
		}

		lowerBound, upperBound = applyConstraint(lowerBound, upperBound, operatorToken, temperature)

		printState(lowerBound, upperBound)
	}
}
