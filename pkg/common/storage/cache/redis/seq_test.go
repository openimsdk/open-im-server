package redis

import (
	"context"
	"testing"
	"time"
)

func TestSeq(t *testing.T) {
	ts := NewTestSeq()
	for i := 1; i < 1000000; i++ {
		var size int64 = 100
		first, err := ts.Malloc(context.Background(), "1", size)
		if err != nil {
			t.Logf("[%d] %s", i, err)
			return
		}
		t.Logf("[%d] %d -> %d", i, first+1, first+size)
		time.Sleep(time.Second / 4)
	}
}

func TestDel(t *testing.T) {
	ts := NewTestSeq()
	t.Log(ts.GetMaxSeq(context.Background(), "1"))

}
