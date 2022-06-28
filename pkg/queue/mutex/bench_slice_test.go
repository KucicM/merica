
package mutex

import (
	"testing"

	"github.com/KucicM/merica/pkg/queue"
)


func BenchmarkMutexSliceSequential(b *testing.B) {
	q := NewSliceQueue[int]()
	for _, c := range queue.GetBenchmarkTable() {
		b.Run(c.Name(), func(b *testing.B) {
			queue.BenchmarkQueueSequential(b, q, c.NumOfWriters, c.NumOfReaders)
		})
	}
}

func BenchmarkMutexSliceParallel(b *testing.B) {
	q := NewSliceQueue[int]()
	for _, c := range queue.GetBenchmarkTable() {
		b.Run(c.Name(), func(b *testing.B) {
			queue.BenchmarkQueueParallel(b, q, c.NumOfWriters, c.NumOfReaders)
		})
	}
}