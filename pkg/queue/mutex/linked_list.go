package mutex

import "sync"

type LinkedListQueue[T any] struct {
	front *node[T]
	back *node[T]
	m *sync.Mutex
}

type node[T any] struct {
	val T
	next *node[T]
}

func NewLinkedListQueue[T any]() *LinkedListQueue[T] {
	h := &node[T]{}
	return &LinkedListQueue[T]{h, h, &sync.Mutex{}}
}

func (q *LinkedListQueue[T]) Enqueue(element T) {
	q.m.Lock()
	defer q.m.Unlock()

	n := &node[T]{element, nil}
	q.back.next = n
	q.back = n
}

func (q *LinkedListQueue[T]) Dequeue() (T, bool) {
	q.m.Lock()
	defer q.m.Unlock()

	var element T
	if q.front.next == nil {
		return element, false
	}


	element = q.front.next.val
	q.front = q.front.next // GC will clean memeroy

	return element, true
}