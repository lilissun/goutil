package goutil_test

import (
	"fmt"

	"github.com/lilissun/goutil"
)

// This example demonstrate how to use Priority Queue
// to sort elements

func Example_priorityQueue_sort() {
	// slice is the underlying storage space for values
	// queue is initialized with slice
	// with half of the values included, i.e., [5, 9, 12]
	slice := []int{5, 9, 12, 0, 3, 7, 4}
	queue := goutil.NewPriorityQueue(
		slice, len(slice)/2,
		func(i, j int) bool { return slice[i] > slice[j] },
	)
	fmt.Println(slice)
	// now, the queue will look like
	//    12
	//   9  5

	// firstly, we push the rest of the values into the queue
	// i.e.,[0, 3, 7, 4]
	for queue.Push() {
	}
	fmt.Println(slice)
	// now, the queue will look like
	//     12
	//   9    7
	// 0  3  5  4

	// then, we pop all the elements out of the queue
	// so in every round, the largest number in the queue
	// is put at the end of the queue
	for queue.Pop() {
	}
	fmt.Println(slice)
	// now, the queue is empty and the slice is sorted

	// Output:
	// [12 9 5 0 3 7 4]
	// [12 9 7 0 3 5 4]
	// [0 3 4 5 7 9 12]
}
