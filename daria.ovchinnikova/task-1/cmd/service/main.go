package main

import "fmt"

func main() {
	var first_operand int
	var second_operand int
	var operation string

	_, err_first_op := fmt.Scan(&first_operand)
	_, err_second_op := fmt.Scan(&second_operand)
	fmt.Scan(&operation)
	if err_first_op != nil {
		fmt.Println("Invalid first operand")
		return
	}
	if err_second_op != nil {
		fmt.Println("Invalid second operand")
		return
	}

	switch operation {
	case "+":
		fmt.Println(first_operand + second_operand)
	case "-":
		fmt.Println(first_operand - second_operand)
	case "*":
		fmt.Println(first_operand * second_operand)
	case "/":
		if second_operand == 0 {
			fmt.Println("Division by zero")
			return
		}
		fmt.Println(first_operand / second_operand)
	default:
		fmt.Println("Invalid operation")
	}
}
