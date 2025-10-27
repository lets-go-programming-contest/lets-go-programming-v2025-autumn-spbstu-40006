package main

import (
	"container/heap"
	"fmt"
)

type IntMinHeap []int

func (h *IntMinHeap) Len() int { return len(*h) }

func (h *IntMinHeap) Less(i, j int) bool { return (*h)[i] < (*h)[j] }

func (h *IntMinHeap) Swap(i, j int) { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *IntMinHeap) Push(x any) {
	v, ok := x.(int)
	if !ok {
		return
	}

	*h = append(*h, v)
}

func (h *IntMinHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]

	*h = old[:n-1]

	return x
}

func main() {
	var dishCount, kPosition int
	if _, err := fmt.Scan(&dishCount); err != nil || dishCount < 1 || dishCount > 10000 {
		return
	}

	ratings := make([]int, dishCount)

	for i := range dishCount {
		if _, err := fmt.Scan(&ratings[i]); err != nil {
			return
		}
	}

	if _, err := fmt.Scan(&kPosition); err != nil || kPosition < 1 || kPosition > dishCount {
		return
	}

	minHeap := &IntMinHeap{}
	heap.Init(minHeap)

	for _, rating := range ratings {
		heap.Push(minHeap, rating)

		if minHeap.Len() > kPosition {
			heap.Pop(minHeap)
		}
	}

	fmt.Println((*minHeap)[0])
}
