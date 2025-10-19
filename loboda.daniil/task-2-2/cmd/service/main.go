package main

import (
	"container/heap"
	"fmt"
)

type IntMinHeap []int

func (h IntMinHeap) Len() int {
	return len(h)
}
func (h IntMinHeap) Less(i, j int) bool {
	return h[i] < h[j]
}
func (h IntMinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
func (h *IntMinHeap) Push(x any) {
	*h = append(*h, x.(int))
}
func (h *IntMinHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func main() {

	var number, kPosition int

	if _, err := fmt.Scan(&number); err != nil || (number < 1 && number > 10000) {
		fmt.Println("Invalid number dishes")
		return
	}
	ratings := make([]int, number)
	maxDish := 0
	for i := 0; i < number; i++ {
		if _, err := fmt.Scan(&ratings[i]); err != nil {
			fmt.Println("Invalid number dishes")
			return
		}
		if ratings[i] > maxDish {
			maxDish = ratings[i]
		}
	}
	if _, err := fmt.Scan(&kPosition); err != nil || (kPosition < -10000 && kPosition > 10000) {
		fmt.Println("Invalid K number dish")
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
