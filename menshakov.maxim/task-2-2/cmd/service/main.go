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

func (h MaxHeap) Len() int {
	return len(h)
}

func (h MaxHeap) Less(i, j int) bool {
	return h[i] > h[j]
}

func (h MaxHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MaxHeap) Push(x any) {
	num, ok := x.(int)
	if !ok {
		return
	}

	*h = append(*h, num)
}

func (h *MaxHeap) Pop() any {
	old := *h
	n := len(old)
	elem := old[n-1]
	*h = old[:n-1]
	return elem
}

func getCount() (int, error) {
	var n int

	if _, err := fmt.Scan(&n); err != nil {
		return 0, fmt.Errorf("failed to read count: %w", err)
	}

	if n < 1 || n > 10000 {
		return 0, fmt.Errorf("%w: %d", ErrInvalidCount, n)
	}

	return n, nil
}

func getRatings(n int) ([]int, error) {
	values := make([]int, n)

	for i := 0; i < n; i++ {
		if _, err := fmt.Scan(&values[i]); err != nil {
			return nil, fmt.Errorf("read rating #%d: %w", i+1, err)
		}

		if values[i] < -10000 || values[i] > 10000 {
			return nil, fmt.Errorf("%w: %d", ErrInvalidRating, values[i])
		}
	}

	return values, nil
}

func getOrder(n int) (int, error) {
	var k int

	if _, err := fmt.Scan(&k); err != nil {
		return 0, fmt.Errorf("failed to read order: %w", err)
	}

	if k < 1 || k > n {
		return 0, fmt.Errorf("%w: %d", ErrInvalidOrder, k)
	}

	return k, nil
}

func createHeap(nums []int) *MaxHeap {
	h := &MaxHeap{}
	heap.Init(h)

	for _, v := range nums {
		heap.Push(h, v)
	}

	return h
}

func extractKth(h *MaxHeap, k int) (int, error) {
	for i := 1; i < k; i++ {
		heap.Pop(h)
	}

	res, ok := heap.Pop(h).(int)
	if !ok {
		return 0, fmt.Errorf("%w in extractKth", ErrHeapType)
	}

	return res, nil
}

func main() {
	count, err := getCount()
	if err != nil {
		return
	}

	ratings, err := getRatings(count)
	if err != nil {
		return
	}

	order, err := getOrder(count)
	if err != nil {
		return
	}

	h := createHeap(ratings)

	result, err := extractKth(h, order)
	if err != nil {
		return
	}

	fmt.Println(result)
}
