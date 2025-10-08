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
		log.Fatal("Invalid first operand")
	}

	_, err = fmt.Scan(&num2)
	if err != nil {
		log.Fatal("Invalid second operand")
	}

	_, err = fmt.Scan(&op)
	if err != nil {
		log.Fatal("Invalid operation")
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
			log.Fatal("Invalid operation")
		}

		fmt.Println(num1 / num2)
	}
}
