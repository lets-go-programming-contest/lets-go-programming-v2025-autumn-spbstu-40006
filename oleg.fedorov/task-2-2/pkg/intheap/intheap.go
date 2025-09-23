package intheap

import "fmt"

type IntHeap []int

func (h *IntHeap) Len() int           { return len(*h) }
func (h *IntHeap) Less(i, j int) bool { return (*h)[i] < (*h)[j] }
func (h *IntHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *IntHeap) Push(x interface{}) {
	num, ok := x.(int)
	if !ok {
		fmt.Printf("non-int value to IntHeap: %v\n", x)

		return
	}

	*h = append(*h, num)
}

func (h *IntHeap) Pop() interface{} {
	old := *h

	n := len(old)
	if n == 0 {
		return nil
	}

	x := old[n-1]

	*h = old[0 : n-1]

	return x
}

func (h *IntHeap) Peek() (int, bool) {
	if len(*h) == 0 {
		return 0, false
	}

	return (*h)[0], true
}
