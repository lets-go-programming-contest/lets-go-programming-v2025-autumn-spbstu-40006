package main

import "fmt"

func main() {
	var num1, num2 int
	var operator string

	_, err := fmt.Scan(&num1)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err = fmt.Scan(&num2)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	fmt.Scan(&operator)

	switch operator {
	case "+":
		fmt.Println(num1 + num2)
	case "-":
		fmt.Println(num1 - num2)
	case "*":
		fmt.Println(num1 * num2)
	case "/":
		if num2 != 0 {
			fmt.Println(num1 / num2)
		} else {
			fmt.Println("Division by zero")
		}
	default:
		fmt.Println("Invalid operation")
	}
}
