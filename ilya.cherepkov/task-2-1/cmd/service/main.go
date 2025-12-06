package main

import "fmt"

func main() {
	var (
		departmentsCount, employeesCount, value int
		operator                                string
	)

	_, err := fmt.Scan(&departmentsCount)
	if err != nil {
		fmt.Println("Invalid input")

		return
	}

	for department := 0; department != departmentsCount; department++ {
		minimum, maximum := 15, 30

		_, err := fmt.Scan(&employeesCount)
		if err != nil {
			fmt.Println("Invalid input")

			return
		}

		for employee := 0; employee != employeesCount; employee++ {
			_, err := fmt.Scan(&operator, &value)
			if err != nil {
				fmt.Println("Invalid input")

				return
			}

			if operator == ">=" && value > minimum {
				minimum = value
			}

			if operator == "<=" && value < maximum {
				maximum = value
			}

			if minimum <= maximum {
				fmt.Println(minimum)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
