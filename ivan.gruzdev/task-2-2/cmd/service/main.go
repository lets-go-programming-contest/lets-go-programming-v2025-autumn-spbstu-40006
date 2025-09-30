package main

import (
	"container/heap"
	"fmt"
	"os"
)

type PriorityQueue []int

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i] < pq[j]
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(int))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq

	n := len(old)

	lastElement := old[n-1]

	*pq = old[:n-1]

	return lastElement

}

func mustScan(a ...interface{}) {
	_, err := fmt.Scan(a...)
	if err != nil {
		os.Exit(1)
	}
}

func main() {
	var (
		countDish, preference int
		pq                    PriorityQueue
	)

	heap.Init(&pq)
	mustScan(&countDish)

	rating := make([]int, countDish)

	for i := 0; i < countDish; i++ {
		mustScan(&rating[i])
	}

	mustScan(&preference)

	for i := 0; i < countDish; i++ {
		if pq.Len() < preference {
			heap.Push(&pq, rating[i])
		} else if pq[0] < rating[i] {
			heap.Pop(&pq)
			heap.Push(&pq, rating[i])
		}
	}

	fmt.Println(heap.Pop(&pq))

}
