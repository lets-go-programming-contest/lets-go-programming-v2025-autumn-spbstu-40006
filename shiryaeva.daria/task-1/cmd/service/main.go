package main

import (
	"errors"
	"fmt"
)

var (
	ErrDivByZero            = errors.New("Division by zero")
	ErrInvalidOperation     = errors.New("Invalid operation")
	ErrInvalidFirstOperand  = errors.New("Invalid first operand")
	ErrInvalidSecondOperand = errors.New("Invalid second operand")
)

func calculate(value1, value2 int, operation string) (int, error) {
	switch operation {
	case "+":
		return value1 + value2, nil
	case "-":
		return value1 - value2, nil
	case "*":
		return value1 * value2, nil
	case "/":
		if value2 == 0 {
			return 0, ErrDivByZero
		}
		return value1 / value2, nil
	default:
		return 0, ErrInvalidOperation
	}
}

func main() {
	var (
		value1, value2 int
		operation      string
	)

	_, err := fmt.Scan(&value1)
	if err != nil {
		fmt.Println(ErrInvalidFirstOperand)
		return
	}

	_, err = fmt.Scan(&value2)
	if err != nil {
		fmt.Println(ErrInvalidSecondOperand)
		return
	}

	_, err = fmt.Scan(&operation)
	if err != nil {
		fmt.Println(ErrInvalidOperation)
		return
	}

	result, err := calculate(value1, value2, operation)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(result)
}
