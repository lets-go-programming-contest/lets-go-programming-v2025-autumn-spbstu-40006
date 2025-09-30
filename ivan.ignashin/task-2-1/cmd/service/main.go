package main

import (
	"fmt"
)

func main() {
	var (
		numberOfOtdels, numberOfWorkers int
		querryOperand                   string
		querryNumber                    int
	)

	fmt.Scan(&numberOfOtdels)

	for range numberOfOtdels {
		var (
			lefBorder   int = 15
			rightBorder int = 30
		)
		fmt.Scan(&numberOfWorkers)

		for range numberOfWorkers {
			_, err := fmt.Scan(&querryOperand, &querryNumber)
			if querryOperand[0] == '<' {
				if err == nil && querryNumber < rightBorder {
					rightBorder = querryNumber
				}
			} else if querryOperand[0] == '>' {
				if err == nil && querryNumber > lefBorder {
					lefBorder = querryNumber
				}
			}

			if rightBorder >= lefBorder {
				fmt.Println(lefBorder)
			} else {
				print(-1)
				break
			}
		}
	}
}
