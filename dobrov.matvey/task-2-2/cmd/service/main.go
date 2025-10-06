package main

import (
	"container/heap"
	"fmt"
)

type IntHeap []int

func (h *IntHeap) Len() int           { return len(*h) }
func (h *IntHeap) Less(i, j int) bool { return (*h)[i] > (*h)[j] }
func (h *IntHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }
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

func main() {
	var countDish, desiredDish int
	_, err := fmt.Scan(&countDish)

	if err != nil {
		return
	}

	if countDish < 1 || countDish > 10000 {
		return
	}

	dishRatings := make([]int, countDish)
	for idx := range countDish {
		_, err = fmt.Scan(&dishRatings[idx])
		if err != nil {
			return
		}

		if dishRatings[idx] < -10000 || dishRatings[idx] > 10000 {
			return
		}
	}

	dishRatingsHeap := &IntHeap{}
	heap.Init(dishRatingsHeap)

	for _, dishRating := range dishRatings {
		heap.Push(dishRatingsHeap, dishRating)
	}

	_, err = fmt.Scan(&desiredDish)
	if err != nil {
		return
	}

	if desiredDish < 1 || desiredDish > countDish {
		return
	}

	for range desiredDish - 1 {
		heap.Pop(dishRatingsHeap)
	}

	fmt.Println(heap.Pop(dishRatingsHeap))
}
