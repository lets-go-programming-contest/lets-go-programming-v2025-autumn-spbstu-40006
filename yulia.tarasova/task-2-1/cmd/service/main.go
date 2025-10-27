package main

import "fmt"

func main() {
	var (
		numDepartments, numWorkers, temp int
		operation                        string
	)

	_, err := fmt.Scanln(&numDepartments)
	if err != nil {
		fmt.Println("Invalid input")

		return
	}

	for range numDepartments {
		minTemp := 15
		maxTemp := 30

		_, err = fmt.Scanln(&numWorkers)
		if err != nil {
			fmt.Println("Invalid input")

			return
		}

		for range numWorkers {
			_, err = fmt.Scanln(&operation, &temp)
			if err != nil {
				fmt.Println("Invalid input")

				return
			}

			if operation == ">=" && temp > minTemp {
				minTemp = temp
			}

			if operation == "<=" && temp < maxTemp {
				maxTemp = temp
			}

			if minTemp > maxTemp {
				fmt.Println("-1")
			} else {
				fmt.Println(minTemp)
			}
		}
	}
}
