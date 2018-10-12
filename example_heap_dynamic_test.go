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

func (queue *ValuePriorityQueue) Pop() (int, bool) {
	if queue.queue.Pop() {
		return queue.Slice[queue.queue.Length()], true
	}
	return 0, false
}

func (queue *ValuePriorityQueue) Top() (int, bool) {
	if queue.queue.IsEmpty() {
		return 0, false
	}
	return queue.Slice[0], true
}

func (queue *ValuePriorityQueue) Update(value int, index int) bool {
	if index >= 0 && index < queue.queue.Length() {
		queue.Slice[index] = value
		queue.queue.Fix(index)
		return true
	}
	return false
}

func ExamplePriorityQueue_dynamicPush() {

	// slice is the underlying storage space for values
	// queue is initialized with a slice [5, 9, 12]
	queue := NewValuePriorityQueue(5, 9, 12)
	fmt.Println(queue.Slice)
	// now, the queue will look like
	//    12
	//   9  5

	// firstly, we push the rest of the values into the queue
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
	// so in every round, the largest number in the queue
	// is put at the end of the queue
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
