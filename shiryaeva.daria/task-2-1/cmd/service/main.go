package main

import (
	"fmt"
)

func main() {
	var (
		numDepartments, numEmployees int
		operator                     string
		temperature                 int
	)

	_, err := fmt.Scan(&numDepartments)
	if err != nil {
		fmt.Println("Invalid input")
		return
	}

	for i := 0; i < numDepartments; i++ {
		minTemp, maxTemp := 15, 30
		possible := true

		_, err := fmt.Scan(&numEmployees)
		if err != nil {
			fmt.Println("Invalid input")
			return
		}
		for j := 0; j < numEmployees; j++ {
			_, err := fmt.Scan(&operator, &temperature)
			if err != nil {
				fmt.Println("Invalid input")
				return
			}

			if operator == ">=" {
				if temperature > minTemp {
					minTemp = temperature
				}
			} else if operator == "<=" {
				if temperature < maxTemp {
					maxTemp = temperature
				}
			}
			if minTemp <= maxTemp && possible {
				fmt.Println(minTemp)
			} else {
				fmt.Println(-1)
				possible = false
			}
		}
	}
}