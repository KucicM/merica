package lockfree_test

import (
	"testing"

	"github.com/KucicM/merica/pkg/queue"
	"github.com/KucicM/merica/pkg/queue/lockfree"
)

func TestLockFreeQueueBasic(t *testing.T) {
	q := lockfree.NewLinkedListQueue[int]()
	err := queue.QueueBasicTest(q)
	if err != nil {
		t.Error(err)
	}
}

func TestLockFreeQueueRandomOps(t *testing.T) {
	q := lockfree.NewLinkedListQueue[int]()
	err := queue.QueueRandomOpsTest(q)
	if err != nil {
		t.Error(err)
	}
}

func TestLockFreeQueueSingleConcurreny(t *testing.T) {
	q := lockfree.NewLinkedListQueue[int]()
	err := queue.QueueConcurrentReadWriteTest(q)
	if err != nil {
		t.Error(err)
	}
}

func TestLockFreeQueueConcurreny(t *testing.T) {
	q := lockfree.NewLinkedListQueue[int]()
	err := queue.QueueConcurrentReadsWritesTest(q)
	if err != nil {
		t.Error(err)
	}
}
