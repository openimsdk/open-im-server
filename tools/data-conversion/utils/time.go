package utils

import "time"

func InitTime(ts ...*time.Time) {
	for i := range ts {
		if ts[i] == nil {
			continue
		}
		if ts[i].IsZero() || ts[i].UnixMicro() < 0 {
			*ts[i] = time.UnixMicro(0)
		}
	}
}
