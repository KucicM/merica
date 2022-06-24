package mutex_test

import (
	"testing"

	"github.com/KucicM/merica/pkg/queue"
	"github.com/KucicM/merica/pkg/queue/mutex"
)

func TestSliceQueueBasic(t *testing.T) {
	q := mutex.NewSliceQueue[int]()
	err := queue.QueueBasicTest(q)
	if err != nil {
		t.Error(err)
	}
}

func TestSliceQueueRandomOps(t *testing.T) {
	q := mutex.NewSliceQueue[int]()
	err := queue.QueueRandomOpsTest(q)
	if err != nil {
		t.Error(err)
	}
}

func TestSliceQueueSingleConcurreny(t *testing.T) {
	q := mutex.NewSliceQueue[int]()
	err := queue.QueueConcurrentReadWriteTest(q)
	if err != nil {
		t.Error(err)
	}
}

func TestSliceQueueConcurreny(t *testing.T) {
	q := mutex.NewSliceQueue[int]()
	err := queue.QueueConcurrentReadsWritesTest(q)
	if err != nil {
		t.Error(err)
	}
}
