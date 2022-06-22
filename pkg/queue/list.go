package queue

import (
	"container/list"
	"fmt"
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

func (q *ListQueue[T]) Dequeue() (T, error) {
	q.m.Lock()
	defer q.m.Unlock()

	var element T
	if q.queue.Len() <= 0 {
		return element, fmt.Errorf("Queue is empty")
	}

	e := q.queue.Front()
	element = e.Value.(T)
	q.queue.Remove(e)

	return element, nil
}

func (q *ListQueue[T]) Size() int {
	q.m.Lock()
	defer q.m.Unlock()
	return q.queue.Len()
}
