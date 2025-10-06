package main

import (
	"container/heap"
	"fmt"
)

type IntHeap []int

func (h *IntHeap) Len() int           { return len(*h) }
func (h *IntHeap) Less(i, j int) bool { return (*h)[i] > (*h)[j] }
func (h *IntHeap) Swap(i, j int)      { (*h)[i], (*h)[j] = (*h)[j], (*h)[i] }

func (h *IntHeap) Push(x any) {
	v, ok := x.(int)
	if !ok {
		return
	}

	*h = append(*h, v)
}

func (h *IntHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]

	return x
}

func readN() (int, error) {
	var n int
	_, err := fmt.Scan(&n)
	if err != nil {
		return 0, err
	}
	if n < 1 || n > 10000 {
		return 0, fmt.Errorf("n out of range")
	}
	return n, nil
}

func readRatings(n int) ([]int, error) {
	ratings := make([]int, n)
	for idx := range n {
		_, err := fmt.Scan(&ratings[idx])
		if err != nil {
			return nil, err
		}
		if ratings[idx] < -10000 || ratings[idx] > 10000 {
			return nil, fmt.Errorf("rating out of range")
		}
	}
	return ratings, nil
}

func readK(n int) (int, error) {
	var k int
	_, err := fmt.Scan(&k)
	if err != nil {
		return 0, err
	}
	if k < 1 || k > n {
		return 0, fmt.Errorf("k out of range")
	}
	return k, nil
}

func buildMaxHeap(a []int) *IntHeap {
	h := &IntHeap{}
	heap.Init(h)
	for _, v := range a {
		heap.Push(h, v)
	}
	return h
}

func kthMax(h *IntHeap, k int) (int, error) {
	for range k - 1 {
		heap.Pop(h)
	}
	v, ok := heap.Pop(h).(int)
	if !ok {
		return 0, fmt.Errorf("unexpected type")
	}
	return v, nil
}

func main() {
	n, err := readN()
	if err != nil {
		return
	}

	ratings, err := readRatings(n)
	if err != nil {
		return
	}

	k, err := readK(n)
	if err != nil {
		return
	}

	h := buildMaxHeap(ratings)
	ans, err := kthMax(h, k)
	if err != nil {
		return
	}
	fmt.Println(ans)
}
