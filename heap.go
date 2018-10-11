package goutil

import (
	"reflect"
)

// PriorityQueue defines a priority queue with pre-allocated memory
type PriorityQueue struct {
	slice    interface{}
	length   int
	capacity int
	less     func(i, j int) bool
	swap     func(i, j int)
}

// NewPriorityQueue creates a priority queue
// with pre-allocated memory space at slice.
// The slice is initially filled with elements of length.
// The maximum capacity of the priority queue is limited
// to the current length of the slice.
// The elements can be compared with the less function
func NewPriorityQueue(
	slice interface{},
	length int,
	less func(i, j int) bool,
) *PriorityQueue {
	value := reflect.ValueOf(slice)
	capacity := value.Len()
	if length > capacity {
		length = capacity
	}
	swap := reflect.Swapper(slice)
	queue := &PriorityQueue{
		slice:    slice,
		length:   length,
		capacity: capacity,
		less:     less,
		swap:     swap,
	}
	for index := queue.length/2 - 1; index >= 0; index-- {
		queue.down(index)
	}
	return queue
}

func (queue *PriorityQueue) down(begin int) bool {
	index := begin
	for {
		left := 2*index + 1
		if left >= queue.length || left < 0 {
			break
		}
		small := left
		right := left + 1
		if right < queue.length && queue.less(right, left) {
			small = right
		}
		if queue.less(small, index) == false {
			break
		}
		queue.swap(small, index)
		index = small
	}
	return index > begin
}

func (queue *PriorityQueue) up(index int) {
	for {
		parent := (index - 1) / 2
		if parent == index || queue.less(index, parent) == false {
			break
		}
		queue.swap(parent, index)
		index = parent
	}
}

// IsEmpty checks if the queue is empty
func (queue *PriorityQueue) IsEmpty() bool {
	return queue.length == 0
}

// IsFull checks if the queue is full
func (queue *PriorityQueue) IsFull() bool {
	return queue.length == queue.capacity
}

// Length of the queue
func (queue *PriorityQueue) Length() int {
	return queue.length
}

// Capacity of the queue
func (queue *PriorityQueue) Capacity() int {
	return queue.capacity
}

// Pop the top element from the queue
// and put it to the current end
func (queue *PriorityQueue) Pop() bool {
	if queue.IsEmpty() {
		return false
	}
	queue.length--
	queue.swap(0, queue.length)
	queue.down(0)
	return true
}

// Push an element located at the end of the queue
func (queue *PriorityQueue) Push() bool {
	if queue.IsFull() {
		return false
	}
	queue.length++
	queue.up(queue.length - 1)
	return true
}

// Fix the queue because the value at index is updated
func (queue *PriorityQueue) Fix(index int) {
	if queue.down(index) == false {
		queue.up(index)
	}
}
