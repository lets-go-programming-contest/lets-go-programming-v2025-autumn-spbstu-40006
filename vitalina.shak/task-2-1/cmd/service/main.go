package main

import (
	"errors"
	"fmt"
)

const (
	initMinTemp = 15
	initMaxTemp = 30
)

var (
	errReadDepartments = errors.New("failed to read departments count")
	errReadEmployees   = errors.New("failed to read employees count")
	errReadRequirement = errors.New("failed to read requirement")
)

func main() {
	var departmentsCount int
	if _, err := fmt.Scan(&departmentsCount); err != nil {
		fmt.Println(errReadDepartments.Error())

		return
	}

	for range departmentsCount {
		if err := processDepartment(); err != nil {
			fmt.Println(err.Error())

			return
		}
	}
}

func processDepartment() error {
	var employeesCount int
	if _, err := fmt.Scan(&employeesCount); err != nil {
		return errReadEmployees
	}

	currentMinTemp := initMinTemp
	currentMaxTemp := initMaxTemp

	for range employeesCount {
		var (
			operation string
			temp      int
		)

		if _, err := fmt.Scan(&operation, &temp); err != nil {
			return errReadRequirement
		}

		switch operation {
		case ">=":
			if temp > currentMinTemp {
				currentMinTemp = temp
			}
		case "<=":
			if temp < currentMaxTemp {
				currentMaxTemp = temp
			}
		default:
			fmt.Println("Invalid operation")

			continue
		}

		if currentMinTemp <= currentMaxTemp {
			fmt.Println(currentMinTemp)

			continue
		}

		fmt.Println(-1)
	}

	return nil
}
