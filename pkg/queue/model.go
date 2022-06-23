package queue

type Queue[T any] interface {
	Enqueue(T)
	Dequeue() (T, bool)
	Size() int
}