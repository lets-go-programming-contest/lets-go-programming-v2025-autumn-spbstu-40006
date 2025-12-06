package main

import "fmt"

func main() {
	var (
		firstNumber, secondNumber int
		operation                 string
	)
	_, err := fmt.Scanln(&firstNumber)
	if err != nil {
		fmt.Println("Invalid first operand")

		return
	}

	_, err = fmt.Scanln(&secondNumber)
	if err != nil {
		fmt.Println("Invalid second operand")

		return
	}

	_, err = fmt.Scanln(&operation)
	if err != nil {
		fmt.Println("Invalid operation")

		return
	}

	switch operation {
	case "+":
		fmt.Println(firstNumber + secondNumber)
	case "-":
		fmt.Println(firstNumber - secondNumber)
	case "*":
		fmt.Println(firstNumber * secondNumber)
	case "/":
		if secondNumber == 0 {
			fmt.Println("Division by zero")

			return
		} else {
			fmt.Println(firstNumber / secondNumber)
		}
	default:
		fmt.Println("Invalid operation")

		return
	}
}
