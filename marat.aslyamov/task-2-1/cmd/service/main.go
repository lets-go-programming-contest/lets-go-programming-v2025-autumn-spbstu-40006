package main

import "fmt"

func main() {
	var (
		N, K, minTemp, maxTemp, temp int
		operand                      string
	)
	_, err := fmt.Scan(&N)
	if err != nil {
		fmt.Println("Incorrect N argument")
		return
	}
	for range N {
		minTemp = 15
		maxTemp = 30
		_, err = fmt.Scan(&K)
		if err != nil {
			fmt.Println("Incorrect K argument")
			return
		}
		for range K {
			_, err = fmt.Scan(&operand, &temp)
			if err == nil && operand[0] == '>' && temp > minTemp {
				minTemp = temp
			} else if err == nil && operand[0] == '<' && temp < maxTemp {
				maxTemp = temp
			}

			if minTemp <= maxTemp {
				fmt.Println(minTemp)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
