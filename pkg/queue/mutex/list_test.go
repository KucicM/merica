package mutex_test

import (
	"testing"

	"github.com/KucicM/merica/pkg/queue"
	"github.com/KucicM/merica/pkg/queue/mutex"
)

func TestListQueueBasic(t *testing.T) {
	q := mutex.NewListQueue[int]()
	err := queue.QueueBasicTest(q)
	if err != nil {
		t.Error(err)
	}
}

func TestListQueueRandomOps(t *testing.T) {
	q := mutex.NewListQueue[int]()
	err := queue.QueueRandomOpsTest(q)
	if err != nil {
		t.Error(err)
	}
}

func TestListQueueSingleConcurreny(t *testing.T) {
	q := mutex.NewListQueue[int]()
	err := queue.QueueConcurrentReadWriteTest(q)
	if err != nil {
		t.Error(err)
	}
}

func TestListQueueConcurreny(t *testing.T) {
	q := mutex.NewListQueue[int]()
	err := queue.QueueConcurrentReadsWritesTest(q)
	if err != nil {
		t.Error(err)
	}
}
