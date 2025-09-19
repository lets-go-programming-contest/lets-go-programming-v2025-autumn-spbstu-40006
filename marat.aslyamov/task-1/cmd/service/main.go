package main

import "fmt"

func main() {
	var i, j, result int
	var sign string
	_, err := fmt.Scan(&i)
	if err != nil {
		fmt.Println("Invalid first operand")
		return
	}
	_, err = fmt.Scan(&j)
	if err != nil {
		fmt.Println("Invalind second operand")
		return
	}
	_, err = fmt.Scan(&sign)
	if err != nil {
		fmt.Println("Invalid operation")
	}
	switch sign {
	case "+":
		result = i + j
	case "-":
		result = i - j
	case "*":
		result = i * j
	case "/":
		if j == 0 {
			fmt.Println("Division by zero")
			return
		}
		result = i / j
	default:
		fmt.Println("Invalid operation")
		return
	}
	fmt.Println(result)
}
