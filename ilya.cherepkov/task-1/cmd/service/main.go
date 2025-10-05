package main

import (
	"errors"
	"fmt"
)

var (
	ErrorInvalidOperation     = errors.New("Invalid operation")
	ErrorInvalidFirstOperand  = errors.New("Invalid first operand")
	ErrorInvalidSecondOperand = errors.New("Invalid second operand")
	ErrorDivisionByZero       = errors.New("Division by zero")
)

func calculate(a, b int, op string) (int, error) {
	switch op {
	case "+":
		return a + b, nil
	case "-":
		return a - b, nil
	case "*":
		return a * b, nil
	case "/":
		if b == 0 {
			return 0, ErrorDivisionByZero
		}
		return a / b, nil
	default:
		return 0, ErrorInvalidOperation
	}
}

func main() {
	var (
		a, b   int
		op     string
		err    error
		result int
	)

	if _, err = fmt.Scan(&a); err != nil {
		fmt.Println(ErrorInvalidFirstOperand)
		return
	}

	if _, err = fmt.Scan(&b); err != nil {
		fmt.Println(ErrorInvalidSecondOperand)
		return
	}

	if _, err = fmt.Scan(&op); err != nil {
		fmt.Println(ErrorInvalidOperation)
		return
	}

	result, err = calculate(a, b, op)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(result)
}
