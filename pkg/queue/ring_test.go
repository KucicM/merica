package queue_test

import (
	"fmt"
	"testing"

	"github.com/KucicM/merica/pkg/queue"
)

func TestRingBasic(t *testing.T) {
	q := queue.NewRingQueue[int]()
	err := QueueBasicTest(q)
	if err != nil {
		t.Error(err)
	}
}

func TestRingRandomOps(t *testing.T) {
	q := queue.NewRingQueue[int]()
	err := QueueRandomOpsTest(q)
	if err != nil {
		t.Error(err)
	}
}

func TestRingSingleConcurreny(t *testing.T) {
	q := queue.NewRingQueue[int]()
	err := QueueConcurrentReadWriteTest(q)
	if err != nil {
		t.Error(err)
	}
}

func TestRingConcurreny(t *testing.T) {
	q := queue.NewRingQueue[int]()
	err := QueueConcurrentReadsWritesTest(q)
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkRingSingleThread(b *testing.B) {
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