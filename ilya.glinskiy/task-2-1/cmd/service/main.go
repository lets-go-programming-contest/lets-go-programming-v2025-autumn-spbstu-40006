package main

import "fmt"

func main() {
	var (
		departmentAmount, employeeAmount       int
		suggestedLimit, lowerLimit, upperLimit int
		str                                    string
		err                                    error
	)

	_, err = fmt.Scan(&departmentAmount)
	if err != nil {
		fmt.Println(-1)

		return
	}

	for range departmentAmount {
		lowerLimit = 15
		upperLimit = 30

		_, err = fmt.Scan(&employeeAmount)
		if err != nil {
			fmt.Println(-1)

			continue
		}

		for range employeeAmount {
			_, err = fmt.Scan(&str, &suggestedLimit)
			if err != nil {
				fmt.Println(-1)

				continue
			}

			if str == ">=" && suggestedLimit > lowerLimit {
				lowerLimit = suggestedLimit
			} else if str == "<=" && suggestedLimit < upperLimit {
				upperLimit = suggestedLimit
			}

			if lowerLimit <= upperLimit {
				fmt.Println(lowerLimit)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
