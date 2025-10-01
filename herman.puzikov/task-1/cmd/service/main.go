package main

import (
	"errors"
	"fmt"
	"os"
)

var (
	ErrorDivisionByZero   = errors.New("Division by zero")
	ErrorInvalidOperation = errors.New("Invalid operation")
)

func calculate(operator string, a, b int) (float64, error) {
	switch operator {
	case "+":
		return float64(a + b), nil
	case "-":
		return float64(a - b), nil
	case "*":
		return float64(a * b), nil
	case "/":
		{
			if b == 0 {
				return 0, ErrorDivisionByZero
			}
			return float64(a) / float64(b), nil
		}
	default:
		return 0, ErrorInvalidOperation
	}
}

func main() {
	var (
		fOp, sOp int
		operator string
	)

	if _, err := fmt.Scan(&fOp); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid first operand\n")
		return
	}

	if _, err := fmt.Scan(&sOp); err != nil {
		fmt.Fprintf(os.Stderr, "Invalid second operand\n")
		return
	}

	if _, err := fmt.Scan(&operator); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", ErrorInvalidOperation)
		return
	}

	if res, err := calculate(operator, fOp, sOp); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	} else {
		fmt.Println(res)
	}
}
