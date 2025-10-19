package main

import "fmt"

func main() {
	var (
		numDepartments, numEmployees, preferredTemp int
		operation                                   string
	)

	_, err := fmt.Scan(&numDepartments)
	if err != nil {
		fmt.Println("Incorrect number of departments")

		return
	}

	for range numDepartments {
		_, err = fmt.Scan(&numEmployees)
		if err != nil {
			fmt.Println("Incorrect number of employees")

			return
		}

		minTemperature := 15
		maxTemperature := 30

		for range numEmployees {
			_, err = fmt.Scan(&operation, &preferredTemp)
			if err != nil {
				fmt.Println("Incorrect temperature value")

				return
			}

			if operation == ">=" && preferredTemp > minTemperature {
				minTemperature = preferredTemp
			} else if operation == "<=" && preferredTemp < maxTemperature {
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
