package waitfree_test

import (
	"testing"

	"github.com/KucicM/merica/pkg/queue"
	"github.com/KucicM/merica/pkg/queue/waitfree"
)

func TestWaitFreeQueueBasic(t *testing.T) {
	q := waitfree.NewLinkedListQueue[int]()
	err := queue.QueueBasicTest(q)
	if err != nil {
		t.Error(err)
	}
}

func TestWaitFreeQueueRandomOps(t *testing.T) {
	q := waitfree.NewLinkedListQueue[int]()
	err := queue.QueueRandomOpsTest(q)
	if err != nil {
		t.Error(err)
	}
}

func TestWaitFreeQueueSingleConcurreny(t *testing.T) {
	q := waitfree.NewLinkedListQueue[int]()
	err := queue.QueueConcurrentReadWriteTest(q)
	if err != nil {
		t.Error(err)
	}
}

func TestWaitFreeQueueConcurreny(t *testing.T) {
	q := waitfree.NewLinkedListQueue[int]()
	err := queue.QueueConcurrentReadsWritesTest(q)
	if err != nil {
		t.Error(err)
	}
}
