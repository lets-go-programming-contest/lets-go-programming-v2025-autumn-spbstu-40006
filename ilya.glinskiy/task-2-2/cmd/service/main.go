package main

import (
	"container/heap"
	"fmt"

	"github.com/bloomkicks/task-2-2/internal/intheap"
)

func main() {
	var (
		dishAmount, num int
		err             error
	)

	dishHeap := &intheap.IntHeap{}

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
	if err != nil || num < 1 || num > dishHeap.Len() {
		fmt.Println("Invalid dish index")

		return
	}

	for dishHeap.Len() > num {
		heap.Pop(dishHeap)
	}

	fmt.Println(heap.Pop(dishHeap))
}
