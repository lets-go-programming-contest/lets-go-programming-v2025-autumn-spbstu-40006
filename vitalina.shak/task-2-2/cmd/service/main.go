package main

import (
	"container/heap"
	"errors"
	"fmt"
	"os"
)

const (
	minDishesCount = 1
	maxDishesCount = 10000

	minScore = -10000
	maxScore = 10000
)

var (
	errReadDishesCnt = errors.New("failed to read dishes count")
	errReadScore     = errors.New("failed to read dish score")
	errReadPref      = errors.New("failed to read preference")
	errInvalidDishes = errors.New("dishes count must be in range [1-10000]")
	errInvalidScore  = errors.New("dish score must be in range [-10000-10000]")
	errInvalidPref   = errors.New("preference must be in range [1-N]")
)

type scoreHeap []int

func (h scoreHeap) Len() int { return len(h) }
func (h scoreHeap) Less(i, k int) bool { return h[i] < h[k] }
func (h scoreHeap) Swap(i, k int) { h[i], h[k] = h[k], h[i] }

func (h *scoreHeap) Push(value any) {
	score, ok := value.(int)
	if !ok {

		return
	}

	*h = append(*h, score)
}

func (h *scoreHeap) Pop() any {
	old := *h
	last := old[len(old)-1]
	*h = old[:len(old)-1]

	return last
}

func main() {
	var dishesCount int
	if _, err := fmt.Scan(&dishesCount); err != nil {
		fmt.Fprintln(os.Stderr, errReadDishesCnt.Error())

		return
	}

	if dishesCount < minDishesCount || dishesCount > maxDishesCount {
		fmt.Fprintln(os.Stderr, errInvalidDishes.Error())

		return
	}

	scores := make([]int, dishesCount)
	for index := 0; index < dishesCount; index++ {
		if _, err := fmt.Scan(&scores[index]); err != nil {
			fmt.Fprintln(os.Stderr, errReadScore.Error())

			return
		}

		if scores[index] < minScore || scores[index] > maxScore {
			fmt.Fprintln(os.Stderr, errInvalidScore.Error())

			return
		}
	}

	var prefIdx int
	if _, err := fmt.Scan(&prefIdx); err != nil {
		fmt.Fprintln(os.Stderr, errReadPref.Error())

		return
	}

	if prefIdx < minDishesCount || prefIdx > dishesCount {
		fmt.Fprintln(os.Stderr, errInvalidPref.Error())

		return
	}

	result := selectBestPref(scores, prefIdx)
	fmt.Println(result)
}

func selectBestPref(scores []int, prefIdx int) int {
	topPref := &scoreHeap{}
	heap.Init(topPref)

	for _, score := range scores {
		heap.Push(topPref, score)

		if topPref.Len() > prefIdx {
			heap.Pop(topPref)
		}
	}

	return (*topPref)[0]
}
