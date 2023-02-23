package fcm

import (
	"Open_IM/internal/push"
	"Open_IM/pkg/common/db/cache"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Push(t *testing.T) {
	var redis cache.MsgCache
	offlinePusher := NewClient(redis)
	err := offlinePusher.Push(context.Background(), []string{"userID1"}, "test", "test", &push.Opts{})
	assert.Nil(t, err)
}
