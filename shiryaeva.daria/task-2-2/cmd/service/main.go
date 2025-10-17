package main

import (
	"container/heap"
	"fmt"
)

type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

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
	var n, k int

	_, err := fmt.Scan(&n)
	if err != nil {

		return
	}

	ratings := make([]int, n)
	for i := 0; i < n; i++ {
		_, err = fmt.Scan(&ratings[i])
		if err != nil {

			return
		}
	}

	_, err = fmt.Scan(&k)
	if err != nil {

		return
	}

	h := &IntHeap{}
	heap.Init(h)

	for _, rating := range ratings {
		heap.Push(h, rating)
	}

	for i := 0; i < k-1; i++ {
		heap.Pop(h)
	}

	result := heap.Pop(h).(int)
	fmt.Println(result)
}
