package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	var N int
	fmt.Fscan(reader, &N) // number of departments

	for i := 0; i < N; i++ {
		var K int
		fmt.Fscan(reader, &K) // number of employees in this department

		low := 15
		high := 30

		for j := 0; j < K; j++ {
			line, _ := reader.ReadString('\n')
			line = strings.TrimSpace(line)

			// Sometimes after fmt.Fscan, the first line can be empty due to newline
			if line == "" {
				j--
				continue
			}

			if strings.HasPrefix(line, ">=") {
				var value int
				fmt.Sscanf(line, ">= %d", &value)
				if value > low {
					low = value
				}
			} else if strings.HasPrefix(line, "<=") {
				var value int
				fmt.Sscanf(line, "<= %d", &value)
				if value < high {
					high = value
				}
			}

			if low <= high {
				fmt.Println(low)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
