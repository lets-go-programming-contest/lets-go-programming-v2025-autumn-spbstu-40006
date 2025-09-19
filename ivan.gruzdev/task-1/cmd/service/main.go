package main

import "fmt"

func main() {
	var value1, value2 int
	var operation string

	_, err1 := fmt.Scan(&value1)
	if err1 != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err2 := fmt.Scan(&value2)
	if err2 != nil {
		fmt.Println("Invalid second operand")
		return
	}

	_, err3 := fmt.Scan(&operation)
	if err3 != nil {
		fmt.Println("Invalid operation")
		return
	}

	switch operation {
	case "+":
		fmt.Println(value1 + value2)
	case "-":
		fmt.Println(value1 - value2)
	case "*":
		fmt.Println(value1 * value2)
	case "/":
		if value2 == 0 {
			fmt.Println("Division by zero")
			return
		}
		fmt.Println(value1 / value2)
	default:
		fmt.Println("Invalid operation")
	}
}
