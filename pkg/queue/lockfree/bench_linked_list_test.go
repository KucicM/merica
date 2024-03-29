package lockfree

import (
	"testing"

	"github.com/KucicM/merica/pkg/queue"
)


func BenchmarkLockFreeLinkedListSequential(b *testing.B) {
	q := NewLinkedListQueue[int]()
	for _, c := range queue.GetBenchmarkTable() {
		b.Run(c.Name(), func(b *testing.B) {
			queue.BenchmarkQueueSequential(b, q, c.NumOfWriters, c.NumOfReaders)
		})
	}
}

func BenchmarkLockFreeLinkedListParallel(b *testing.B) {
	q := NewLinkedListQueue[int]()
	for _, c := range queue.GetBenchmarkTable() {
		b.Run(c.Name(), func(b *testing.B) {
			queue.BenchmarkQueueParallel(b, q, c.NumOfWriters, c.NumOfReaders)
		})
	}
}