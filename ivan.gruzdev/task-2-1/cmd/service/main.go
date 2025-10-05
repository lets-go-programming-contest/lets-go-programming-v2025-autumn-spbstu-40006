package main

import (
	"fmt"
	"os"
)

func mustScan(a ...interface{}) {
	_, err := fmt.Scan(a...)
	if err != nil {
		os.Exit(1)
	}
}

func main() {
	var (
		countDepartment, countStaff, temperature int
		operation                                string
	)

	mustScan(&countDepartment)

	for range countDepartment {
		mustScan(&countStaff)

		minTemp, maxTemp := 15, 30

		for range countStaff {
			mustScan(&operation, &temperature)

			if operation == ">=" {
				if temperature > minTemp {
					minTemp = temperature
				}
			} else {
				if temperature < maxTemp {
					maxTemp = temperature
				}
			}

			if minTemp > maxTemp {
				fmt.Println(-1)
			} else {
				fmt.Println(minTemp)
			}
		}
	}
}
