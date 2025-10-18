package main

import (
	"container/heap"
	"fmt"

	"github.com/SpeaarIt/task-2-2/pkg/intheap"
)

func main() {
	var totalDishes, targetPosition int

	_, err := fmt.Scan(&totalDishes)
	if err != nil {
		fmt.Println("Incorrect amount of dishes")

		return
	}

	if totalDishes <= 0 {
		fmt.Println("Amount of dishes must be positive number")

		return
	}

	dishRatings := make([]int, totalDishes)
	for i := range dishRatings {
		_, err = fmt.Scan(&dishRatings[i])
		if err != nil {
			fmt.Println("Incorrect rating")

			return
		}
	}

	_, err = fmt.Scan(&targetPosition)
	if err != nil {
		fmt.Println("Incorrect number of dish")

		return
	}

	if targetPosition <= 0 || targetPosition > totalDishes {
		fmt.Println("Number of dish must be in [1...N]")

		return
	}

	ratingHeap := &intheap.IntHeap{}
	heap.Init(ratingHeap)

	for _, rating := range dishRatings {
		heap.Push(ratingHeap, rating)
	}

	for currentPosition := 1; currentPosition < targetPosition; currentPosition++ {
		heap.Pop(ratingHeap)
	}

	selectedDishRating := heap.Pop(ratingHeap)
	fmt.Println(selectedDishRating)
}
