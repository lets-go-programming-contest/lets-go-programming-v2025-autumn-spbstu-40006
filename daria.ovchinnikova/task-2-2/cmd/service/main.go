package main

import (
	"container/heap"
	"fmt"
)

type MaxHeap []int

func (heap MaxHeap) Len() int {
	return len(heap)
}

func (heap MaxHeap) Less(i, j int) bool {
	return heap[i] > heap[j]
}

func (heap MaxHeap) Swap(i, j int) {
	heap[i], heap[j] = heap[j], heap[i]
}

func (heap *MaxHeap) Push(x interface{}) {
	*heap = append(*heap, x.(int))
}

func (heap *MaxHeap) Pop() interface{} {
	old := *heap
	n := len(old)
	x := old[n-1]
	*heap = old[0 : n-1]
	return x
}

func main() {
	var numberOfDishes, kNumber, choice int
	_, err := fmt.Scan(&numberOfDishes)
	if err != nil {
		fmt.Println("Invalid number of dishes")
		return
	}

	ratings := make([]int, numberOfDishes)
	for i := 0; i < numberOfDishes; i++ {
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

	for i := 0; i < kNumber; i++ {
		choice = heap.Pop(ratingsHeap).(int)
	}

	fmt.Println(choice)

}
