package statistics

import (
	"time"
)

type Statistics struct {
	Count      *uint64
	ModuleName string
	PrintArgs  string
	SleepTime  int
}

func (s *Statistics) output() {
	t := time.NewTicker(time.Duration(s.SleepTime) * time.Second)
	defer t.Stop()
	//var sum uint64
	for {
		//sum = *s.Count
		select {
		case <-t.C:
		}
		//log.NewWarn("", " system stat ", s.ModuleName, s.PrintArgs, *s.Count-sum, "total:", *s.Count)
	}
}

func NewStatistics(count *uint64, moduleName, printArgs string, sleepTime int) *Statistics {
	p := &Statistics{Count: count, ModuleName: moduleName, SleepTime: sleepTime, PrintArgs: printArgs}
	go p.output()
	return p
}
