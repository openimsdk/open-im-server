package localcache

import (
	"sync"
	"time"
)

type Execute[K comparable, V any] func(K, V)

type Task[K comparable, V any] struct {
	key   K
	value V
}

type TimeWheel[K comparable, V any] struct {
	ticker     *time.Ticker
	slots      [][]Task[K, V]
	currentPos int
	size       int
	slotMutex  sync.Mutex
	execute    Execute[K, V]
}

func NewTimeWheel[K comparable, V any](size int, tickDuration time.Duration, execute Execute[K, V]) *TimeWheel[K, V] {
	return &TimeWheel[K, V]{
		ticker:     time.NewTicker(tickDuration),
		slots:      make([][]Task[K, V], size),
		currentPos: 0,
		size:       size,
		execute:    execute,
	}
}

func (tw *TimeWheel[K, V]) Start() {
	for range tw.ticker.C {
		tw.tick()
	}
}

func (tw *TimeWheel[K, V]) Stop() {
	tw.ticker.Stop()
}

func (tw *TimeWheel[K, V]) tick() {
	tw.slotMutex.Lock()
	defer tw.slotMutex.Unlock()

	tasks := tw.slots[tw.currentPos]
	tw.slots[tw.currentPos] = nil
	if len(tasks) > 0 {
		go func(tasks []Task[K, V]) {
			for _, task := range tasks {
				tw.execute(task.key, task.value)
			}
		}(tasks)
	}

	tw.currentPos = (tw.currentPos + 1) % tw.size
}

func (tw *TimeWheel[K, V]) AddTask(delay int, task Task[K, V]) {
	if delay < 0 || delay >= tw.size {
		return
	}

	tw.slotMutex.Lock()
	defer tw.slotMutex.Unlock()

	pos := (tw.currentPos + delay) % tw.size
	tw.slots[pos] = append(tw.slots[pos], task)
}
