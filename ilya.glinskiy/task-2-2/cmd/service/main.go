package main

import (
	"container/heap"
	"fmt"
)

type IntHeap []int

func (h IntHeap) Len() int {
	return len(h)
}

func (h IntHeap) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h IntHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *IntHeap) Push(x interface{}) {
	*h = append(*h, x.(int))
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func main() {
	var (
		dishAmount, num int
		err             error
	)

	var dishHeap *IntHeap = &IntHeap{}
	heap.Init(dishHeap)

	_, err = fmt.Scan(&dishAmount)
	if err != nil {
		fmt.Println("Invalid dish amount")

		return
	}

	for range dishAmount {
		_, err = fmt.Scan(&num)
		if err != nil {
			fmt.Println("Invalid dish rating")

			return
		}

		heap.Push(dishHeap, num)
	}

	_, err = fmt.Scan(&num)
	if err != nil || num > dishHeap.Len() {
		fmt.Println("Invalid dish index")

		return
	}

	for dishHeap.Len() > num {
		heap.Pop(dishHeap)
	}

	fmt.Println(heap.Pop(dishHeap))
}
