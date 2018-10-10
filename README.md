# Utilities for Go Programming Language

This project aims to improve the user experience of some native functions and objects,
e.g., Heap, File Scanner.

## Heap

Heap is implemented following the `sort.Slice(...)` interface.
It takes a slice as its full storage space and a length as its currently filled length.
It also accepts a comparison function less which is operated on the slice.
Notice that the slice cannot be appended in size as it will create a new slice object,
which voids the less comparison function.

Officially, the Priority Queue is implemented as follows.

```go
// An Item is something we manage in a priority queue.
type Item struct {
    value    string // The value of the item; arbitrary.
    priority int    // The priority of the item in the queue.
    // The index is needed by update and is maintained by the heap.Interface methods.
    index int // The index of the item in the heap.
}

// A PriorityQueue implements heap.Interface and holds Items.
type PriorityQueue []*Item

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
    // We want Pop to give us the highest, not lowest, priority so we use greater than here.
    return pq[i].priority > pq[j].priority
}

func (pq PriorityQueue) Swap(i, j int) {
    pq[i], pq[j] = pq[j], pq[i]
    pq[i].index = i
    pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
    n := len(*pq)
    item := x.(*Item)
    item.index = n
    *pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
    old := *pq
    n := len(old)
    item := old[n-1]
    item.index = -1 // for safety
    *pq = old[0 : n-1]
    return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update(item *Item, value string, priority int) {
    item.value = value
    item.priority = priority
    heap.Fix(pq, item.index)
}
```

It is better than the official implementation in the sense that
the users do not need to define a custom type with `heap.Interface`.
Meanwhile, the comparison function and its targets is defined
based on the actual type held in `slice`.
Thus, we do not need to rewrite the whole `PriorityQueue`
repeatedly for different element types.

On the other hand, this `PriorityQueue` implementation is not self-contained.
When `queue.Push()` is used, it pushes the element
located at `queue.GetLength()` in `slice` to the queue.
Similarly, when `queue.Pop()` is used, it pops the first element out of queue
and puts the element to the position `queue.GetLength()` in `slice`.
Moreover, its memory space must be re-allocated before the initialization of the queue.

Heap must be used with lock in multi-threading context.