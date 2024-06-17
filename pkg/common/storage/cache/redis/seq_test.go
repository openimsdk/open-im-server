package redis

import (
	"context"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSeq(t *testing.T) {
	ts := NewTestSeq()
	var (
		wg    sync.WaitGroup
		speed atomic.Int64
	)

	const count = 256
	wg.Add(count)
	for i := 0; i < count; i++ {
		index := i + 1
		go func() {
			defer wg.Done()
			var size int64 = 1
			cID := strconv.Itoa(index * 100)
			for i := 1; ; i++ {
				first, err := ts.mgo.Malloc(context.Background(), cID, size) // mongo
				//first, err := ts.Malloc(context.Background(), cID, size) // redis
				if err != nil {
					t.Logf("[%d-%d] %s %s", index, i, cID, err)
					return
				}
				speed.Add(size)
				_ = first
				//t.Logf("[%d] %d -> %d", i, first+1, first+size)
			}
		}()
	}

	done := make(chan struct{})

	go func() {
		wg.Wait()
		close(done)
	}()

	ticker := time.NewTicker(time.Second)

	for {
		select {
		case <-done:
			ticker.Stop()
			return
		case <-ticker.C:
			value := speed.Swap(0)
			t.Logf("speed: %d/s", value)
		}
	}

	//for i := 1; i < 1000000; i++ {
	//	var size int64 = 100
	//	first, err := ts.Malloc(context.Background(), "1", size)
	//	if err != nil {
	//		t.Logf("[%d] %s", i, err)
	//		return
	//	}
	//	t.Logf("[%d] %d -> %d", i, first+1, first+size)
	//	time.Sleep(time.Second / 4)
	//}
}

func TestDel(t *testing.T) {
	ts := NewTestSeq()
	t.Log(ts.GetMaxSeq(context.Background(), "1"))

}
