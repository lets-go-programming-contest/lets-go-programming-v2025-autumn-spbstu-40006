package main

import "fmt"

const (
	baseLower = 15
	baseUpper = 30
)

func main() {
	var departmentsCount int
	if _, err := fmt.Scan(&departmentsCount); err != nil {
		return
	}

	for range departmentsCount {
		var workersCount int
		if _, err := fmt.Scan(&workersCount); err != nil {
			return
		}

		lowerBound := baseLower
		upperBound := baseUpper

		for range workersCount {
			var operator string
			var temperature int
			if _, err := fmt.Scan(&operator, &temperature); err != nil {
				return
			}

			if operator == ">=" {
				if temperature > lowerBound {
					lowerBound = temperature
				}
			} else if operator == "<=" {
				if temperature < upperBound {
					upperBound = temperature
				}
			} else {
				return
			}

			if lowerBound <= upperBound {
				fmt.Println(lowerBound)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
