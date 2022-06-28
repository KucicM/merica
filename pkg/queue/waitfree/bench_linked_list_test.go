package waitfree

import (
	"testing"

	"github.com/KucicM/merica/pkg/queue"
)


func BenchmarkWaitFreeLinkedListSequential(b *testing.B) {
	q := NewLinkedListQueue[int]()
	for _, c := range queue.GetBenchmarkTable() {
		b.Run(c.Name(), func(b *testing.B) {
			queue.BenchmarkQueueSequential(b, q, c.NumOfWriters, c.NumOfReaders)
		})
	}
}

func BenchmarkWaitFreeLinkedListParallel(b *testing.B) {
	q := NewLinkedListQueue[int]()
	for _, c := range queue.GetBenchmarkTable() {
		b.Run(c.Name(), func(b *testing.B) {
			queue.BenchmarkQueueParallel(b, q, c.NumOfWriters, c.NumOfReaders)
		})
	}
}