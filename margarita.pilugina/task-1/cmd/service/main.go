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

func sum(a, b int) (int, error)  { return a + b, nil }
func sub(a, b int) (int, error)  { return a - b, nil }
func mult(a, b int) (int, error) { return a * b, nil }
func div(a, b int) (int, error) {
	if b == 0 {
		return 0, ErrDivByZero
	}
	return a / b, nil
}

func calculator(a, b int, operation string) (int, error) {
	switch operation {
	case "+":
		return sum(a, b)
	case "-":
		return sub(a, b)
	case "*":
		return mult(a, b)
	case "/":
		return div(a, b)
	default:
		return 0, ErrInvalidOperation
	}
}

func main() {
	var (
		a, b int
		op   string
	)

	if _, err := fmt.Scan(&a); err != nil {
		fmt.Println(ErrInvalidFirstOperand)
		return
	}

	if _, err := fmt.Scan(&b); err != nil {
		fmt.Println(ErrInvalidSecondOperand)
		return
	}

	if _, err := fmt.Scan(&op); err != nil {
		fmt.Println(ErrInvalidOperation)
		return
	}

	result, err := calculator(a, b, op)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(result)
}
