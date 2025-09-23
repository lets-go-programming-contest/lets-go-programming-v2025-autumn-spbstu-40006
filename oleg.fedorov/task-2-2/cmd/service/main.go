package main

import (
	"container/heap"
	"fmt"

	"github.com/dizey5k/task-2-2/internal/utils"
	"github.com/dizey5k/task-2-2/pkg/intheap"
)

func main() {
	ratings, kthNumber := utils.ReadRatingsAndK()
	if ratings == nil {
		return
	}

	if len(ratings) == 0 {
		fmt.Println("0")
		return
	}

	h := &intheap.IntHeap{}
	heap.Init(h)

	for _, rating := range ratings {
		if h.Len() < kthNumber {
			heap.Push(h, rating)
			continue
		}

		if rating > h.Peek() {
			heap.Pop(h)
			heap.Push(h, rating)
		}
	}

	result := h.Peek()
	fmt.Println(result)
}
