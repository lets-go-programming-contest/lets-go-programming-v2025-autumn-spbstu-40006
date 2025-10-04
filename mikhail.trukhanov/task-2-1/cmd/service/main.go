package main

import "fmt"

func main() {
	var (
		numOfDepartments, numOfWorkers, temp int
		sign                                 string
	)

	_, err := fmt.Scan(&numOfDepartments)
	if err != nil {
		fmt.Println("Invalid input")

		return
	}

	for i := 0; i != numOfDepartments; i++ {
		minTemp, maxTemp := 15, 30

		_, err := fmt.Scan(&numOfWorkers)
		if err != nil {
			fmt.Println("Invalid input")

			return
		}

		for j := 0; j != numOfWorkers; j++ {
			_, err := fmt.Scan(&sign, &temp)
			if err != nil {
				fmt.Println("Invalid input")

				return
			}

			if sign == ">=" && temp > minTemp {
				minTemp = temp
			}

			if sign == "<=" && temp < maxTemp {
				maxTemp = temp
			}

			if maxTemp >= minTemp {
				fmt.Println(minTemp)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
