# Utilities for Go Programming Language

This project aims to improve the user experience
of some native functions and structures in Go standard libraries.

## Heap

Heap is implemented following the `sort.Slice(...)` interface.
It takes a slice as its full storage space and a length as its currently filled length.
It also accepts a comparison function less which is operated on the slice.
Notice that the slice cannot be appended in size as it will create a new slice object,
which voids the less comparison function.

Officially,
the [Priority Queue](https://golang.org/pkg/container/heap/#example__priorityQueue)
is implemented as follows.

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

// Update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) Update(item *Item, value string, priority int) {
    item.value = value
    item.priority = priority
    heap.Fix(pq, item.index)
}
```

Our `PriorityQueue` implementation is better than the official one in the sense that
the users do not need to define a custom type with `heap.Interface`.
Meanwhile, the comparison function and its targets is defined
based on the actual type held in `slice`.
Thus, we do not need to rewrite the whole `PriorityQueue`
repeatedly for different element types.
For example, the following code heap-sort `documents` in accending order.

```go
type Document struct {
    ID    int64
    Score float64
}

func NewDocument(id int64, score float64) *Document {
    return &Document{
        ID:    id,
        Score: score,
    }
}

queue := NewPriorityQueue(
    documents, len(documents),
    func(i, j int) bool {
        return documents[i].Score > documents[j].Score
    },
)
for queue.Pop() {
}
```

On the other hand, our `PriorityQueue` implementation is not self-contained.
The memory space of `slice` must be pre-allocated before the initialization of the queue.
It also means that `slice` cannot be dynamically extended to a larger size after the queue is built.
(In this case, it is recommended to create a new `PriorityQueue` with the extended `slice`.)
When `queue.Push()` is used, it pushes the element
located at `queue.GetLength()` in `slice` to the queue
rather than taking an element as parameters.
Similarly, when `queue.Pop()` is used, it pops the first element out of queue
and puts the element to the position `queue.GetLength()` in `slice`, 
rather than returning the poped element.
When an element located at `index` is updated in `slice`,
the user need to invoke `queue.Fix(index)` to update the queue.
For example, the following code finds the document with kth largest score.
(The document with kth largest score is held at `documents[0]`)

```go
queue := NewPriorityQueue(
    documents[:k], k,
    func(i, j int) bool {
        return documents[i].Score < documents[j].Score
    },
)
for _, document := range documents[k:] {
    if document.Score > documents[0].Score {
        documents[0] = document
        queue.Fix(0)
    }
}
```

Heap must be used with lock in multi-threading context.