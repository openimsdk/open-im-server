package fcm

import (
	"Open_IM/internal/push"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Push(t *testing.T) {
	offlinePusher := NewFcm()
	resp, err := offlinePusher.Push([]string{"test_uid"}, "test", "test", "12321", push.PushOpts{})
	assert.Nil(t, err)
	fmt.Println(resp)
}
