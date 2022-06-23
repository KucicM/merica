package queue_test

import (
	"fmt"
	"testing"

	"github.com/KucicM/merica/pkg/queue"
)

func TestBasic(t *testing.T) {
	q := queue.NewSliceQueue[int]()
	err := QueueBasicTest(q)
	if err != nil {
		t.Error(err)
	}
}

func TestRandomOps(t *testing.T) {
	q := queue.NewSliceQueue[int]()
	err := QueueRandomOpsTest(q)
	if err != nil {
		t.Error(err)
	}
}

func TestSingleConcurreny(t *testing.T) {
	q := queue.NewSliceQueue[int]()
	err := QueueConcurrentReadWriteTest(q)
	if err != nil {
		t.Error(err)
	}
}

func TestConcurreny(t *testing.T) {
	q := queue.NewSliceQueue[int]()
	err := QueueConcurrentReadsWritesTest(q)
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkSingleThread(b *testing.B) {
	for n := 10; n < 100_000_000; n *= 10 {
		b.Run(fmt.Sprintf("batch_size_%d", n), func (b *testing.B)  {
			q := queue.NewSliceQueue[int]()
			for i := 0; i < b.N; i++ {
				q.Enqueue(i)
			}

			for i := 0; i < b.N; i++ {
				q.Dequeue()
			}
		})
	}
}