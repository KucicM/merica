package queue

import (
	"sync"
)

const minSize int = 8

type RingQueue[T any] struct {
	ringBuffer []T
	head int
	tail int
	size int
	m     *sync.Mutex
}

func NewRingQueue[T any]() *RingQueue[T] {
	return &RingQueue[T]{
		ringBuffer: make([]T, minSize),
		m:     &sync.Mutex{},
	}
}

func (q *RingQueue[T]) Enqueue(element T) {
	q.m.Lock()
	defer q.m.Unlock()

	if q.size == len(q.ringBuffer) {
		q.resize()
	}

	q.ringBuffer[q.tail] = element
	q.tail = (q.tail + 1) & (len(q.ringBuffer) - 1)
	q.size++
}

func (q *RingQueue[T]) Dequeue() (T, bool) {
	q.m.Lock()
	defer q.m.Unlock()

	var element T
	if q.size <= 0 {
		return element, false
	}

	element = q.ringBuffer[q.head]
	var null T
	q.ringBuffer[q.head] = null

	q.head = (q.head + 1) & (len(q.ringBuffer) - 1)
	q.size--

	if q.size > minSize && (q.size<<2) == len(q.ringBuffer) {
		q.resize()
	}


	return element, true
}

func (q *RingQueue[T]) resize() {
		newRingBuffer := make([]T, q.size<<1)

		if q.tail > q.head {
			copy(newRingBuffer, q.ringBuffer[q.head:q.tail])
		} else {
			n := copy(newRingBuffer, q.ringBuffer[q.head:])
			copy(newRingBuffer[n:], q.ringBuffer[:q.tail])
		}

		q.head = 0
		q.tail = q.size
		q.ringBuffer = newRingBuffer
}
