package main

import "fmt"

func main() {
	var (
		departmentAmount, employeeAmount, temperature int
		curTempLimit, lowerTempLimit, upperTempLimit  int
		str                                           string
		err                                           error
	)

	_, err = fmt.Scan(&departmentAmount)
	if err != nil || departmentAmount < 1 || departmentAmount > 1000 {
		fmt.Println(-1)

		return
	}

	for range departmentAmount {
		temperature = -1
		lowerTempLimit = 15
		upperTempLimit = 30

		_, err = fmt.Scan(&employeeAmount)
		if err != nil || employeeAmount < 1 || employeeAmount > 1000 {
			fmt.Println(-1)

			continue
		}

		for range employeeAmount {
			_, err = fmt.Scan(&str, &curTempLimit)
			if err != nil {
				fmt.Println(-1)

				continue
			}

			switch {
			case str == ">=" && curTempLimit <= upperTempLimit:
				lowerTempLimit = max(lowerTempLimit, curTempLimit)
				temperature = max(lowerTempLimit, temperature)
			case str == "<=" && curTempLimit >= lowerTempLimit:
				upperTempLimit = min(upperTempLimit, curTempLimit)
				temperature = min(upperTempLimit, temperature)
			}

			fmt.Println(temperature)
		}
	}
}
