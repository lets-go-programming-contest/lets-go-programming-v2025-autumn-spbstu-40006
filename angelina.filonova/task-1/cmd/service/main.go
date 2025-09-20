package main

import (
	"errors"
	"fmt"
)

var (
	ErrorDivisionByZero = errors.New("Division by zero")
	ErrorInvalidOperation = errors.New("Invalid operation")
	ErrorInvalidFirstOperand = errors.New("Invalid first operand")
	ErrorInvalidSecondOperand = errors.New("Invalid second operand")
)

func add(a int, b int) int {
	return a + b
}

func sub(a int, b int) int {
	return a - b
}

func multiply(a int, b int) int {
	return a * b
}

func divide(a int, b int) (int, error) {
	if b == 0 {
		return 0, ErrorDivisionByZero
	}
	return a / b, nil
}

func main() {
	var num1, num2 int
	var operand string
	fmt.Scan(&num1, &num2, &operand)

	var result int
	var err error 

	switch operand {
	case "+":
		result = add(num1, num2)
	case "-":
		result = sub(num1, num2)
	case "*":
		result = multiply(num1, num2)
	case "/":
		result, err = divide(num1, num2) 
	default:
		err = ErrorInvalidOperation
	}

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(result)
}
