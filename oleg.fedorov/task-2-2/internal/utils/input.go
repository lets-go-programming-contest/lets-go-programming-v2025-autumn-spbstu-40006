package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func ReadRatingsAndK() ([]int, int, error) {
	scanner := bufio.NewScanner(os.Stdin)

	numOfDishes, err := readNumberOfDishes(scanner)
	if err != nil {
		return nil, 0, err
	}

	if numOfDishes == 0 {
		return []int{}, 0, nil
	}

	ratings, err := readRatings(scanner, numOfDishes)
	if err != nil {
		return nil, 0, err
	}

	kthNumber, err := readKthNumber(scanner, numOfDishes)
	if err != nil {
		return nil, 0, err
	}

	return ratings, kthNumber, nil
}

func readNumberOfDishes(scanner *bufio.Scanner) (int, error) {
	scanner.Scan()

	numOfDishes, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return 0, fmt.Errorf("invalid number of dishes: %v", err)
	}

	if numOfDishes < 0 || numOfDishes > 10000 {
		return 0, fmt.Errorf("number of dishes out of range: %d", numOfDishes)
	}

	return numOfDishes, nil
}

func readRatings(scanner *bufio.Scanner, numOfDishes int) ([]int, error) {
	scanner.Scan()

	ratingText := scanner.Text()
	ratingStrs := strings.Fields(ratingText)

	if len(ratingStrs) != numOfDishes {
		return nil, fmt.Errorf("number of ratings doesn't match number of dishes")
	}

	ratings := make([]int, numOfDishes)
	for i, str := range ratingStrs {
		rating, err := strconv.Atoi(str)
		if err != nil {
			return nil, fmt.Errorf("invalid rating: %s", str)
		}

		if rating < -10000 || rating > 10000 {
			return nil, fmt.Errorf("rating out of range: %d", rating)
		}

		ratings[i] = rating
	}

	return ratings, nil
}

func readKthNumber(scanner *bufio.Scanner, numOfDishes int) (int, error) {
	scanner.Scan()

	kthNumber, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return 0, fmt.Errorf("invalid kth number: %v", err)
	}

	if kthNumber < 1 || kthNumber > numOfDishes {
		return 0, fmt.Errorf("kth number out of range: %d", kthNumber)
	}

	return kthNumber, nil
}
