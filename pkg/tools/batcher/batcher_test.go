package batcher

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/openimsdk/tools/utils/stringutil"
)

func TestBatcher(t *testing.T) {
	config := Config{
		size:     1000,
		buffer:   10,
		worker:   10,
		interval: 5 * time.Millisecond,
	}

	b := New[string](
		WithSize(config.size),
		WithBuffer(config.buffer),
		WithWorker(config.worker),
		WithInterval(config.interval),
		WithSyncWait(true),
	)

	// Mock Do function to simply print values for demonstration
	b.Do = func(ctx context.Context, channelID int, vals *Msg[string]) {
		t.Logf("Channel %d Processed batch: %v", channelID, vals)
	}
	b.OnComplete = func(lastMessage *string, totalCount int) {
		t.Logf("Completed processing with last message: %v, total count: %d", *lastMessage, totalCount)
	}
	b.Sharding = func(key string) int {
		hashCode := stringutil.GetHashCode(key)
		return int(hashCode) % config.worker
	}
	b.Key = func(data *string) string {
		return *data
	}

	err := b.Start()
	if err != nil {
		t.Fatal(err)
	}

	// Test normal data processing
	for i := 0; i < 10000; i++ {
		data := "data" + fmt.Sprintf("%d", i)
		if err := b.Put(context.Background(), &data); err != nil {
			t.Fatal(err)
		}
	}

	time.Sleep(time.Duration(1) * time.Second)
	start := time.Now()
	// Wait for all processing to finish
	b.Close()

	elapsed := time.Since(start)
	t.Logf("Close took %s", elapsed)

	if len(b.data) != 0 {
		t.Error("Data channel should be empty after closing")
	}
}
