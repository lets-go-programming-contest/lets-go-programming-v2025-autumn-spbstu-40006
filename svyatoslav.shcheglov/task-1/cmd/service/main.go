package main

import "fmt"

func main() {
	var (
		value1, value2 int
		sign           string
	)
	_, err := fmt.Scan(&value1)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err = fmt.Scan(&value2)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	_, err = fmt.Scan(&sign)
	if err != nil {
		fmt.Println("Invalid operation")
		return
	}

	switch sign {
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
		return
	}
}
