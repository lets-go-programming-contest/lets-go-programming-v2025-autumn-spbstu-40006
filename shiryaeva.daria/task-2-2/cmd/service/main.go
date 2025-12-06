package main

import (
	"container/heap"
	"fmt"
)

type IntHeap []int

func (h *IntHeap) Len() int { return len(*h) }

func (h *IntHeap) Less(i, j int) bool { return (*h)[i] > (*h)[j] }

func (h *IntHeap) Swap(i, j int) { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *IntHeap) Push(x interface{}) {
	if value, ok := x.(int); ok {
		*h = append(*h, value)
	}
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

func main() {
	var numberOfElements, kElement int

	_, err := fmt.Scan(&numberOfElements)
	if err != nil {
		return
	}

	ratings := make([]int, numberOfElements)
	for index := range numberOfElements {
		_, err = fmt.Scan(&ratings[index])
		if err != nil {
			return
		}
	}

	_, err = fmt.Scan(&kElement)
	if err != nil {
		return
	}

	ratingsHeap := &IntHeap{}
	heap.Init(ratingsHeap)

	for _, rating := range ratings {
		heap.Push(ratingsHeap, rating)
	}

	for range kElement - 1 {
		heap.Pop(ratingsHeap)
	}

	result, ok := heap.Pop(ratingsHeap).(int)
	if !ok {
		return
	}

	fmt.Println(result)
}
