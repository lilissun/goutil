package goutil_test

import (
	"fmt"

	"github.com/lilissun/goutil"
)

// This example demonstrate how to use Priority Queue
// to find the kth largest number in a list

func ExamplePriorityQueue_top_k() {
	// slice is the underlying storage space for values
	// queue is initialized with slice with size of k
	slice := []int{5, 9, 12, 0, 3, 7, 4}
	k := 3
	queue := goutil.NewPriorityQueue(
		slice, k,
		func(i, j int) bool { return slice[i] < slice[j] },
	)
	fmt.Println(slice)
	// now, the queue will look like
	//     5
	//   9  12

	// for the rest of values in the slice,
	// we compare them with the first (smallest) value in the slice
	// and replace the first value if they are larger
	for _, value := range slice[k:] {
		if value > slice[0] {
			slice[0] = value
			queue.Fix(0)
		}
	}
	fmt.Println(slice)
	// now, the queue will look like
	//     7
	//   9  12

	// lastly, we output the first value as the kth largest number in the list
	fmt.Println(slice[0])

	// Output:
	// [5 9 12 0 3 7 4]
	// [7 9 12 0 3 7 4]
	// 7
}
