package main

import (
	"container/heap"
	"fmt"
)

type IntHeap []int

func (h *IntHeap) Len() int {
	return len(*h)
}

func (h *IntHeap) Less(i, j int) bool {
	return (*h)[i] > (*h)[j]
}

func (h *IntHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *IntHeap) Push(x interface{}) {
	value, ok := x.(int)
	if !ok {
		return
	}

	*h = append(*h, value)
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]

	return x
}

func main() {
	var dishAmount, dishNum int

	_, err := fmt.Scan(&dishAmount)
	if err != nil {
		fmt.Println("Incorrect amount of dishes")

		return
	}

	if dishAmount <= 0 {
		fmt.Println("Amount of dishes must be positive number")

		return
	}

	ratings := make([]int, dishAmount)
	for i := range ratings {
		_, err = fmt.Scan(&ratings[i])
		if err != nil {
			fmt.Println("Incorrect raiting")

			return
		}
	}

	_, err = fmt.Scan(&dishNum)
	if err != nil {
		fmt.Println("Incorrect number of dish")

		return
	}

	if dishNum <= 0 || dishNum > dishAmount {
		fmt.Println("Number of dish must be in [1...N]")

		return
	}

	dishHeap := &IntHeap{}
	heap.Init(dishHeap)

	for _, rating := range ratings {
		heap.Push(dishHeap, rating)
	}

	for i := 1; i < dishNum; i++ {
		heap.Pop(dishHeap)
	}

	fmt.Println(heap.Pop(dishHeap))
}
