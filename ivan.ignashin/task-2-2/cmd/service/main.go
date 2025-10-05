package main

import (
	"container/heap"
	"fmt"
)

type MaxHeap []int

func (h *MaxHeap) Len() int           { return len(*h) }
func (h *MaxHeap) Less(i, j int) bool { return (*h)[i] > (*h)[j] }
func (h *MaxHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *MaxHeap) Push(x interface{}) {
	value, ok := x.(int)
	if !ok {
		return
	}

	*h = append(*h, value)
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

	_, err := fmt.Scan(&numberOfElements)
	if err != nil {
		return
	}

	ratings := make([]int, numberOfElements)
	for i := range numberOfElements {
		_, err = fmt.Scan(&ratings[i])
		if err != nil {
			return
		}
	}

	_, err = fmt.Scan(&kElement)
	if err != nil {
		return
	}

	ratingsHeap := &MaxHeap{}
	heap.Init(ratingsHeap)

	for _, rating := range ratings {
		heap.Push(ratingsHeap, rating)
	}

	for range kElement - 1 {
		heap.Pop(ratingsHeap)
	}

	result := heap.Pop(ratingsHeap)
	fmt.Println(result)
}
