package queue

import (
	"fmt"
	"sync"
)

type ArrayMutexQueue[T any] struct {
	queue []T
	m *sync.Mutex
}

func NewArrayMutexQueue[T any]() *ArrayMutexQueue[T] {
	return &ArrayMutexQueue[T]{
		queue: make([]T, 0),
		m: &sync.Mutex{},
	}
}

func (q *ArrayMutexQueue[T]) Enqueue(element T) {
	q.m.Lock()
	defer q.m.Unlock()
	q.queue = append(q.queue, element)
}

func (q *ArrayMutexQueue[T]) Dequeue() (T, error) {
	q.m.Lock()
	defer q.m.Unlock()

	var element T
	if len(q.queue) <= 0 {
		return element, fmt.Errorf("Queue is empty")
	}

	element = q.queue[0]
	q.queue = q.queue[1:]

	return element, nil
}

func (q *ArrayMutexQueue[T]) Size() int {
	q.m.Lock()
	defer q.m.Unlock()
	return len(q.queue)
}
