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
		minTemp, maxTemp                         int = 15, 30
		countDepartment, countStaff, temperature int
		operation                                string
	)

	mustScan(&countDepartment)

	for i := 0; i < countDepartment; i++ {
		mustScan(&countStaff)
		minTemp, maxTemp = 15, 30

		for j := 0; j < countStaff; j++ {
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
