// package main

package main

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidNum       = errors.New("не является числом")
	ErrInvalidParameter = errors.New("не является параметром")
)

func setTemperature(parameter string, temperature int, maxTemp, minTemp *int) (int, error) {
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
		return -1, ErrInvalidParameter
	}

	if *minTemp > *maxTemp {
		return -1, ErrInvalidParameter
	}

	return *minTemp, nil
}

func main() {
	var (
		numDep, dearColleagues, temperature int
		parameter                           string
	)

	if _, err := fmt.Scan(&numDep); err != nil {
		fmt.Println(ErrInvalidNum)

		return
	}

	for range numDep {
		maxTemp, minTemp := 30, 15

		if _, err := fmt.Scan(&dearColleagues); err != nil {
			fmt.Println(ErrInvalidNum)

			return
		}

		for range dearColleagues {
			if _, err := fmt.Scan(&parameter, &temperature); err != nil {
				fmt.Println(ErrInvalidNum)

				return
			}

			out, err := setTemperature(parameter, temperature, &maxTemp, &minTemp)
			if err != nil {
				fmt.Println(ErrInvalidParameter)

				return
			}

			fmt.Println(out)
		}
	}
}
