package queue

import "sync"

type CustomListQueue[T any] struct {
	front *node[T]
	rear *node[T]
	m *sync.Mutex
}

type node[T any] struct {
	val T
	next *node[T]
}

func NewCustomListQueue[T any]() *CustomListQueue[T] {
	return &CustomListQueue[T]{nil, nil, &sync.Mutex{}}
}

func (q *CustomListQueue[T]) Enqueue(element T) {
	q.m.Lock()
	defer q.m.Unlock()

	n := &node[T]{element, nil}

	if q.rear == nil {
		q.rear = n
		q.front = n
		return
	}

	q.rear.next = n
	q.rear = n
}

func (q *CustomListQueue[T]) Dequeue() (T, bool) {
	q.m.Lock()
	defer q.m.Unlock()

	var element T
	if q.front == nil {
		return element, false
	}

	element = q.front.val
	q.front = q.front.next // GC will clean memeroy

	if q.front == nil {
		q.rear = nil
	}

	return element, true
}