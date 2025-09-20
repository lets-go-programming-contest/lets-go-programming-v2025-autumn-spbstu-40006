package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()

	depCount, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return
	}

	for index := 0; index < depCount; index++ {
		scanner.Scan()

		workersNumber, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return
		}

		minTemp := 15
		maxTemp := 30

		for jndex := 0; jndex < workersNumber; jndex++ {
			scanner.Scan()
			line := scanner.Text()
			parts := strings.Fields(line)
			operation := parts[0]

			temp, err := strconv.Atoi(parts[1])
			if err != nil {
				return
			}

			switch operation {
			case ">=":
				if temp > minTemp {
					minTemp = temp
				}
			case "<=":
				if temp < maxTemp {
					maxTemp = temp
				}
			}

			if minTemp <= maxTemp {
				fmt.Println(minTemp)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
