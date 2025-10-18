package main

import "fmt"

const (
	baseLower = 15
	baseUpper = 30
)

func readConstraint() (operator string, temperature int, valid bool) {
	var op string
	var temp int

	fieldsRead, err := fmt.Scan(&op, &temp)
	if err != nil || fieldsRead != 2 {
		return "", 0, false
	}
	if op != ">=" && op != "<=" {
		return "", 0, false
	}
	return op, temp, true
}

func applyConstraint(lower, upper int, operator string, temperature int) (int, int) {
	switch operator {
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
		operator, temperature, valid := readConstraint()
		if !valid {
			return
		}
		lowerBound, upperBound = applyConstraint(lowerBound, upperBound, operator, temperature)
		printState(lowerBound, upperBound)
	}
}
