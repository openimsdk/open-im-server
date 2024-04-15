package memAsyncQueue

import (
	"testing"
	"time"
)

// TestPushSuccess tests the successful pushing of data into the queue.
func TestPushSuccess(t *testing.T) {
	queue := &MemoryQueue{}
	queue.Initialize(func(data any) {}, 1, 5) // Small buffer size for test

	// Try to push data that should succeed
	err := queue.Push("test data")
	if err != nil {
		t.Errorf("Push should succeed, but got error: %v", err)
	}
}

// TestPushFailWhenFull tests that pushing to a full queue results in an error.
func TestPushFailWhenFull(t *testing.T) {
	queue := &MemoryQueue{}
	queue.Initialize(func(data any) {
		time.Sleep(100 * time.Millisecond) // Simulate work to delay processing
	}, 1, 1) // Very small buffer to fill quickly

	queue.Push("data 1")        // Fill the buffer
	err := queue.Push("data 2") // This should fail

	if err == nil {
		t.Error("Expected an error when pushing to full queue, but got none")
	}
}

// TestPushFailWhenStopped tests that pushing to a stopped queue results in an error.
func TestPushFailWhenStopped(t *testing.T) {
	queue := &MemoryQueue{}
	queue.Initialize(func(data any) {}, 1, 1)

	queue.Stop() // Stop the queue before pushing
	err := queue.Push("test data")

	if err == nil {
		t.Error("Expected an error when pushing to stopped queue, but got none")
	}
}

// TestQueueOperationSequence tests a sequence of operations to ensure the queue handles them correctly.
func TestQueueOperationSequence(t *testing.T) {
	queue := &MemoryQueue{}
	queue.Initialize(func(data any) {}, 1, 2)

	// Sequence of pushes and a stop
	err := queue.Push("data 1")
	if err != nil {
		t.Errorf("Failed to push data 1: %v", err)
	}

	err = queue.Push("data 2")
	if err != nil {
		t.Errorf("Failed to push data 2: %v", err)
	}

	queue.Stop()               // Stop the queue
	err = queue.Push("data 3") // This push should fail
	if err == nil {
		t.Error("Expected an error when pushing after stop, but got none")
	}
}

// TestBlockingOnFull tests that the queue does not block indefinitely when full.
func TestBlockingOnFull(t *testing.T) {
	queue := &MemoryQueue{}
	queue.Initialize(func(data any) {
		time.Sleep(1 * time.Second) // Simulate a long processing time
	}, 1, 1)

	queue.Push("data 1") // Fill the queue

	start := time.Now()
	err := queue.Push("data 2") // This should time out
	duration := time.Since(start)

	if err == nil {
		t.Error("Expected an error due to full queue, but got none")
	}

	if duration >= time.Second {
		t.Errorf("Push blocked for too long, duration: %v", duration)
	}
}
