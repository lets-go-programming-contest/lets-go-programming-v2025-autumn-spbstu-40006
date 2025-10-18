package main

import "fmt"

const (
	baseLower = 15
	baseUpper = 30
)

func applyConstraint(lower, upper int, operatorToken string, temperature int) (int, int) {
	if operatorToken == ">=" {
		if temperature > lower {
			lower = temperature
		}

		return lower, upper
	}

	if operatorToken == "<=" {
		if temperature < upper {
			upper = temperature
		}

		return lower, upper
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
	var departmentsCount, employeesCount int
	if _, err := fmt.Scan(&departmentsCount, &employeesCount); err != nil {
		return
	}

	for range departmentsCount {
		lowerBound, upperBound := baseLower, baseUpper

		for range employeesCount {
			var operatorToken string
			var temperature int
			if _, err := fmt.Scan(&operatorToken, &temperature); err != nil {
				return
			}

			lowerBound, upperBound = applyConstraint(lowerBound, upperBound, operatorToken, temperature)
			printState(lowerBound, upperBound)
		}
	}
}
