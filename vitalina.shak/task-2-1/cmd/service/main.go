package main

import (
	"errors"
	"fmt"
)

const (
	initMinTemp = 15
	initMaxTemp = 30
)

func main() {
	var departmentsCount int
	if _, err := fmt.Scan(&departmentsCount); err != nil {
		fmt.Println("Failed to read departments count")
		return
	}

	for departmentIndex := 0; departmentIndex < departmentsCount; departmentIndex++ {
		if err := processDepartment(); err != nil {
			fmt.Println(err.Error())
			return
		}
	}
}

func processDepartment() error {
	var employeesCount int
	if _, err := fmt.Scan(&employeesCount); err != nil {
		return errors.New("Failed to read employees count")
	}

	currentMinTemp := initMinTemp
	currentMaxTemp := initMaxTemp

	for employeeIndex := 0; employeeIndex < employeesCount; employeeIndex++ {
		var operation string
		var temp int

		if _, err := fmt.Scan(&operation, &temp); err != nil {
			return errors.New("Failed to read requirement")
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
		} else {
			fmt.Println(-1)
		}
	}

	return nil
}
