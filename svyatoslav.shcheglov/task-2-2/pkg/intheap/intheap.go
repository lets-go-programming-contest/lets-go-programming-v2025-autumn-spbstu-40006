package intheap

import "fmt"

type IntHeap []int

func (h *IntHeap) Len() int {
	return len(*h)
}

func (h *IntHeap) Less(i, j int) bool {
	return (*h)[i] > (*h)[j]
}

func (h *IntHeap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *IntHeap) Push(x any) {
	value, ok := x.(int)
	if !ok {
		fmt.Printf("IntHeap: expected int type, got %T\n", x)

		return
	}

	*h = append(*h, value)
}

func (h *IntHeap) Pop() any {
	oldHeap := *h
	heapSize := len(oldHeap)
	topElement := oldHeap[heapSize-1]
	*h = oldHeap[:heapSize-1]

	return topElement
}
