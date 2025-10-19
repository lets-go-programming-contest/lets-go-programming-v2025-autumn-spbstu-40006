package intheap

type Tree []int

func (h *Tree) Len() int {
	return len(*h)
}

func (h *Tree) Less(i, j int) bool {
	return (*h)[i] > (*h)[j]
}

func (h *Tree) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *Tree) Push(x interface{}) {
	value, ok := x.(int)
	if ok {
		*h = append(*h, value)
	}
}

func (h *Tree) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}
