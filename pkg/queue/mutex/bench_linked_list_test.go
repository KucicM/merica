
package mutex

import (
	"testing"

	"github.com/KucicM/merica/pkg/queue"
)


func BenchmarkMutexLinkedListSequential(b *testing.B) {
	q := NewLinkedListQueue[int]()
	for _, c := range queue.GetBenchmarkTable() {
		b.Run(c.Name(), func(b *testing.B) {
			queue.BenchmarkQueueSequential(b, q, c.NumOfWriters, c.NumOfReaders)
		})
	}
}

func BenchmarkMutexLinkedListParallel(b *testing.B) {
	q := NewLinkedListQueue[int]()
	for _, c := range queue.GetBenchmarkTable() {
		b.Run(c.Name(), func(b *testing.B) {
			queue.BenchmarkQueueParallel(b, q, c.NumOfWriters, c.NumOfReaders)
		})
	}
}