package main

import (
	"container/heap"
	"fmt"
	"os"

	"github.com/filon6/task-2-2/pkg/intheap"
)

func main() {
	var countDishes int
	if _, err := fmt.Scan(&countDishes); err != nil {
		fmt.Fprintln(os.Stderr, "error: failed to read countDishes:", err)

		return
	}

	dishValues := make([]int, countDishes)
	for i := range countDishes {
		if _, err := fmt.Scan(&dishValues[i]); err != nil {
			fmt.Fprintln(os.Stderr, "error: failed to read dish value:", err)

			return
		}
	}

	var favoriteIndex int
	if _, err := fmt.Scan(&favoriteIndex); err != nil {
		fmt.Fprintln(os.Stderr, "error: failed to read favorite index:", err)

		return
	}

	if favoriteIndex < 1 || favoriteIndex > countDishes {
		fmt.Fprintln(os.Stderr, "error: favoriteIndex must be positive")

		return
	}

	intHeap := &intheap.IntHeap{}
	heap.Init(intHeap)

	for i := range countDishes {
		currentValue := dishValues[i]
		if intHeap.Len() < favoriteIndex {
			heap.Push(intHeap, currentValue)

			continue
		}

		if currentValue > (*intHeap)[0] {
			heap.Pop(intHeap)
			heap.Push(intHeap, currentValue)
		}
	}

	fmt.Println((*intHeap)[0])
}
