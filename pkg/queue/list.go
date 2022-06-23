package queue

import (
	"container/list"
	"sync"
)

type ListQueue[T any] struct {
	queue *list.List
	m     *sync.Mutex
}

func NewListQueue[T any]() *ListQueue[T] {
	return &ListQueue[T]{
		queue: list.New(),
		m:     &sync.Mutex{},
	}
}

func (q *ListQueue[T]) Enqueue(element T) {
	q.m.Lock()
	defer q.m.Unlock()
	q.queue.PushBack(element)
}

func (q *ListQueue[T]) Dequeue() (T, bool) {
	q.m.Lock()
	defer q.m.Unlock()

	var element T
	if q.queue.Len() <= 0 {
		return element, false
	}

	e := q.queue.Front()
	element = e.Value.(T)
	q.queue.Remove(e)

	return element, true
}
