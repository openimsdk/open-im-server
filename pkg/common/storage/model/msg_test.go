package model

import "testing"

func TestMsgDocModelGetLimitForSingleDoc(t *testing.T) {
	model := new(MsgDocModel)
	for _, test := range []struct{ seq, want int64 }{{1, 0}, {99, 98}, {100, 99}, {101, 0}} {
		if got := model.GetLimitForSingleDoc(test.seq); got != test.want {
			t.Errorf("GetLimitForSingleDoc(%d) = %d, want %d", test.seq, got, test.want)
		}
	}
}
