
package queue_test

import (
	"fmt"
	"testing"

	"github.com/KucicM/merica/pkg/queue"
)

func TestCustumListBasic(t *testing.T) {
	q := queue.NewCustomListQueue[int]()
	err := QueueBasicTest(q)
	if err != nil {
		t.Error(err)
	}
}

func TestCustomListRandomOps(t *testing.T) {
	q := queue.NewCustomListQueue[int]()
	err := QueueRandomOpsTest(q)
	if err != nil {
		t.Error(err)
	}
}

func TestCustomListSingleConcurreny(t *testing.T) {
	q := queue.NewCustomListQueue[int]()
	err := QueueConcurrentReadWriteTest(q)
	if err != nil {
		t.Error(err)
	}
}

func TestCustomListConcurreny(t *testing.T) {
	q := queue.NewCustomListQueue[int]()
	err := QueueConcurrentReadsWritesTest(q)
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkCustomListSingleThread(b *testing.B) {
	for n := 10; n < 100_000_000; n *= 10 {
		b.Run(fmt.Sprintf("batch_size_%d", n), func (b *testing.B)  {
			q := queue.NewCustomListQueue[int]()
			for i := 0; i < b.N; i++ {
				q.Enqueue(i)
			}

			for i := 0; i < b.N; i++ {
				q.Dequeue()
			}
		})
	}
}