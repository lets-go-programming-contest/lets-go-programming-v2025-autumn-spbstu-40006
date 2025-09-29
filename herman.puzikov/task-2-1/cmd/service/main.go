package main

import (
	"fmt"
	"os"
)

func main() {
	var n int
	if _, err := fmt.Scan(&n); err != nil {
		return
	}

	for range n {
		var k int
		if _, err := fmt.Scan(&k); err != nil {
			fmt.Fprintln(os.Stderr, "couldn't read the department count:", err)
		}

		lowerBound, higherBound := 15, 30

		for range k {
			var operator string
			var t int
			if _, err := fmt.Scan(&operator, &t); err != nil {
				fmt.Fprintln(os.Stderr, "couldn't read employee:", err)
				return
			}

			switch operator {
			case "<=":
				if t < higherBound {
					higherBound = t
				}
			case ">=":
				if t > lowerBound {
					lowerBound = t
				}
			default:
				fmt.Fprintln(os.Stderr, "unsupported operator:", operator)
				return
			}

			if lowerBound > higherBound {
				fmt.Println(-1)
			} else {
				fmt.Println(higherBound)
			}
		}
	}
}
