package main

import (
	"container/heap"
	"fmt"

	"github.com/dizey5k/task-2-2/internal/utils"
	"github.com/dizey5k/task-2-2/pkg/intheap"
)

func main() {
	ratings, kthNumber, err := utils.ReadRatingsAndK()
	if err != nil {
		fmt.Printf("error reading: %v\n", err)

		return
	}

	if len(ratings) == 0 {
		fmt.Println("0")

		return
	}

	highest := &intheap.IntHeap{}
	heap.Init(highest)

	for _, rating := range ratings {
		if highest.Len() < kthNumber {
			heap.Push(highest, rating)

			continue
		}

		if minimum, ok := highest.Peek(); ok && rating > minimum {
			heap.Pop(highest)
			heap.Push(highest, rating)
		}
	}

	if result, ok := highest.Peek(); ok {
		fmt.Println(result)
	} else {
		fmt.Println("0")
	}
}
