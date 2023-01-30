package utils

import "testing"

func TestGenID(t *testing.T) {
	m := map[string]struct{}{}
	for i := 0; i < 2000; i++ {
		got := GenID()
		if _, ok := m[got]; !ok {
			m[got] = struct{}{}
		} else {
			t.Error("id generate error", got)
		}
	}
}
