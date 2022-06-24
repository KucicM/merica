package lockfree

import (
	"sync/atomic"
	"unsafe"
)


type LinkedListQueue[T any] struct {
	front unsafe.Pointer
	back unsafe.Pointer
}

type node[T any] struct {
	val T
	next unsafe.Pointer
}

func NewLinkedListQueue[T any]() *LinkedListQueue[T] {
	n := unsafe.Pointer(new(node[T]))
	return &LinkedListQueue[T]{n, n}
}

func (q *LinkedListQueue[T]) Enqueue(element T) {
	new_node := unsafe.Pointer(&node[T]{element, nil})
	var old_back unsafe.Pointer

	// when compare and swap succedes in one thread all other threads are blocked until 
	// the back is updated
	for {
		old_back = atomic.LoadPointer(&q.back)
		if atomic.CompareAndSwapPointer(&(*node[T])(old_back).next, nil, new_node) {
			break
		}
	}
	atomic.CompareAndSwapPointer(&q.back, old_back, new_node)
}


func (q *LinkedListQueue[T]) Dequeue() (T, bool) {
	var old_front, next unsafe.Pointer
	for {
		old_front = atomic.LoadPointer(&q.front)
		next = atomic.LoadPointer(&((*node[T])(old_front).next))
		if next == nil || atomic.CompareAndSwapPointer(&q.front, old_front, next) {
			break
		}
	}

	var e T
	var ok bool
	if next != nil {
		e = (*node[T])(next).val
		ok = true
	}
	return e, ok
}
