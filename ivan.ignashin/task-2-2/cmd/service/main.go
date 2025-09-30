package main

import (
	"container/heap"
	"fmt"
)

type MaxHeap []int

func (h MaxHeap) Len() int           { return len(h) }
func (h MaxHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h MaxHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MaxHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *MaxHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func main() {
	var numberOfElements, kElement int
	fmt.Scan(&numberOfElements)

	ratings := make([]int, numberOfElements)
	for i := 0; i < numberOfElements; i++ {
		fmt.Scan(&ratings[i])
	}

	fmt.Scan(&kElement)

	ratingsHeap := &MaxHeap{}
	heap.Init(ratingsHeap)

	for _, rating := range ratings {
		heap.Push(ratingsHeap, rating)
	}

	for i := 0; i < kElement-1; i++ {
		heap.Pop(ratingsHeap)
	}

	result := heap.Pop(ratingsHeap)
	fmt.Println(result)
}
