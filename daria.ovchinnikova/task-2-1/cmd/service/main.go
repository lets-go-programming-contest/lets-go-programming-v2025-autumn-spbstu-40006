package main

import "fmt"

func main() {
	var (
		numberOfDepatrs, numberOfWorkers int
		minTemp, maxTemp                 int
		sign                             string
		valueTemp                        int
	)

	_, err := fmt.Scan(&numberOfDepatrs)
	if err != nil {
		fmt.Println("Invalid number of departments")

		return
	}

	for range numberOfDepatrs {
		minTemp = 15
		maxTemp = 30

		_, err = fmt.Scan(&numberOfWorkers)
		if err != nil {
			fmt.Println("Invalid number of workers")

			return
		}

		for range numberOfWorkers {
			_, err = fmt.Scan(&sign, &valueTemp)
			if err != nil {
				fmt.Println("Invalid input format")

				return
			}

			if sign == "<=" {
				if valueTemp < maxTemp {
					maxTemp = valueTemp
				}
			} else if sign == ">=" {
				if valueTemp > minTemp {
					minTemp = valueTemp
				}
			}

			if minTemp <= maxTemp {
				fmt.Println(minTemp)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
