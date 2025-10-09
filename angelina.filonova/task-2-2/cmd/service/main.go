package main

import (
	"container/heap"
	"fmt"

	"github.com/filon6/task-2-2/pkg/intheap"
)

func main() {
	var countDishes, favorite int
	fmt.Scan(&countDishes)

	arr := make([]int, countDishes)
	for i := range countDishes {
		fmt.Scan(&arr[i])
	}

	fmt.Scan(&favorite)

	h := &intheap.IntHeap{}
	heap.Init(h)

	for i := range countDishes {
		term := arr[i]
		if h.Len() < favorite {
			heap.Push(h, term)
		} else if term > (*h)[0] {
			heap.Pop(h)
			heap.Push(h, term)
		}
	}

	fmt.Println((*h)[0])
}
