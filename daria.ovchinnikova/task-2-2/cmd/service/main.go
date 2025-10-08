package main

import (
	"container/heap"
	"fmt"
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
	var numberOfDishes, kNumber int
	_, err := fmt.Scan(&numberOfDishes)
	if err != nil {
		fmt.Println("Invalid number of dishes")

		return
	}

	ratings := make([]int, numberOfDishes)
	for i := range numberOfDishes {
		_, err = fmt.Scan(&ratings[i])
		if err != nil {
			fmt.Println("Invalid rating")

			return
		}
	}

	_, err = fmt.Scan(&kNumber)
	if err != nil {
		fmt.Println("Invalid k-number")

		return
	}

	ratingsHeap := &MaxHeap{}
	heap.Init(ratingsHeap)

	for _, rating := range ratings {
		heap.Push(ratingsHeap, rating)
	}

	for range kNumber - 1 {
		heap.Pop(ratingsHeap)
	}

	choice := heap.Pop(ratingsHeap)
	fmt.Println(choice)

}
