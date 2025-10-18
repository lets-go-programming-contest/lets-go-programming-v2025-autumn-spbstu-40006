package main

import (
	"container/heap"
	"fmt"
)

type PriorityQueue []int

func (pq *PriorityQueue) Len() int {
	return len(*pq)
}

func (pq *PriorityQueue) Less(i, j int) bool {
	return (*pq)[i] < (*pq)[j]
}

func (pq *PriorityQueue) Swap(i, j int) {
	(*pq)[i], (*pq)[j] = (*pq)[j], (*pq)[i]
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

func mustScan(a ...interface{}) error {
	if _, err := fmt.Scan(a...); err != nil {
		return fmt.Errorf("error reading the input: %w", err)
	}

	return nil
}

func main() {
	var (
		countDish, preference int
		priorityQueue         PriorityQueue
	)

	heap.Init(&priorityQueue)

	if err := mustScan(&countDish); err != nil {
		fmt.Printf("Error reading dish count: %v\n", err)

		return
	}

	rating := make([]int, countDish)

	for i := range countDish {
		if err := mustScan(&rating[i]); err != nil {
			fmt.Printf("Error reading raiting[%d]: %v\n", i, err)

			return
		}
	}

	if err := mustScan(&preference); err != nil {
		fmt.Printf("Error reading preference: %v\n", err)

		return
	}

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
