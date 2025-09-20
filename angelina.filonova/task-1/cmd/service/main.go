package main

import "fmt"

func add(a int, b int) int {
	return a + b;
}

func sub(a int, b int) int {
	return a - b;
}

func multiply(a int, b int) int {
	return a * b;
}

func divide(a int, b int) int {
	return a / b;
}

func main() {

	var num1, num2 int
	var operand string
	fmt.Scan(&num1, &num2, &operand)

	switch operand {
	case "+":
		add(num1, num2)
	 case "-":
        sub(num1, num2)
    case "*":
        multiply(num1, num2)
    case "/":
        divide(num1, num2)

	}
}