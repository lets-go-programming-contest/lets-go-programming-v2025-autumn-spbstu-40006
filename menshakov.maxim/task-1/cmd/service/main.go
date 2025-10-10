package main

import "fmt"

func main() {
	var (
		num_1, num_2 int
		op_1         string
	)

	_, err := fmt.Scan(&num_1)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}

	_, err = fmt.Scan(&num_2)
	if err != nil {
		fmt.Println("Invalid second operand")
		return
	}

	_, err = fmt.Scan(&op_1)
	if err != nil || (op_1 != "+" && op_1 != "-" && op_1 != "*" && op_1 != "/") {
		fmt.Println("Invalid operation")
		return
	}

	switch op_1 {
	case "+":
		fmt.Println(num_1 + num_2)
	case "-":
		fmt.Println(num_1 - num_2)
	case "*":
		fmt.Println(num_1 * num_2)
	case "/":
		if num_2 != 0 {
			fmt.Println(num_1 / num_2)
		} else {
			fmt.Println("Division by zero")
		}
	}
}