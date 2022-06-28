
package mutex

import (
	"testing"

	"github.com/KucicM/merica/pkg/queue"
)


func BenchmarkMutexRingSequential(b *testing.B) {
	q := NewRingQueue[int]()
	for _, c := range queue.GetBenchmarkTable() {
		b.Run(c.Name(), func(b *testing.B) {
			queue.BenchmarkQueueSequential(b, q, c.NumOfWriters, c.NumOfReaders)
		})
	}
}

func BenchmarkMutexRingParallel(b *testing.B) {
	q := NewRingQueue[int]()
	for _, c := range queue.GetBenchmarkTable() {
		b.Run(c.Name(), func(b *testing.B) {
			queue.BenchmarkQueueParallel(b, q, c.NumOfWriters, c.NumOfReaders)
		})
	}
}