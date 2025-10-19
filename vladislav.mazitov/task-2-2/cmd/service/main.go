package main

import (
	"container/heap"
	"fmt"

	"github.com/identicalaffiliation/task-2-2/pkg/intheap"
)

func newHeap() *intheap.Tree {
	h := &intheap.Tree{}
	heap.Init(h)

	return h
}

func main() {
	var countDishes, numbers int

	_, err := fmt.Scan(&countDishes)
	if err != nil {
		return
	}

	ratings := make([]int, countDishes)
	for index := range countDishes {
		_, err := fmt.Scan(&ratings[index])
		if err != nil {
			return
		}
	}

	_, err = fmt.Scan(&numbers)
	if err != nil {
		return
	}

	ratingsHeap := newHeap()

	for _, rating := range ratings {
		heap.Push(ratingsHeap, rating)
	}

	for range numbers - 1 {
		heap.Pop(ratingsHeap)
	}

	result, ok := heap.Pop(ratingsHeap).(int)
	if !ok {
		return
	}

	fmt.Println(result)
}
