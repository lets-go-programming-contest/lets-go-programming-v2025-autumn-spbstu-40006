package utils

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

var (
	ErrInvalidNumberOfDishes = errors.New("invalid number of dishes")
	ErrNumberOfDishesRange   = errors.New("number of dishes out of range")
	ErrRatingsCountMismatch  = errors.New("number of ratings doesn't match number of dishes")
	ErrInvalidRating         = errors.New("invalid rating")
	ErrRatingOutOfRange      = errors.New("rating out of range")
	ErrInvalidKthNumber      = errors.New("invalid kth number")
	ErrKthNumberOutOfRange   = errors.New("kth number out of range")
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
		return 0, ErrInvalidNumberOfDishes
	}

	if numOfDishes < 0 || numOfDishes > 10000 {
		return 0, ErrNumberOfDishesRange
	}

	return numOfDishes, nil
}

func readRatings(scanner *bufio.Scanner, numOfDishes int) ([]int, error) {
	scanner.Scan()

	ratingText := scanner.Text()
	ratingStrs := strings.Fields(ratingText)

	if len(ratingStrs) != numOfDishes {
		return nil, ErrRatingsCountMismatch
	}

	ratings := make([]int, numOfDishes)

	for index, str := range ratingStrs {
		rating, err := strconv.Atoi(str)
		if err != nil {
			return nil, ErrInvalidRating
		}

		if rating < -10000 || rating > 10000 {
			return nil, ErrRatingOutOfRange
		}

		ratings[index] = rating
	}

	return ratings, nil
}

func readKthNumber(scanner *bufio.Scanner, numOfDishes int) (int, error) {
	scanner.Scan()

	kthNumber, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return 0, ErrInvalidKthNumber
	}

	if kthNumber < 1 || kthNumber > numOfDishes {
		return 0, ErrKthNumberOutOfRange
	}

	return kthNumber, nil
}
