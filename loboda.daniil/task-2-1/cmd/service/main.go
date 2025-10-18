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
			var (
				operatorToken string
				temperature   int
			)

			if _, err := fmt.Scan(&operatorToken, &temperature); err != nil {
				return
			}

			switch operatorToken {
			case ">=":
				if temperature > lowerBound {
					lowerBound = temperature
				}
			case "<=":
				if temperature < upperBound {
					upperBound = temperature
				}
			default:
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
