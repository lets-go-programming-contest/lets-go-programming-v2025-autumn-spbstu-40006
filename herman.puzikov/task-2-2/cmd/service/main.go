package main

import (
	"container/heap"
	"fmt"
	"os"

	"github.com/Segfault-chan/task-2-2/pkg/intheap"
)

func main() {
	var dishesNum int
	if _, err := fmt.Scan(&dishesNum); err != nil {
		fmt.Fprintln(os.Stderr, "failed to read N:", err)

		return
	}

	values := make([]int, dishesNum)
	for i := range dishesNum {
		if _, err := fmt.Scan(&values[i]); err != nil {
			fmt.Fprintln(os.Stderr, "failed to read ai:", err)

			return
		}
	}

	var indexK int
	if _, err := fmt.Scan(&indexK); err != nil {
		fmt.Fprintln(os.Stderr, "failed to read k:", err)

		return
	}

	if indexK < 1 || indexK > dishesNum {
		fmt.Fprintln(os.Stderr, "k must be in [1..N]")

		return
	}

	intHeap := &intheap.IntHeap{}
	heap.Init(intHeap)

	for i := range dishesNum {
		currentValue := values[i]
		if intHeap.Len() < indexK {
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
