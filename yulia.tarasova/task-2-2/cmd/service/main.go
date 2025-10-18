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
	var dishesNumb, kPreferDish, result int

	_, err := fmt.Scanln(&dishesNumb)
	if err != nil {
		fmt.Println("Invalid input")

		return
	}

	dishList := make([]int, dishesNumb)
	for i := range dishesNumb {
		_, err = fmt.Scan(&dishList[i])
		if err != nil {
			fmt.Println("Invalid input")

			return
		}
	}

	_, err = fmt.Scan(&kPreferDish)
	if err != nil {
		fmt.Println("Invalid input")

		return
	}

	dishesHeap := &IntHeap{}
	heap.Init(dishesHeap)

	for i := range dishesNumb {
		heap.Push(dishesHeap, dishList[i])
	}

	for i := 0; i < kPreferDish; i++ {
		result = heap.Pop(dishesHeap).(int)
	}

	fmt.Println(result)

}
