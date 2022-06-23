
package queue_test

import (
	"fmt"
	"testing"

	"github.com/KucicM/merica/pkg/queue"
)

func TestListBasic(t *testing.T) {
	q := queue.NewListQueue[int]()
	err := QueueBasicTest(q)
	if err != nil {
		t.Error(err)
	}
}

func TestListRandomOps(t *testing.T) {
	q := queue.NewListQueue[int]()
	err := QueueRandomOpsTest(q)
	if err != nil {
		t.Error(err)
	}
}

func TestListSingleConcurreny(t *testing.T) {
	q := queue.NewListQueue[int]()
	err := QueueConcurrentReadWriteTest(q)
	if err != nil {
		t.Error(err)
	}
}

func TestListConcurreny(t *testing.T) {
	q := queue.NewListQueue[int]()
	err := QueueConcurrentReadsWritesTest(q)
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkListSingleThread(b *testing.B) {
	for n := 10; n < 100_000_000; n *= 10 {
		b.Run(fmt.Sprintf("batch_size_%d", n), func (b *testing.B)  {
			q := queue.NewListQueue[int]()
			for i := 0; i < b.N; i++ {
				q.Enqueue(i)
			}

			for i := 0; i < b.N; i++ {
				q.Dequeue()
			}
		})
	}
}