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
	value, ok := x.(int)
	if !ok {
		return
	}
	*pq = append(*pq, value)
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
		priorityQueue         PriorityQueue
	)

	heap.Init(&priorityQueue)
	mustScan(&countDish)

	rating := make([]int, countDish)

	for i := 0; i < countDish; i++ {
		mustScan(&rating[i])
	}

	mustScan(&preference)

	for i := range countDish {
		if priorityQueue.Len() < preference {
			heap.Push(&priorityQueue, rating[i])
		} else if priorityQueue[0] < rating[i] {
			heap.Pop(&priorityQueue)
			heap.Push(&priorityQueue, rating[i])
		}
	}

	fmt.Println(heap.Pop(&priorityQueue))

}
