package main

import "fmt"

func main() {

	var first int
	var second int
	var operator string

	if _, err := fmt.Scanln(&first); err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	if _, err := fmt.Scanln(&second); err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	if _, err := fmt.Scanln(&operator); err != nil {
		fmt.Println("Invalid operation")
		return
	}

	switch operator {

	case "+":
		fmt.Println(first + second)
	case "-":
		fmt.Println(first - second)
	case "*":
		fmt.Println(first * second)
	case "/":

		if second == 0 {
			fmt.Println("Division by zero")
			return
		}

		fmt.Println(first / second)

	default:
		fmt.Println("Invalid operation")
	}
}
