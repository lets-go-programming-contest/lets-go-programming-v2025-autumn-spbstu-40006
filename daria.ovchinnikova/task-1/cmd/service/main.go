package main

import "fmt"

func main() {
	var firstOperand int
	var secondOperand int
	var operation string

	_, errFirstOp := fmt.Scan(&firstOperand)
	_, errSecondOp := fmt.Scan(&secondOperand)
	_, errOp := fmt.Scan(&operation)
	if errFirstOp != nil {
		fmt.Println("Invalid first operand")
		return
	}
	if errSecondOp != nil {
		fmt.Println("Invalid second operand")
		return
	}
	if errOp != nil {
		fmt.Println("Invalid operation")
		return
	}

	switch operation {
	case "+":
		fmt.Println(firstOperand + secondOperand)
	case "-":
		fmt.Println(firstOperand - secondOperand)
	case "*":
		fmt.Println(firstOperand * secondOperand)
	case "/":
		if secondOperand == 0 {
			fmt.Println("Division by zero")
			return
		}
		fmt.Println(firstOperand / secondOperand)
	default:
		fmt.Println("Invalid operation")
	}
}
