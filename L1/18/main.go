package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type ICounter interface {
	Increment()
	GetValue() int32
}

func worker(counter ICounter, times int) {
	for range times {
		counter.Increment()
	}
}

// AtomicCounter implements ICounter using atomic operations
type AtomicCounter struct {
	value int32
}

func (ac *AtomicCounter) Increment() {
	atomic.AddInt32(&ac.value, 1)
}

func (ac *AtomicCounter) GetValue() int32 {
	return ac.value
}

// MutexCounter implements ICounter using RWMutex
type MutexCounter struct {
	value int32
	mu    sync.RWMutex
}

func (mc *MutexCounter) Increment() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.value++
}

func (mc *MutexCounter) GetValue() int32 {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	return mc.value
}

func main() {
	numWorker := 5
	times := 10
	var counter ICounter
	var wg sync.WaitGroup

	// Run Atomic-based counter
	counter = &AtomicCounter{}
	for range numWorker {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(counter, times)
		}()
	}
	wg.Wait()
	fmt.Println(counter.GetValue())

	// Run Mutex-based counter
	counter = &MutexCounter{}
	for range numWorker {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(counter, times)
		}()
	}
	wg.Wait()
	fmt.Println(counter.GetValue())

	// ! NOTE: both counters should be 50 (numWorker * times = 50)
}
