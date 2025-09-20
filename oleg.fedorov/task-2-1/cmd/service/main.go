package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type IntHeap []int

func (h IntHeap) Len() int           { return len(h) }
func (h IntHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntHeap) Push(x interface{}) {
	if num, ok := x.(int); ok {
		*h = append(*h, num)
	}
}

func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]

	return x
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	numOfDishes, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return
	}
	if numOfDishes < 1 || numOfDishes > 10000 {
		return
	}

	scanner.Scan()
	rText := scanner.Text()
	ratingStrs := strings.Fields(rText)
	if len(ratingStrs) != numOfDishes {
		return
	}

	ratings := make([]int, numOfDishes)
	for i, str := range ratingStrs {
		rating, err := strconv.Atoi(str)
		if err != nil {
			return
		}
		if rating < -10000 || rating > 10000 {
			return
		}
		ratings[i] = rating
	}

	scanner.Scan()
	kthNumber, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return
	}
	if kthNumber < 1 || kthNumber > numOfDishes {
		return
	}

	intHeap := &IntHeap{}
	heap.Init(intHeap)

	for _, rating := range ratings {
		if intHeap.Len() < kthNumber {
			heap.Push(intHeap, rating)
		} else if rating > (*intHeap)[0] {
			heap.Pop(intHeap)
			heap.Push(intHeap, rating)
		}
	}

	result := (*intHeap)[0]
	fmt.Println(result)
}
