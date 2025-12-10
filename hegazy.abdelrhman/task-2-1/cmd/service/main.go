package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	var departments, workers int

	if _, err := fmt.Fscan(reader, &departments); err != nil {
		return
	}

	for range departments {
		if _, err := fmt.Fscan(reader, &workers); err != nil {
			return
		}

		minTemp := 15
		maxTemp := 30

		for range workers {
			var (
				operator    string
				temperature int
			)

			if _, err := fmt.Fscan(reader, &operator, &temperature); err != nil {
				return
			}

			switch operator {
			case ">=":
				if temperature > minTemp {
					minTemp = temperature
				}
			case "<=":
				if temperature < maxTemp {
					maxTemp = temperature
				}
			}

			if minTemp > maxTemp {
				fmt.Println(-1)
			} else {
				fmt.Println(minTemp)
			}
		}
	}
}
