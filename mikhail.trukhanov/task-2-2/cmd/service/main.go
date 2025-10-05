package main

import (
	"container/heap"
	"fmt"

	"github.com/Mishaa105/task-2-2/pkg/intheap"
)

func main() {
	var amountOfDishes, dishNumber, num int

	dishesHeap := intheap.New()

	_, err := fmt.Scan(&amountOfDishes)
	if err != nil {
		fmt.Println("Invalid input")

		return
	}

	for n := 0; n != amountOfDishes; n++ {
		_, err := fmt.Scan(&num)
		if err != nil {
			fmt.Println("Invalid input")

			return
		}

		heap.Push(dishesHeap, num)
	}

	_, err = fmt.Scan(&dishNumber)
	if err != nil {
		fmt.Println("Invalid input")

		return
	}

	for n := 0; n != dishNumber; n++ {
		if n == dishNumber-1 {
			fmt.Println((*dishesHeap)[0])
		} else {
			heap.Pop(dishesHeap)
		}
	}
}
