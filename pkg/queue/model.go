package queue

type Queue[T any] interface {
	Enqueue(T)
	Dequeue() (T, error)
	Size() int
}