// package main

package main

import (
	"fmt"
)

func setTemperature(parameter string, temperature int, maxTemp, minTemp *int) int {
	switch parameter {
	case ">=":
		if temperature > *minTemp {
			*minTemp = temperature
		}
	case "<=":
		if temperature < *maxTemp {
			*maxTemp = temperature
		}
	default:
		return -1
	}

	if *minTemp > *maxTemp {
		return -1
	}

	return *minTemp
}

func main() {
	var (
		numDep, dearColleagues, temperature int
		parameter                           string
	)

	if _, err := fmt.Scan(&numDep); err != nil {
		fmt.Println(-1)

		return
	}

	for range numDep {
		maxTemp, minTemp := 30, 15

		if _, err := fmt.Scan(&dearColleagues); err != nil {
			fmt.Println(-1)

			return
		}

		for range dearColleagues {
			if _, err := fmt.Scan(&parameter, &temperature); err != nil {
				fmt.Println(-1)

				return
			}

			out := setTemperature(parameter, temperature, &maxTemp, &minTemp)

			fmt.Println(out)
		}
	}
}
