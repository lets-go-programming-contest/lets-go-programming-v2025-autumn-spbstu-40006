package main

import "fmt"

func main() {
	var numDepartments int

	_, err := fmt.Scan(&numDepartments)
	if err != nil {
		fmt.Println("Invalid number of departments")

		return
	}

	for range numDepartments {
		minT := 15
		maxT := 30

		var numWorkers int

		_, err = fmt.Scan(&numWorkers)
		if err != nil {
			fmt.Println("Invalid number of workers")

			return
		}

		for range numWorkers {
			var operation string

			var Temp int

			_, err = fmt.Scan(&operation, &Temp)
			if err != nil {
				fmt.Println("Invalid operation")

				return
			}

			switch operation {
			case ">=":
				if Temp > minT {
					minT = Temp
				}
			case "<=":
				if Temp < maxT {
					maxT = Temp
				}
			default:
				fmt.Println("Invalid operation")

				return
			}

			if minT <= maxT {
				fmt.Println(minT)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
