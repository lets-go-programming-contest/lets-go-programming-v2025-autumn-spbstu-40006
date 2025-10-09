package main

import "fmt"

func main() {
	var N int

	_, err := fmt.Scan(&N)
	if err != nil {
		fmt.Println("Invalid number of departments")
		return
	}

	for i := 0; i < N; i++ {
		minT := 15
		maxT := 30
		var K int

		_, err = fmt.Scan(&K)
		if err != nil {
			fmt.Println("Invalid number of workers")
			return
		}

		for j := 0; j < K; j++ {
			var op string
			var T int
			_, err = fmt.Scan(&op, &T)

			if err == nil && op == ">=" {
				if T > minT {
					minT = T
				}
			} else if err == nil && op == "<=" {
				if T < maxT {
					maxT = T
				}
			}

			if minT <= maxT {
				fmt.Println(minT)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
