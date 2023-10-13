package utils

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func FindAndInsert[V2 any, V3 schema.Tabler](v2db *gorm.DB, v3db *gorm.DB, fn func(V2) (V3, bool)) (string, error) {
	const batchSize = 100
	var t V3
	name := t.TableName()
	if err := v3db.AutoMigrate(&t); err != nil {
		return name, fmt.Errorf("auto migrate v3 %s failed %w", name, err)
	}
	for i := 0; ; i++ {
		var v2s []V2
		if err := v2db.Offset(i * batchSize).Limit(batchSize).Find(&v2s).Error; err != nil {
			return name, fmt.Errorf("find v2 %s failed %w", name, err)
		}
		if len(v2s) == 0 {
			return name, nil
		}
		v3s := make([]V3, 0, len(v2s))
		for _, v := range v2s {
			res, ok := fn(v)
			if ok {
				v3s = append(v3s, res)
			}
		}
		if len(v3s) == 0 {
			continue
		}
		if err := v3db.Create(&v3s).Error; err != nil {
			return name, fmt.Errorf("insert v3 %s failed %w", name, err)
		}
	}
}

type TakeList []Task

func (l *TakeList) Append(fn ...Task) {
	*l = append(*l, fn...)
}

type Task func() (string, error)

func RunTask(concurrency int, tasks TakeList) []string {
	if len(tasks) == 0 {
		return []string{}
	}
	if concurrency < 1 {
		concurrency = 1
	}
	if concurrency > len(tasks) {
		concurrency = len(tasks)
	}

	taskCh := make(chan func() (string, error), 4)
	go func() {
		defer close(taskCh)
		for i := range tasks {
			taskCh <- tasks[i]
		}
	}()

	var lock sync.Mutex
	var failedTables []string

	var wg sync.WaitGroup
	wg.Add(concurrency)
	var count int64

	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for task := range taskCh {
				name, err := task()
				index := atomic.AddInt64(&count, 1)
				if err == nil {
					log.Printf("[%d/%d] %s success\n", index, len(tasks), name)
				} else {
					lock.Lock()
					failedTables = append(failedTables, name)
					lock.Unlock()
					log.Printf("[%d/%d] %s failed %s\n", index, len(tasks), name, err)
					return
				}
			}
		}()
	}

	wg.Wait()
	if len(failedTables) == 0 {
		log.Println("all tables success")
	} else {
		log.Printf("failed tables %d: %+v\n", len(failedTables), failedTables)
	}

	return failedTables
}
