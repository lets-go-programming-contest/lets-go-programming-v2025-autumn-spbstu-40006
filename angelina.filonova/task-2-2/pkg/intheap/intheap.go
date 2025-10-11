package intheap

type IntHeap []int

func (h IntHeap) Len() int {
	return len(h)
}

func (h IntHeap) Less(i, j int) bool {
	return h[i] < h[j]
}

func (h IntHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *IntHeap) Push(x any) {
	value, ok := x.(int)
	if !ok {
		panic("IntHeap: invalid type")
	}

	*h = append(*h, value)
}

func (h *IntHeap) Pop() any {
	oldHeap := *h
	n := len(oldHeap)
	lastElement := oldHeap[n-1]
	*h = oldHeap[0 : n-1]

	return lastElement
}
