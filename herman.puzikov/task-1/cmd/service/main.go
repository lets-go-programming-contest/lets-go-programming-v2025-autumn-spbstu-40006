package main

import (
	"errors"
	"fmt"
)

var (
	ErrorDivisionByZero       = errors.New("Division by zero")
	ErrorInvalidOperation     = errors.New("Invalid operation")
	ErrorInvalidFirstOperand  = errors.New("Invalid first operand")
	ErrorInvalidSecondOperand = errors.New("Invalid second operand")
)

func main() {
	var fOp, sOp int
	var operand string
	var functionsMap = map[string]func(int, int) (float64, error){
		"+": func(a, b int) (float64, error) { return float64(a + b), nil },
		"-": func(a, b int) (float64, error) { return float64(a - b), nil },
		"*": func(a, b int) (float64, error) { return float64(a * b), nil },
		"/": func(a, b int) (float64, error) {
			if b == 0 {
				return 0, ErrorDivisionByZero
			}
			return float64(a) / float64(b), nil
		},
	}

	if _, err := fmt.Scan(&fOp); err != nil {
		fmt.Println(ErrorInvalidFirstOperand)
		return
	}
	if _, err := fmt.Scan(&sOp); err != nil {
		fmt.Println(ErrorInvalidSecondOperand)
		return
	}
	if _, err := fmt.Scan(&operand); err != nil {
		fmt.Println(ErrorInvalidOperation)
		return
	}

	var operation, ok = functionsMap[operand]
	if !ok {
		fmt.Println(ErrorInvalidOperation)
		return
	}

	if res, err := operation(fOp, sOp); err != nil {
		fmt.Println(ErrorDivisionByZero)
	} else {
		fmt.Println(res)
	}
}
