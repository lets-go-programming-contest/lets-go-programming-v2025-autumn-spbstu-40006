package main

import (
	"container/heap"
	"errors"
	"fmt"
)

var (
	ErrInvalidCount  = errors.New("invalid count range")
	ErrInvalidRating = errors.New("invalid rating value")
	ErrInvalidOrder  = errors.New("invalid order position")
	ErrHeapType      = errors.New("unexpected type from heap")
)

type MaxHeap []int

func (h *MaxHeap) Len() int {
	return len(*h)
}

func (h *MaxHeap) Less(i, j int) bool {
	return (*h)[i] > (*h)[j]
}

func (h *MaxHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *MaxHeap) Push(x any) {
	value, ok := x.(int)
	if !ok {
		return
	}

	*h = append(*h, value)
}

func (h *MaxHeap) Pop() any {
	oldHeap := *h
	n := len(oldHeap)
	elem := oldHeap[n-1]
	*h = oldHeap[:n-1]

	return elem
}

func readCount() (int, error) {
	var count int

	if _, err := fmt.Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to read count: %w", err)
	}

	if count < 1 || count > 10000 {
		return 0, fmt.Errorf("%w: %d", ErrInvalidCount, count)
	}

	return count, nil
}

func readRatings(count int) ([]int, error) {
	ratings := make([]int, count)

	for idx := range ratings {
		if _, err := fmt.Scan(&ratings[idx]); err != nil {
			return nil, fmt.Errorf("failed to read rating #%d: %w", idx+1, err)
		}

		if ratings[idx] < -10000 || ratings[idx] > 10000 {
			return nil, fmt.Errorf("%w: %d", ErrInvalidRating, ratings[idx])
		}
	}

	return ratings, nil
}

func readOrder(count int) (int, error) {
	var order int

	if _, err := fmt.Scan(&order); err != nil {
		return 0, fmt.Errorf("failed to read order: %w", err)
	}

	if order < 1 || order > count {
		return 0, fmt.Errorf("%w: %d", ErrInvalidOrder, order)
	}

	return order, nil
}

func buildMaxHeap(values []int) *MaxHeap {
	maxHeap := &MaxHeap{}
	heap.Init(maxHeap)

	for _, val := range values {
		heap.Push(maxHeap, val)
	}

	return maxHeap
}

func extractKthMax(maxHeap *MaxHeap, k int) (int, error) {
	for idx := 1; idx < k; idx++ {
		heap.Pop(maxHeap)
	}

	elem, ok := heap.Pop(maxHeap).(int)
	if !ok {
		return 0, fmt.Errorf("%w in extractKthMax", ErrHeapType)
	}

	return elem, nil
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

	result, err := extractKthMax(maxHeap, order)
	if err != nil {
		return
	}

	fmt.Println(result)
}
