package memAsyncQueue

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// AsyncQueue is the interface responsible for asynchronous processing of functions.
type AsyncQueue interface {
	Initialize(processFunc func(), workerCount int, bufferSize int)
	Push(task func()) error
}

// MemoryQueue is an implementation of the AsyncQueue interface using a channel to process functions.
type MemoryQueue struct {
	taskChan  chan func()
	wg        sync.WaitGroup
	isStopped bool
	stopMutex sync.Mutex // Mutex to protect access to isStopped
}

func NewMemoryQueue(workerCount int, bufferSize int) *MemoryQueue {
	mq := &MemoryQueue{}                   // Create a new instance of MemoryQueue
	mq.Initialize(workerCount, bufferSize) // Initialize it with specified parameters
	return mq
}

// Initialize sets up the worker nodes and the buffer size of the channel,
// starting internal goroutines to handle tasks from the channel.
func (mq *MemoryQueue) Initialize(workerCount int, bufferSize int) {
	mq.taskChan = make(chan func(), bufferSize) // Initialize the channel with the provided buffer size.
	mq.isStopped = false

	// Start multiple goroutines based on the specified workerCount.
	for i := 0; i < workerCount; i++ {
		mq.wg.Add(1)
		go func(workerID int) {
			defer mq.wg.Done()
			for task := range mq.taskChan {
				fmt.Printf("Worker %d: Executing task\n", workerID)
				task() // Execute the function
			}
		}(i)
	}
}

// Push submits a function to the queue.
// Returns an error if the queue is stopped or if the queue is full.
func (mq *MemoryQueue) Push(task func()) error {
	mq.stopMutex.Lock()
	if mq.isStopped {
		mq.stopMutex.Unlock()
		return errors.New("push failed: queue is stopped")
	}
	mq.stopMutex.Unlock()

	select {
	case mq.taskChan <- task:
		return nil
	case <-time.After(time.Millisecond * 100): // Timeout to prevent deadlock/blocking
		return errors.New("push failed: queue is full")
	}
}

// Stop is used to terminate the internal goroutines and close the channel.
func (mq *MemoryQueue) Stop() {
	mq.stopMutex.Lock()
	mq.isStopped = true
	close(mq.taskChan)
	mq.stopMutex.Unlock()
	mq.wg.Wait()
}
