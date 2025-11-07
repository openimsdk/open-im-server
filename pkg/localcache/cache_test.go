package localcache

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	c := New[string](WithExpirationEvict())
	//c := New[string]()
	ctx := context.Background()

	const (
		num  = 10000
		tNum = 10000
		kNum = 100000
		pNum = 100
	)

	getKey := func(v uint64) string {
		return fmt.Sprintf("key_%d", v%kNum)
	}

	start := time.Now()
	t.Log("start", start)

	var (
		get atomic.Int64
		del atomic.Int64
	)

	incrGet := func() {
		if v := get.Add(1); v%pNum == 0 {
			//t.Log("#get count", v/pNum)
		}
	}
	incrDel := func() {
		if v := del.Add(1); v%pNum == 0 {
			//t.Log("@del count", v/pNum)
		}
	}

	var wg sync.WaitGroup

	for i := 0; i < tNum; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			for i := 0; i < num; i++ {
				c.Get(ctx, getKey(rand.Uint64()), func(ctx context.Context) (string, error) {
					return fmt.Sprintf("index_%d", i), nil
				})
				incrGet()
			}
		}()

		go func() {
			defer wg.Done()
			time.Sleep(time.Second / 10)
			for i := 0; i < num; i++ {
				c.Del(ctx, getKey(rand.Uint64()))
				incrDel()
			}
		}()
	}

	wg.Wait()
	end := time.Now()
	t.Log("end", end)
	t.Log("time", end.Sub(start))
	t.Log("get", get.Load())
	t.Log("del", del.Load())
	// 137.35s
}
