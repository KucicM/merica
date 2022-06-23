package queue

import (
	"sync"
)

type SliceQueue[T any] struct {
	queue []T
	m *sync.Mutex
}

func NewSliceQueue[T any]() *SliceQueue[T] {
	return &SliceQueue[T]{
		queue: make([]T, 0),
		m: &sync.Mutex{},
	}
}

func (q *SliceQueue[T]) Enqueue(element T) {
	q.m.Lock()
	defer q.m.Unlock()
	q.queue = append(q.queue, element)
}

func (q *SliceQueue[T]) Dequeue() (T, bool) {
	q.m.Lock()
	defer q.m.Unlock()

	var element T
	if len(q.queue) <= 0 {
		return element, false
	}

	element = q.queue[0]
	q.queue = q.queue[1:]

	return element, true
}

func (q *SliceQueue[T]) Size() int {
	q.m.Lock()
	defer q.m.Unlock()
	return len(q.queue)
}
