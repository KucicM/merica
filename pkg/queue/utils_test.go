package queue_test

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"

	"github.com/KucicM/merica/pkg/queue"
)


func QueueBasicTest(q queue.Queue[int]) error {

	q.Enqueue(3)
	size := q.Size()
	if size != 1 {
		return fmt.Errorf("Queue not size 1 but %d", size)
	}

	q.Enqueue(1)
	size = q.Size()
	if size != 2 {
		return fmt.Errorf("Queue not size 2 but %d", size)
	}

	q.Enqueue(2)
	size = q.Size()
	if size != 3 {
		return fmt.Errorf("Queue not size 3 but %d", size)
	}

	val, ok := q.Dequeue()
	if !ok {
		return fmt.Errorf("No value")
	}

	if val != 3 {
		return fmt.Errorf("Expected 3 got %d", val)
	}

	val, ok = q.Dequeue()
	if !ok {
		return fmt.Errorf("No value")
	}

	if val != 1 {
		return fmt.Errorf("Expected 1 got %d", val)
	}

	val, ok = q.Dequeue()
	if !ok {
		return fmt.Errorf("No value")
	}

	if val != 2 {
		return fmt.Errorf("Expected 2 got %d", val)
	}

	size = q.Size()
	if size != 0 {
		return fmt.Errorf("Expected empty queue got size %d", size)
	}

	val, ok = q.Dequeue()
	if ok {
		return fmt.Errorf("Expected no value got %d", val)
	}

	return nil
}

func QueueRandomOpsTest(q queue.Queue[int]) error {
	testSize := 100_000
	queueSize := 0
	lastDeque := -1
	lastEque := -1

	for i := 0; i < testSize; i++ {
		switch rand.Intn(3) {
		case 0: 
			size := q.Size()
			if size != queueSize {
				return fmt.Errorf("Expected size of %d got size %d", queueSize, size)
			}
		case 1:
			lastEque++
			queueSize++
			q.Enqueue(lastEque)
		case 2:
			val, ok := q.Dequeue()
			if queueSize == 0 && ok {
				return fmt.Errorf("Expected empty queue got val %v", val)
			}

			if queueSize != 0 {
				lastDeque++
				if val != lastDeque {
					return fmt.Errorf("Expected to get %d got %d", lastDeque, val)
				}
				queueSize--
			}
		}
	}

	return nil
}

// single writer and single reader
// tests order as well
func QueueConcurrentReadWriteTest(q queue.Queue[int]) error {
	testSize := 100_000

	writeWg := sync.WaitGroup{}
	writeWg.Add(1)
	var writeDone int32

	go func() {
		defer writeWg.Done()
		for i := 0; i < testSize; i++ {
			q.Enqueue(i)
		}
		atomic.AddInt32(&writeDone, 1)
	}()

	go func() {
		// just to cover race check
		q.Size() 
	}()

	var err error
	readWg := sync.WaitGroup{}
	readWg.Add(1)

	go func() {
		defer readWg.Done()
		nextExpected := 0
		ok := true
		var element int
		for  ok || atomic.LoadInt32(&writeDone) == 0 {
			element, ok = q.Dequeue()
			if ok {
				if element != nextExpected {
					err = fmt.Errorf("Order not ok, expected %d got %d", nextExpected, element)
					return
				}
				nextExpected++
			}
		}
		if nextExpected != testSize {
			err = fmt.Errorf("Expected Enqueue count to be %d but it was %d", testSize, nextExpected)
		}
	}()

	writeWg.Wait()
	readWg.Wait()

	return err
}

// multiple writers and readers
// does not test order
func QueueConcurrentReadsWritesTest(q queue.Queue[int]) error {
	testSize := 100_000

	numberOfWriters := 10
	writeWg := sync.WaitGroup{}

	for i := 0; i < numberOfWriters; i++ {
		writeWg.Add(1)
		go func(wId int) {
			defer writeWg.Done()
			for j := wId; j < testSize; j++ {
				if j % numberOfWriters == wId {
					q.Enqueue(j)
				}
			}
		}(i)
	}

	var writeDone int32

	go func() {
		writeWg.Wait()
		atomic.AddInt32(&writeDone, 1)
	}()

	for i := 0; i < 10; i++ {
		go func() {
			q.Size()
		}()
	}

	recivedElements := make([]bool, testSize)
	lock := sync.Mutex{}


	readWg := sync.WaitGroup{}
	numberOfReaders := 10
	for i := 0; i < numberOfReaders; i++ {
		readWg.Add(1)
		go func() {
			defer readWg.Done()
			res := make([]int, 0)

			ok := true
			var element int
			for  ok || atomic.LoadInt32(&writeDone) == 0 {
				if element, ok = q.Dequeue(); ok {
					res = append(res, element)
				}
			}

			lock.Lock()
			defer lock.Unlock()
			for _, id := range res {
				recivedElements[id] = true
			}
		}()
	}

	readWg.Wait()
	for i, v := range recivedElements {
		if !v {
			return fmt.Errorf("Did not recive value %d", i)
		}
	}

	return nil
}