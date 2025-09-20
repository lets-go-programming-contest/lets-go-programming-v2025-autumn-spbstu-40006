package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
)

var (
	ErrDivisionByZero   = errors.New("Division by zero")
	ErrInvalidOperation = errors.New("Invalid operation")
)

func isOperationSupported(operation string) bool {
	supported := map[string]bool{
		"+": true,
		"-": true,
		"*": true,
		"/": true,
	}
	_, exists := supported[operation]
	return exists
}

func calculate(a, b int, operation string) (int, error) {
	if !isOperationSupported(operation) {
		return 0, ErrInvalidOperation
	}

	switch operation {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, ErrDivisionByZero
		}
		return a / b, nil
	default:
		return 0, ErrInvalidOperation
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	firstOperand := scanner.Text()
	a, err := strconv.Atoi(firstOperand)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	scanner.Scan()
	secondOperand := scanner.Text()
	b, err := strconv.Atoi(secondOperand)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	scanner.Scan()
	operation := scanner.Text()
	if !isOperationSupported(operation) {
		fmt.Println("Invalid operation")
		return
	}

	result, err := calculate(a, b, operation)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(result)
}
