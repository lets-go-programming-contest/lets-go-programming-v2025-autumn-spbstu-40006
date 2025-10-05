package main

import "fmt"

func main() {
	var (
		firstNumber, secondNumber int
		symbolOperation           string
	)

	_, err := fmt.Scan(&firstNumber)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err = fmt.Scan(&secondNumber)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	_, err = fmt.Scan(&symbolOperation)
	if err != nil {
		fmt.Println("Invalid operation")
		return
	}

	switch symbolOperation {
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
		}
		fmt.Println(firstNumber / secondNumber)
	default:
		fmt.Println("Invalid operation")
	}
}
