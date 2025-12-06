package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	if !scanner.Scan() || scanner.Text() == "" {
		fmt.Println("Invalid first operand")
		return
	}
	x, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	if !scanner.Scan() || scanner.Text() == "" {
		fmt.Println("Invalid second operand")
		return
	}
	y, err := strconv.Atoi(scanner.Text())
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	if !scanner.Scan() || scanner.Text() == "" {
		fmt.Println("Invalid operation")
		return
	}
	op := scanner.Text()

	switch op {
	case "+":
		fmt.Println(x + y)
	case "-":
		fmt.Println(x - y)
	case "*":
		fmt.Println(x * y)
	case "/":
		if y == 0 {
			fmt.Println("Division by zero")
			return
		}
		fmt.Println(float64(x) / float64(y))
	default:
		fmt.Println("Invalid operation")
	}
}
