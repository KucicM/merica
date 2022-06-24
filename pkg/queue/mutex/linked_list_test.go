package mutex_test

import (
	"testing"

	"github.com/KucicM/merica/pkg/queue"
	"github.com/KucicM/merica/pkg/queue/mutex"
)

func TestCustumListQueueBasic(t *testing.T) {
	q := mutex.NewLinkedListQueue[int]()
	err := queue.QueueBasicTest(q)
	if err != nil {
		t.Error(err)
	}
}

func TestCustomListQueueRandomOps(t *testing.T) {
	q := mutex.NewLinkedListQueue[int]()
	err := queue.QueueRandomOpsTest(q)
	if err != nil {
		t.Error(err)
	}
}

func TestCustomListQueueSingleConcurreny(t *testing.T) {
	q := mutex.NewLinkedListQueue[int]()
	err := queue.QueueConcurrentReadWriteTest(q)
	if err != nil {
		t.Error(err)
	}
}

func TestCustomListQueueConcurreny(t *testing.T) {
	q := mutex.NewLinkedListQueue[int]()
	err := queue.QueueConcurrentReadsWritesTest(q)
	if err != nil {
		t.Error(err)
	}
}