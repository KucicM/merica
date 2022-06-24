// utils for tests
package queue

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
)


func QueueBasicTest(q Queue[int]) error {

	q.Enqueue(3)
	q.Enqueue(1)
	q.Enqueue(2)

	val, ok := q.Dequeue()
	if !ok {
		return fmt.Errorf("no value")
	}

	if val != 3 {
		return fmt.Errorf("expected 3 got %d", val)
	}

	val, ok = q.Dequeue()
	if !ok {
		return fmt.Errorf("no value")
	}

	if val != 1 {
		return fmt.Errorf("expected 1 got %d", val)
	}

	val, ok = q.Dequeue()
	if !ok {
		return fmt.Errorf("no value")
	}

	if val != 2 {
		return fmt.Errorf("expected 2 got %d", val)
	}

	val, ok = q.Dequeue()
	if ok {
		return fmt.Errorf("expected no value got %d", val)
	}

	return nil
}

func QueueRandomOpsTest(q Queue[int]) error {
	testSize := 100_000
	queueSize := 0
	lastDeque := -1
	lastEque := -1

	for i := 0; i < testSize; i++ {
		switch rand.Intn(2) {
		case 0:
			lastEque++
			queueSize++
			q.Enqueue(lastEque)
		case 1:
			val, ok := q.Dequeue()
			if queueSize == 0 && ok {
				return fmt.Errorf("expected empty queue got val %v", val)
			}

			if queueSize != 0 {
				lastDeque++
				if val != lastDeque {
					return fmt.Errorf("expected to get %d got %d", lastDeque, val)
				}
				queueSize--
			}
		}
	}

	return nil
}

// single writer and single reader
// tests order as well
func QueueConcurrentReadWriteTest(q Queue[int]) error {
	testSize := 500_000

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
					err = fmt.Errorf("order not ok, expected %d got %d", nextExpected, element)
					return
				}
				nextExpected++
			}
		}
		if nextExpected != testSize {
			err = fmt.Errorf("expected Enqueue count to be %d but it was %d", testSize, nextExpected)
		}
	}()

	writeWg.Wait()
	readWg.Wait()

	return err
}

// multiple writers and readers
// does not test order
func QueueConcurrentReadsWritesTest(q Queue[int]) error {
	testSize := 500_000

	numberOfWriters := 50
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

	recivedElements := make([]bool, testSize)
	lock := sync.Mutex{}


	readWg := sync.WaitGroup{}
	numberOfReaders := 50
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
			return fmt.Errorf("did not recive value %d", i)
		}
	}

	return nil
}