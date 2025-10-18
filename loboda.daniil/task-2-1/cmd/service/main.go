package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	baseLower = 15
	baseUpper = 30
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var departments, employees int

	if _, err := fmt.Fscan(in, &departments, &employees); err != nil {

		return
	}

	for i := 0; i < departments; i++ {
		lower, upper := baseLower, baseUpper
		for j := 0; j < employees; j++ {
			var op string
			var t int

			if _, err := fmt.Fscan(in, &op, &t); err != nil {

				return
			}
			switch op {
			case ">=":
				if t > lower {
					lower = t
				}
			case "<=":
				if t < upper {
					upper = t
				}
			default:

				return
			}

			if lower <= upper {
				fmt.Fprintln(out, lower)
			} else {
				fmt.Fprintln(out, -1)
			}
		}
	}
}
