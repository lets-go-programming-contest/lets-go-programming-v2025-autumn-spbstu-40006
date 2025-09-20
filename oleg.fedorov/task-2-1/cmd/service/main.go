package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type IntHeap []int

func (h *IntHeap) Len() int           { return len(*h) }
func (h *IntHeap) Less(i, j int) bool { return (*h)[i] < (*h)[j] }
func (h *IntHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *IntHeap) Push(x interface{}) {
	if num, ok := x.(int); ok {
		*h = append(*h, num)
	}
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

func readInput() ([]int, int, error) {
	scanner := bufio.NewScanner(os.Stdin)

	if !scanner.Scan() {
		return nil, 0, fmt.Errorf("err: %w", scanner.Err())
	}

	numOfDishes, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return nil, 0, fmt.Errorf("err: %w", err)
	}

	if numOfDishes < 1 || numOfDishes > 10000 {
		return nil, 0, fmt.Errorf("err: from 1 to 10000")
	}

	if !scanner.Scan() {
		return nil, 0, fmt.Errorf("err: %w", scanner.Err())
	}

	ratingText := scanner.Text()
	ratingStrs := strings.Fields(ratingText)
	if len(ratingStrs) != numOfDishes {
		return nil, 0, fmt.Errorf("incorrect rating count: %d, %d", numOfDishes, len(ratingStrs))
	}

	ratings := make([]int, numOfDishes)
	for index, str := range ratingStrs {
		rating, err := strconv.Atoi(str)
		if err != nil {
			return nil, 0, fmt.Errorf("err: %w", err)
		}

		if rating < -10000 || rating > 10000 {
			return nil, 0, fmt.Errorf("err (from -10000 to 10000): %d", rating)
		}

		ratings[index] = rating
	}

	if !scanner.Scan() {
		return nil, 0, fmt.Errorf("err: %w", scanner.Err())
	}

	kthNumber, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return nil, 0, fmt.Errorf("err: %w", err)
	}

	if kthNumber < 1 || kthNumber > numOfDishes {
		return nil, 0, fmt.Errorf("err (from 1 to %d)", numOfDishes)
	}

	return ratings, kthNumber, nil
}

func main() {
	ratings, kthNumber, err := readInput()
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	intHeap := &IntHeap{}
	heap.Init(intHeap)

	for _, rating := range ratings {
		if intHeap.Len() < kthNumber {
			heap.Push(intHeap, rating)
			continue
		}

		if rating > (*intHeap)[0] {
			heap.Pop(intHeap)
			heap.Push(intHeap, rating)
		}
	}

	result := (*intHeap)[0]
	fmt.Println(result)
}
