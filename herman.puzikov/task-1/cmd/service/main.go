package main

import (
	"errors"
	"fmt"
	"os"
)

var (
	ErrorDivisionByZero   = errors.New("division by zero")
	ErrorInvalidOperation = errors.New("invalid operation")
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

	if _, err := fmt.Scan(&fOp, &sOp, &operator); err != nil {
		fmt.Fprintf(os.Stderr, "There has been an error while reading input: %v\n", err)
		return
	}

	if res, err := calculate(operator, fOp, sOp); err != nil {
		fmt.Fprintf(os.Stderr, "There has been an error while calculating: %v\n", err)
	} else {
		fmt.Println(res)
	}
}
