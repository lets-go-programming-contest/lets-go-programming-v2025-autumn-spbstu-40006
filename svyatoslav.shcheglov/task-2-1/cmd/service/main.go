package main

import "fmt"

func main() {
	var (
		numDepartments, numEmployees, preferredTemp int
		operation                                   string
	)
	_, err := fmt.Scan(&numDepartments)
	if err != nil || numDepartments > 1000 {
		fmt.Println("Incorrect number of departments")
		return
	}

	for i := 0; i < numDepartments; i++ {
		_, err = fmt.Scan(&numEmployees)
		if err != nil || numEmployees > 1000 {
			fmt.Println("Incorrect number of employees")
			return
		}

		minTemperature := 15
		maxTemperature := 30

		for j := 0; j < numEmployees; j++ {
			_, err = fmt.Scan(&operation, &preferredTemp)
			if err != nil {
				fmt.Println("Incorrect number if temperature")
				return
			}

			if operation[0] == '>' && preferredTemp > minTemperature {
				minTemperature = preferredTemp
			} else if operation[0] == '<' && preferredTemp < maxTemperature {
				maxTemperature = preferredTemp
			}

			if minTemperature <= maxTemperature {
				fmt.Println(minTemperature)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
