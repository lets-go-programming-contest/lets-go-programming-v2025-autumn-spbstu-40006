package main

import (
	"container/heap"
	"fmt"

	"github.com/monka6/task-2-2/intheap"
)

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

	ratingsHeap := &intheap.MaxHeap{}
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
