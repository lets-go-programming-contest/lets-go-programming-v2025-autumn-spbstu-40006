package main

import (
	"fmt"
	"log"
)

func main() {
	var num1, num2 int
	var op string

	_, err := fmt.Scan(&num1)
	if err != nil {
		log.Println("Invalid first operand")
		return
	}

	_, err = fmt.Scan(&num2)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	_, err = fmt.Scan(&op)
	if err != nil {
		fmt.Println("Invalid operation")
		return
	}

	switch op {
	case "+":
		fmt.Println(num1 + num2)
	case "-":
		fmt.Println(num1 - num2)
	case "*":
		fmt.Println(num1 * num2)
	case "/":
		if num2 == 0 {
			fmt.Println("Invalid operation")
			return
		}

		fmt.Println(num1 / num2)
	default:
		fmt.Println("Invalid operation")
		return
	}
}
