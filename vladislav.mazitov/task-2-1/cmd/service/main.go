package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func istreamInt(reader *bufio.Reader) (int, error) {
	input, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	input = strings.TrimSpace(input)

	val, err := strconv.Atoi(input)
	if err != nil {
		return 0, err
	}

	return val, nil
}

func istreamString(reader *bufio.Reader) (string, error) {
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(input), nil
}

func checkTemperatures(check string, minValue *int, maxValue *int) int {
	switch check[:2] {
	case ">=":
		temp, err := strconv.Atoi(strings.TrimSpace(check[3:]))
		if err != nil {
			log.Fatal(err)
		}

		if temp > *minValue {
			*minValue = temp
		}
	case "<=":
		temp, err := strconv.Atoi(strings.TrimSpace(check[3:]))
		if err != nil {
			log.Fatal(err)
		}

		if temp < *maxValue {
			*maxValue = temp
		}
	}

	if *minValue <= *maxValue {

		return *minValue
	} else {

		return -1
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	num, err := istreamInt(reader)
	if err != nil {
		log.Fatal(err)
	}

	for range num {
		key, err := istreamInt(reader)
		if err != nil {
			log.Fatal(err)
		}

		var mnVal, mxVal = 15, 30

		for range key {
			str, err := istreamString(reader)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(checkTemperatures(str, &mnVal, &mxVal))
		}
	}
}
