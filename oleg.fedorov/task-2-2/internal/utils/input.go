package utils

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func ReadRatingsAndK() ([]int, int) {
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()

	numOfDishes, err := strconv.Atoi(scanner.Text())
	if err != nil || numOfDishes < 0 || numOfDishes > 10000 {
		return nil, 0
	}

	if numOfDishes == 0 {
		return []int{}, 0
	}

	scanner.Scan()

	ratingText := scanner.Text()
	ratingStrs := strings.Fields(ratingText)

	if len(ratingStrs) != numOfDishes {
		return nil, 0
	}

	ratings := make([]int, numOfDishes)

	for i, str := range ratingStrs {
		rating, err := strconv.Atoi(str)
		if err != nil || rating < -10000 || rating > 10000 {
			return nil, 0
		}

		ratings[i] = rating
	}

	scanner.Scan()

	kthNumber, err := strconv.Atoi(scanner.Text())
	if err != nil || kthNumber < 1 || kthNumber > numOfDishes {
		return nil, 0
	}

	return ratings, kthNumber
}
