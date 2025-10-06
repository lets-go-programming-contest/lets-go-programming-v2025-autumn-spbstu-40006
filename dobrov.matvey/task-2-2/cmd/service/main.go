package main

import (
	"container/heap"
	"errors"
	"fmt"
)

var (
	ErrCountOutOfRange  = errors.New("count out of range")
	ErrRatingOutOfRange = errors.New("rating out of range")
	ErrOrderOutOfRange  = errors.New("order out of range")
	ErrUnexpectedType   = errors.New("unexpected type from heap.Pop")
)

type IntHeap []int

func (h *IntHeap) Len() int {
	return len(*h)
}

func (h *IntHeap) Less(i, j int) bool {
	return (*h)[i] > (*h)[j]
}

func (h *IntHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *IntHeap) Push(x any) {
	v, ok := x.(int)
	if !ok {
		return
	}

	*h = append(*h, v)
}

func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]

	return x
}

func readCount() (int, error) {
	var count int

	_, err := fmt.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("readCount scan: %w", err)
	}

	if count < 1 || count > 10000 {
		return 0, fmt.Errorf("readCount: %w", ErrCountOutOfRange)
	}

	return count, nil
}

func readRatings(count int) ([]int, error) {
	ratings := make([]int, count)

	for idx := range count {
		_, err := fmt.Scan(&ratings[idx])
		if err != nil {
			return nil, fmt.Errorf("readRatings scan #%d: %w", idx, err)
		}

		if ratings[idx] < -10000 || ratings[idx] > 10000 {
			return nil, fmt.Errorf("readRatings: %w", ErrRatingOutOfRange)
		}
	}

	return ratings, nil
}

func readOrder(count int) (int, error) {
	var order int

	_, err := fmt.Scan(&order)
	if err != nil {
		return 0, fmt.Errorf("readOrder scan: %w", err)
	}

	if order < 1 || order > count {
		return 0, fmt.Errorf("readOrder: %w", ErrOrderOutOfRange)
	}

	return order, nil
}

func buildMaxHeap(values []int) *IntHeap {
	maxHeap := &IntHeap{}
	heap.Init(maxHeap)

	for _, v := range values {
		heap.Push(maxHeap, v)
	}

	return maxHeap
}

func kthMax(h *IntHeap, order int) (int, error) {
	for range order - 1 {
		heap.Pop(h)
	}

	val, ok := heap.Pop(h).(int)
	if !ok {
		return 0, fmt.Errorf("kthMax: %w", ErrUnexpectedType)
	}

	return val, nil
}

func main() {
	count, err := readCount()
	if err != nil {
		return
	}

	ratings, err := readRatings(count)
	if err != nil {
		return
	}

	order, err := readOrder(count)
	if err != nil {
		return
	}

	maxHeap := buildMaxHeap(ratings)

	ans, err := kthMax(maxHeap, order)
	if err != nil {
		return
	}

	fmt.Println(ans)
}
