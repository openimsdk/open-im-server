package pkg

import "time"

type Config struct {
	TaskPath     string
	ProgressPath string
	Concurrency  int
	Retry        int
	Timeout      time.Duration
	Api          string
	UserID       string
	Secret       string
	TempDir      string
	CacheSize    int64
}
