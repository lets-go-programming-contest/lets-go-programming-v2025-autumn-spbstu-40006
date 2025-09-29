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

	for dept := 0; dept < n; dept++ {
		var k int
		fmt.Scan(&k)

		lo, hi := 15, 30

		for i := 0; i < k; i++ {
			var op string
			var t int
			if _, err := fmt.Scan(&op, &t); err != nil {
				fmt.Fprintln(os.Stderr, "couldn't read employee", i, ":", err)
				return
			}

			switch op {
			case "<=":
				if t < hi {
					hi = t
				}
			case ">=":
				if t > lo {
					lo = t
				}
			default:
				fmt.Fprintln(os.Stderr, "unsupported operator:", op)
				return
			}

			if lo > hi {
				fmt.Println(-1)
			} else {
				fmt.Println(hi)
			}
		}
	}
}
