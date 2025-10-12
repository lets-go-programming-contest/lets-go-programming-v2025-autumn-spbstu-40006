package intheap

import "container/heap"

type intheap []int

func New() *intheap {
	h := &intheap{}
	heap.Init(h)

	return h
}

func (h *intheap) Len() int {
	return len(*h)
}

func (h *intheap) Less(i, j int) bool {
	return (*h)[i] > (*h)[j]
}

func (h *intheap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *intheap) Push(x any) {
	val, ok := x.(int)
	if !ok {
		return
	}

	*h = append(*h, val)
}

func (h *intheap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}
