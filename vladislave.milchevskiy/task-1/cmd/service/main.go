package main

import (
	"fmt"
)

func main() {
	var (
		operand1Str, operand2Str int
		operation                string
	)

	operand1, err1 := fmt.Scan(&operand1Str)
	if err1 != nil {
		fmt.Println("Invalid first operand")
		return
	}

	operand2, err2 := fmt.Scan(&operand2Str)
	if err2 != nil {
		fmt.Println("Invalid second operand")
		return
	}

	var result int
	switch operation {
	case "+":
		result = operand1 + operand2
	case "-":
		result = operand1 - operand2
	case "*":
		result = operand1 * operand2
	case "/":
		if operand2 == 0 {
			fmt.Println("Division by zero")
			return
		}
		result = operand1 / operand2
	default:
		fmt.Println("Invalid operation")
		return
	}

	fmt.Println(result)
}

