package main

import "fmt"

func main() {
	var (
		departNum, emplsNum, minTemp, maxTemp, temp int
		operand                                     string
	)

	_, err := fmt.Scan(&departNum)
	if err != nil {
		fmt.Println("Incorrect number of departments")

		return
	}

	for range departNum {
		_, err = fmt.Scan(&emplsNum)

		minTemp = 15
		maxTemp = 30

		if err != nil {
			fmt.Println("Incorrect number of employees")

			return
		}

		for range emplsNum {
			_, err = fmt.Scan(&operand, &temp)
			if err != nil {
				fmt.Println("Incorrect number if temperature")

				return
			}

			if operand[0] == '>' && temp > minTemp {
				minTemp = temp
			} else if operand[0] == '<' && temp < maxTemp {
				maxTemp = temp
			}

			if minTemp <= maxTemp {
				fmt.Println(minTemp)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
