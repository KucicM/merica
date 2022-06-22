package queue_test

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/KucicM/merica/pkg/queue"
)

func TestBasic(t *testing.T) {
	q := queue.NewArrayMutexQueue[int]()

	q.Enqueue(3)
	size := q.Size()
	if size != 1 {
		t.Errorf("Queue not size 1 but %d", size)
	}

	q.Enqueue(1)
	size = q.Size()
	if size != 2 {
		t.Errorf("Queue not size 2 but %d", size)
	}

	q.Enqueue(2)
	size = q.Size()
	if size != 3 {
		t.Errorf("Queue not size 3 but %d", size)
	}

	val, err := q.Dequeue()
	if err != nil {
		t.Error(err)
	}

	if val != 3 {
		t.Errorf("Expected 3 got %d", val)
	}

	val, err = q.Dequeue()
	if err != nil {
		t.Error(err)
	}

	if val != 1 {
		t.Errorf("Expected 1 got %d", val)
	}

	val, err = q.Dequeue()
	if err != nil {
		t.Error(err)
	}

	if val != 2 {
		t.Errorf("Expected 2 got %d", val)
	}

	size = q.Size()
	if size != 0 {
		t.Errorf("Expected empty queue got size %d", size)
	}

	_, err = q.Dequeue()
	if err == nil {
		t.Error("Expected error got nil")
	}
}

func TestConcurreny(t *testing.T) {
	q := queue.NewArrayMutexQueue[int]()

	testSize := 10
	wgW := sync.WaitGroup{}

	writerSize := 20
	for i := 0; i < writerSize; i++ {
		wgW.Add(1)
		go func(id int) {
			defer wgW.Done()
			for j := 0; j < testSize; j++ {
				if j % writerSize == id {
					q.Enqueue(j)
				}
			}
		}(i)
	}

	var writeDone int32
	go func(){
		wgW.Wait()
		atomic.AddInt32(&writeDone, 1)
	}()

	check := make([]bool, testSize)

	wgR := sync.WaitGroup{}
	for i := 0; i < 20; i++ {
		wgR.Add(1)
		go func() {
			defer wgR.Done()
			var err error
			var element int
			for err == nil || atomic.LoadInt32(&writeDone) == 0 {
				element, err = q.Dequeue()
				if err == nil {
					check[element] = true
				}
			}
		}()
	}

	wgR.Wait()

	for i, val := range check {
		if !val {
			t.Fatalf("Did not recive %d", i)
		}
	}

}

func BenchmarkSingleThread(b *testing.B) {
	for n := 10; n < 100_000_000; n *= 10 {
		b.Run(fmt.Sprintf("batch_size_%d", n), func (b *testing.B)  {
			q := queue.NewArrayMutexQueue[int]()
			for i := 0; i < b.N; i++ {
				q.Size()
				q.Enqueue(i)
			}

			for i := 0; i < b.N; i++ {
				q.Size()
				q.Dequeue()
			}
		})
	}
}