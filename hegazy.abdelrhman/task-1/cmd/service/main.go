package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Read first operand
	firstInput, _ := reader.ReadString('\n')
	firstInput = strings.TrimSpace(firstInput)
	firstOperand, err1 := strconv.Atoi(firstInput)
	if err1 != nil {
		fmt.Println("Invalid first operand")
		return
	}

	// Read second operand
	secondInput, _ := reader.ReadString('\n')
	secondInput = strings.TrimSpace(secondInput)
	secondOperand, err2 := strconv.Atoi(secondInput)
	if err2 != nil {
		fmt.Println("Invalid second operand")
		return
	}

	// Read operation
	op, _ := reader.ReadString('\n')
	op = strings.TrimSpace(op)

	// Perform calculation
	switch op {
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
