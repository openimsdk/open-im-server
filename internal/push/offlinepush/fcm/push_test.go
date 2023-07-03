package fcm

import (
	"context"
	"testing"

	"github.com/OpenIMSDK/Open-IM-Server/internal/push/offlinepush"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	"github.com/stretchr/testify/assert"
)

func Test_Push(t *testing.T) {
	var redis cache.MsgModel
	offlinePusher := NewClient(redis)
	err := offlinePusher.Push(context.Background(), []string{"userID1"}, "test", "test", &offlinepush.Opts{})
	assert.Nil(t, err)
}
