package main

import "fmt"

func main() {
	const baseLo, baseHi = 15, 30

	var k int
	if n, err := fmt.Scan(&k); err != nil || n != 1 || k < 0 {
		return
	}

	lo, hi := baseLo, baseHi

	for i := 0; i < k; i++ {
		var op string
		var t int

		n, err := fmt.Scan(&op, &t)
		if err != nil || n != 2 {
			return
		}

		switch op {
		case ">=":
			if t > lo {
				lo = t
			}
		case "<=":
			if t < hi {
				hi = t
			}
		default:
			return
		}

		if lo <= hi {
			fmt.Println(lo)
		} else {
			fmt.Println(-1)
		}
	}
}
