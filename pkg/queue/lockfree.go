package queue

import (
	"sync/atomic"
	"unsafe"
)


type LockFreeQueue[T any] struct {
	front unsafe.Pointer
	back unsafe.Pointer
}

type lockFreeNode[T any] struct {
	val T
	next unsafe.Pointer
}

func NewLockFreeQueue[T any]() *LockFreeQueue[T] {
	n := unsafe.Pointer(new(lockFreeNode[T]))
	return &LockFreeQueue[T]{n, n}
}

func (q *LockFreeQueue[T]) Enqueue(element T) {
	new_node := unsafe.Pointer(&lockFreeNode[T]{element, nil})
	var old_back unsafe.Pointer

	// when compare and swap succedes in one thread all other threads are blocked until 
	// the back is updated
	for {
		old_back = atomic.LoadPointer(&q.back)
		if atomic.CompareAndSwapPointer(&(*lockFreeNode[T])(old_back).next, nil, new_node) {
			break
		}
	}
	atomic.CompareAndSwapPointer(&q.back, old_back, new_node)
}


func (q *LockFreeQueue[T]) Dequeue() (T, bool) {
	var old_front, next unsafe.Pointer
	for {
		old_front = atomic.LoadPointer(&q.front)
		next = atomic.LoadPointer(&((*lockFreeNode[T])(old_front).next))
		if next == nil || atomic.CompareAndSwapPointer(&q.front, old_front, next) {
			break
		}
	}

	var e T
	var ok bool
	if next != nil {
		e = (*lockFreeNode[T])(next).val
		ok = true
	}
	return e, ok
}
