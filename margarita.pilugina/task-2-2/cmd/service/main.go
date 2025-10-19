package main

import (
	"container/heap"
	"fmt"
	"os"

	"github.com/MargotBush/task-2-2/pkg/intheap"
)

func main() {
	var dishesNum int
	if _, err := fmt.Scan(&dishesNum); err != nil {
		fmt.Fprintln(os.Stderr, "failed to read dishes number:", err)

		return
	}

	values := make([]int, dishesNum)
	for i := range dishesNum {
		if _, err := fmt.Scan(&values[i]); err != nil {
			fmt.Fprintln(os.Stderr, "failed to read dish rating:", err)

			return
		}
	}

	var favIndex int
	if _, err := fmt.Scan(&favIndex); err != nil {
		fmt.Fprintln(os.Stderr, "failed to read favorite index:", err)

		return
	}

	if favIndex < 1 || favIndex > dishesNum {
		fmt.Fprintln(os.Stderr, "favorite index must be in [1..N]")

		return
	}

	dishHeap := &intheap.IntHeap{}
	heap.Init(dishHeap)

	for i := range dishesNum {
		current := values[i]

		if dishHeap.Len() < favIndex {
			heap.Push(dishHeap, current)

			continue
		}

		if current > (*dishHeap)[0] {
			heap.Pop(dishHeap)
			heap.Push(dishHeap, current)
		}
	}

	fmt.Println((*dishHeap)[0])
}
