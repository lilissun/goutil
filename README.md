# Utilities for Go Programming Language

This project aims to improve the user experience of some native functions and objects,
e.g., Heap, File Scanner.

## Heap

Heap is implemented following the `sort.Slice(...)` interface.
It takes a slice as its full storage space and a length as its currently filled length.
It also accepts a comparison function less which is operated on the slice.
Notice that the slice cannot be appended in size as it will create a new slice object,
which voids the less comparison function.

It is better than the official implementation in the sense that
the users do not need to define a custom type with `heap.Interface`.
Meanwhile, the comparison function is defined
based on the actual type held in the slice.
Thus, we can support all sorts of comparison funcions.

On the other hand, the structure is not self-contained.
When `queue.Push()` is used, it pushes the element
located at `queue.GetLength()` in `slice` to the queue.
Similarly, when `queue.Pop()` is used, it pops the first element out of queue
and set it to `queue.GetLength()`.
