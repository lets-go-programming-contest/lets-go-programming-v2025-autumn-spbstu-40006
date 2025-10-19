package main

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrInvalidOperator    = errors.New("invalid operator")
	ErrInvalidTemperature = errors.New("invalid temperature")
)

func checkBorder(operation string, tempStr string) (string, int, error) {
	if operation != "<=" && operation != ">=" {
		return "", 0, ErrInvalidOperator
	}

	tempInt, err := strconv.Atoi(tempStr)
	if err != nil {
		return "", 0, ErrInvalidTemperature
	}

	if tempInt > 30 || tempInt < 15 {
		return "", 0, ErrInvalidTemperature
	}

	return operation, tempInt, nil
}

func processEmployees(employees int) {
	minTemp := 15
	maxTemp := 30
	impossibleCondition := false

	for range employees {
		var (
			oper string
			temp string
		)

		_, err := fmt.Scan(&oper, &temp)
		if err != nil {
			fmt.Println("Error reading operator and temperature:", err)
			fmt.Println(-1)

			impossibleCondition = true

			continue
		}

		operation, temperature, err := checkBorder(oper, temp)
		if err != nil {
			fmt.Println("Error has occurred:", err)
			fmt.Println(-1)

			impossibleCondition = true

			continue
		}

		if impossibleCondition {
			fmt.Println(-1)

			continue
		}

		switch operation {
		case ">=":
			minTemp = max(minTemp, temperature)
		case "<=":
			maxTemp = min(maxTemp, temperature)
		}

		if minTemp > maxTemp {
			fmt.Println(-1)

			impossibleCondition = true

			continue
		}

		fmt.Println(minTemp)
	}
}

func main() {
	var departments int

	_, err := fmt.Scan(&departments)
	if err != nil {
		fmt.Println("Error reading number of departments:", err)

		return
	}

	for range departments {
		var employees int

		_, err := fmt.Scan(&employees)
		if err != nil {
			fmt.Println("Error reading number of employees:", err)

			return
		}

		processEmployees(employees)
	}
}
