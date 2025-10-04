package main

import (
	"fmt"
)

func main() {
	var (
		numberOfOtdels, numberOfWorkers int
		querryOperand                   string
		querryNumber                    int
		leftBorder, rightBorder         int
	)

	_, err := fmt.Scan(&numberOfOtdels)
	if err != nil {
		fmt.Println("Error! Imvalod number of departments.")
	}

	for range numberOfOtdels {
		leftBorder = 15
		rightBorder = 30

		_, err := fmt.Scan(&numberOfWorkers)
		if err != nil {
			fmt.Println("Error! Imvalod number of workers.")
		}

		for range numberOfWorkers {
			_, err := fmt.Scan(&querryOperand, &querryNumber)
			if querryOperand[0] == '<' {
				if err == nil && querryNumber < rightBorder {
					rightBorder = querryNumber
				}
			} else if querryOperand[0] == '>' {
				if err == nil && querryNumber > leftBorder {
					leftBorder = querryNumber
				}
			}

			if rightBorder >= leftBorder {
				fmt.Println(leftBorder)
			} else {
				println(-1)
			}
		}
	}
}
