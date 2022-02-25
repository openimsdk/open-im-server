package statistics

import "time"

type Statistics struct {
	Count *uint64
	Dr    int
}

func (s *Statistics) output() {
	for {
		time.Sleep(time.Duration(s.Dr) * time.Second)

	}
}

func NewStatistics(count *uint64, dr int) *Statistics {
	p := &Statistics{Count: count}
	go p.output()
}
