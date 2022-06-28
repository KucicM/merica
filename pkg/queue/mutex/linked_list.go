package mutex

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type LinkedListQueue[T any] struct {
	front *node[T]
	back *node[T]
	frontLock *sync.Mutex
	backLock *sync.Mutex
}

type node[T any] struct {
	val T
	next *node[T]
}

// without this, go detects race because initial state 
// back and front have same next... 
func (n *node[T]) atomicLoadNext() *node[T] {
	return (*node[T])(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&n.next))))
}

// same...
func (n *node[T]) atomicStoreNext(node *node[T]) {
	addr := (*unsafe.Pointer)(unsafe.Pointer(&n.next))
	atomic.StorePointer(addr, unsafe.Pointer(node))
}

func NewLinkedListQueue[T any]() *LinkedListQueue[T] {
	h := &node[T]{}
	return &LinkedListQueue[T]{h, h, &sync.Mutex{}, &sync.Mutex{}}
}

func (q *LinkedListQueue[T]) Enqueue(element T) {
	q.backLock.Lock()
	defer q.backLock.Unlock()

	n := &node[T]{element, nil}
	q.back.atomicStoreNext(n)
	q.back = n
}

func (q *LinkedListQueue[T]) Dequeue() (T, bool) {
	q.frontLock.Lock()
	defer q.frontLock.Unlock()

	var element T
	if q.front.atomicLoadNext() == nil {
		return element, false
	}


	element = q.front.next.val
	q.front = q.front.next // GC will clean memeroy

	return element, true
}