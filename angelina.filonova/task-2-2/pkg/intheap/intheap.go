package intheap

import "fmt"

type IntHeap []int

func (h *IntHeap) Len() int           { return len(*h) }
func (h *IntHeap) Less(i, j int) bool { return (*h)[i] < (*h)[j] }
func (h *IntHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *IntHeap) Push(x any) {
	value, ok := x.(int)
	if !ok {
		fmt.Printf("IntHeap: invalid type %T\n", x)

		return
	}

	*h = append(*h, value)
}

func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	last := old[n-1]
	*h = old[:n-1]

	return last
}
