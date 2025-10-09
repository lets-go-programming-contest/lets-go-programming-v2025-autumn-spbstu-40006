package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	var departments int
	if _, err := fmt.Fscan(reader, &departments); err != nil {
		fmt.Fprintln(os.Stderr, "failed to read number of departments:", err)

		return
	}

	for deptIndex := range make([]struct{}, departments) {
		_ = deptIndex // silence unused var if not needed

		var employees int
		if _, err := fmt.Fscan(reader, &employees); err != nil {
			fmt.Fprintln(os.Stderr, "failed to read number of employees:", err)

			return
		}

		low := 15
		high := 30

		for empIndex := 0; empIndex < employees; empIndex++ {
			line, err := reader.ReadString('\n')
			if err != nil {
				fmt.Fprintln(os.Stderr, "failed to read employee preference:", err)

				return
			}

			line = strings.TrimSpace(line)
			if line == "" {
				empIndex--

				continue
			}

			var value int
			if strings.HasPrefix(line, ">=") {
				if n, err := fmt.Sscanf(line, ">= %d", &value); err == nil && n == 1 && value > low {
					low = value
				}
			} else if strings.HasPrefix(line, "<=") {
				if n, err := fmt.Sscanf(line, "<= %d", &value); err == nil && n == 1 && value < high {
					high = value
				}
			}
		}

		if low <= high {
			fmt.Println(low)
		} else {
			fmt.Println(-1)
		}
	}
}