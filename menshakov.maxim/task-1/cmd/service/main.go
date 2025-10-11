package main

import "fmt"

func main() {
	var num1, num2 int
	var operator string

	if _, err := fmt.Scan(&num1); err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	if _, err := fmt.Scan(&num2); err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	if _, err := fmt.Scan(&operator); err != nil || (operator != "+" && operator != "-" && operator != "*" && operator != "/") {
		fmt.Println("Invalid operation")
		return
	}

	switch operator {
	case "+":
		fmt.Println(num1 + num2)
	case "-":
		fmt.Println(num1 - num2)
	case "*":
		fmt.Println(num1 * num2)
	case "/":
		if num2 == 0 {
			fmt.Println("Division by zero")
		} else {
			fmt.Println(num1 / num2)
		}
	}
}
