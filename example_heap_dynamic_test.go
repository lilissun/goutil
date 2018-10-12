// This example demonstrate how to use Priority queue
// with dynamic size of underlying storage space
package goutil_test

import (
	"fmt"

	"github.com/lilissun/goutil"
)

// In order to enable dynamic size of storage space,
// we need to define a type to hold the slice and the queue
type ValuePriorityQueue struct {
	Slice []int
	queue *goutil.PriorityQueue
}

// NewValuePriorityQueue init the queue with slice as its initial storage
func NewValuePriorityQueue(slice ...int) *ValuePriorityQueue {
	queue := &ValuePriorityQueue{Slice: slice}
	queue.queue = goutil.NewPriorityQueue(
		queue.Slice, len(queue.Slice),
		func(i, j int) bool { return queue.Slice[i] > queue.Slice[j] },
	)
	return queue
}

// when an element is pushed to the queue
// we need to ensure that it is not full
// otherwise, we must enlarge its capacity
// by allocating a new slice with double its size
// notice that the slice and queue are updated at the same time
func (queue *ValuePriorityQueue) Push(value int) {
	if queue.queue.IsFull() {
		slice := queue.Slice
		length := queue.queue.Length()
		queue.Slice = make([]int, 2*length+1)
		copy(queue.Slice, slice)
		queue.Slice[length] = value
		queue.queue = goutil.NewPriorityQueue(
			queue.Slice, length+1,
			func(i, j int) bool { return queue.Slice[i] > queue.Slice[j] },
		)
		return
	}
	queue.Slice[queue.queue.Length()] = value
	queue.queue.Push()
}

// when a value is poped, we need to ensure that it is not empty
func (queue *ValuePriorityQueue) Pop() (int, bool) {
	if queue.queue.Pop() {
		return queue.Slice[queue.queue.Length()], true
	}
	return 0, false
}

// when the top value is enquired, we need to ensure it is not empty
func (queue *ValuePriorityQueue) Top() (int, bool) {
	if queue.queue.IsEmpty() {
		return 0, false
	}
	return queue.Slice[0], true
}

// when a value is updated at the index,
// we additionally ask the queue to fix the index
func (queue *ValuePriorityQueue) Update(value int, index int) bool {
	if index >= 0 && index < queue.queue.Length() {
		queue.Slice[index] = value
		queue.queue.Fix(index)
		return true
	}
	return false
}

func ExamplePriorityQueue_dynamicSize() {
	// queue is initialized with a slice [5, 9, 12]
	queue := NewValuePriorityQueue(5, 9, 12)
	fmt.Println(queue.Slice)
	// now, the queue will look like
	//    12
	//   9  5

	// then, we push more values into the queue
	// i.e.,[0, 3, 7, 4]
	for _, value := range []int{0, 3, 7, 4} {
		queue.Push(value)
	}
	fmt.Println(queue.Slice)
	// now, the queue will look like
	//     12
	//   9    7
	// 0  3  5  4

	// then, we pop all the elements out of the queue
	// in a descending order
	slice := make([]int, 0, len(queue.Slice))
	for {
		value, more := queue.Pop()
		if more == false {
			break
		}
		slice = append(slice, value)
	}
	fmt.Println(slice)
	// now, the queue is empty and the slice is sorted

	// Output:
	// [12 9 5]
	// [12 9 7 0 3 5 4]
	// [12 9 7 5 4 3 0]
}
