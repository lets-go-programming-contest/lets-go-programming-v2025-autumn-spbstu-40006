package main

import (
	"container/heap"
	"fmt"
)

type MinHeap []int

func (h *MinHeap) Len() int           { return len(*h) }
func (h *MinHeap) Less(i, j int) bool { return (*h)[i] < (*h)[j] }
func (h *MinHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *MinHeap) Push(x interface{}) {
	value, ok := x.(int)
	if !ok {
		panic("MinHeap.Push: value is not int")
	}

	*h = append(*h, value)
}

func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

func main() {
	var dishCount, preferenceIndex int

	_, err := fmt.Scan(&dishCount)
	if err != nil {
		fmt.Println("Invalid input")

		return
	}

	scores := make([]int, dishCount)
	for i := range dishCount {
		_, err = fmt.Scan(&scores[i])
		if err != nil {
			fmt.Println("Invalid input")

			return
		}
	}

	_, err = fmt.Scan(&preferenceIndex)
	if err != nil {
		fmt.Println("Invalid input")

		return
	}

	heapContainer := &MinHeap{}
	heap.Init(heapContainer)

	for i := range preferenceIndex {
		heap.Push(heapContainer, scores[i])
	}

	for i := preferenceIndex; i < len(scores); i++ {
		if scores[i] > (*heapContainer)[0] {
			heap.Pop(heapContainer)
			heap.Push(heapContainer, scores[i])
		}
	}

	fmt.Println((*heapContainer)[0])
}
