package mgo

import (
	"testing"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/protocol/constant"
)

func TestGroupRequestHandlerUpdate(t *testing.T) {
	handledTime := time.Date(2026, time.September, 17, 12, 0, 0, 0, time.UTC)
	update := groupRequestHandlerUpdate(&model.GroupRequest{
		HandleResult: constant.GroupResponseAgree,
		HandledMsg:   "approved",
		HandleUserID: "admin",
		HandledTime:  handledTime,
	})

	if _, ok := update["handle_msg"]; ok {
		t.Fatal("handler update contains obsolete handle_msg field")
	}
	want := map[string]any{
		"handled_msg": "approved", "handle_result": int32(constant.GroupResponseAgree),
		"handle_user_id": "admin", "handled_time": handledTime,
	}
	for key, expected := range want {
		if got := update[key]; got != expected {
			t.Fatalf("handler update %s = %v, want %v", key, got, expected)
		}
	}
}
