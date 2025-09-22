package main

import (
	"fmt"
)

func plus(a, b int) int {
	return a + b
}

func minus(a, b int) int {
	return a - b
}

func multiply(a, b int) int {
	return a * b
}

func devide(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("devide by zero")
	}
	return a / b, nil
}

func work() {
	var (
		first, second int
		operation     string
	)

	if _, err := fmt.Scan(&first, &second, &operation); err != nil {
		fmt.Printf("Invalid values\n")
		return
	}

	switch operation {
	case "+":
		fmt.Printf("Result is: %d", plus(first, second))

	case "-":
		fmt.Printf("Result is: %d", minus(first, second))
	case "*":
		fmt.Printf("Result is: %d", multiply(first, second))
	case "/":
		temp, err := devide(first, second)
		if err != nil {
			fmt.Printf("Error: %v", err)
			return
		}
		fmt.Printf("Result is: %d", temp)
	default:
		fmt.Println("Invalid operation")
	}
}

func main() {
	work()
}
