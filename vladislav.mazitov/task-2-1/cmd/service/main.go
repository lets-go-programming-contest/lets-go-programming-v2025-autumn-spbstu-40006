package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	input = strings.TrimSpace(input)
	n, err := strconv.Atoi(input)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < n; i++ {
		secondInput, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		secondInput = strings.TrimSpace(secondInput)
		k, err := strconv.Atoi(secondInput)
		if err != nil {
			log.Fatal(err)
		}
		var min, max int = 15, 30
		for j := 0; j < k; j++ {
			thirdInput, err := reader.ReadString('\n')
			if err != nil {
				log.Fatal(err)
			}
			check := (strings.TrimSpace(thirdInput))[:2]
			value, err := strconv.Atoi(strings.TrimSpace(thirdInput)[3:])
			if err != nil {
				log.Fatal(err)
			}
			switch check {
			case ">=":
				if value > min {
					min = value
				}
			case "<=":
				if value < max {
					max = value
				}
			}
			if min <= max {
				fmt.Println(min)
			} else {
				fmt.Println(-1)
			}
		}
	}
}
