package queue_test

import (
	// "fmt"
	"fmt"
	"testing"

	"github.com/KucicM/merica/pkg/queue"
)

func TestLockFreeBasic(t *testing.T) {
	q := queue.NewLockFreeQueue[int]()
	err := QueueBasicTest(q)
	if err != nil {
		t.Error(err)
	}
}

func TestLockFreeRandomOps(t *testing.T) {
	q := queue.NewLockFreeQueue[int]()
	err := QueueRandomOpsTest(q)
	if err != nil {
		t.Error(err)
	}
}

func TestLockFreeSingleConcurreny(t *testing.T) {
	q := queue.NewLockFreeQueue[int]()
	err := QueueConcurrentReadWriteTest(q)
	if err != nil {
		t.Error(err)
	}
}

func TestLockFreeConcurreny(t *testing.T) {
	q := queue.NewLockFreeQueue[int]()
	err := QueueConcurrentReadsWritesTest(q)
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkLockFreeSingleThread(b *testing.B) {
	for n := 10; n < 100_000_000; n *= 10 {
		b.Run(fmt.Sprintf("batch_size_%d", n), func(b *testing.B) {
			q := queue.NewLockFreeQueue[int]()
			for i := 0; i < b.N; i++ {
				q.Enqueue(i)
			}

			for i := 0; i < b.N; i++ {
				q.Dequeue()
			}
		})
	}
}
