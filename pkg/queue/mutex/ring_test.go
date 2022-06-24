package mutex_test

import (
	"testing"

	"github.com/KucicM/merica/pkg/queue"
	"github.com/KucicM/merica/pkg/queue/mutex"
)

func TestRingQueueBasic(t *testing.T) {
	q := mutex.NewRingQueue[int]()
	err := queue.QueueBasicTest(q)
	if err != nil {
		t.Error(err)
	}
}

func TestRingQueueRandomOps(t *testing.T) {
	q := mutex.NewRingQueue[int]()
	err := queue.QueueRandomOpsTest(q)
	if err != nil {
		t.Error(err)
	}
}

func TestRingQueueSingleConcurreny(t *testing.T) {
	q := mutex.NewRingQueue[int]()
	err := queue.QueueConcurrentReadWriteTest(q)
	if err != nil {
		t.Error(err)
	}
}

func TestRingQueueConcurreny(t *testing.T) {
	q := mutex.NewRingQueue[int]()
	err := queue.QueueConcurrentReadsWritesTest(q)
	if err != nil {
		t.Error(err)
	}
}
