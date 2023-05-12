package msg

import (
	"context"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
)

const GlOBALLOCK = "GLOBAL_LOCK"

type MessageLocker interface {
	LockMessageTypeKey(ctx context.Context, clientMsgID, typeKey string) (err error)
	UnLockMessageTypeKey(ctx context.Context, clientMsgID string, typeKey string) error
	LockGlobalMessage(ctx context.Context, clientMsgID string) (err error)
	UnLockGlobalMessage(ctx context.Context, clientMsgID string) (err error)
}
type LockerMessage struct {
	cache cache.MsgModel
}

func NewLockerMessage(cache cache.MsgModel) *LockerMessage {
	return &LockerMessage{cache: cache}
}
func (l *LockerMessage) LockMessageTypeKey(ctx context.Context, clientMsgID, typeKey string) (err error) {
	for i := 0; i < 3; i++ {
		err = l.cache.LockMessageTypeKey(ctx, clientMsgID, typeKey)
		if err != nil {
			time.Sleep(time.Millisecond * 100)
			continue
		} else {
			break
		}
	}
	return err

}
func (l *LockerMessage) LockGlobalMessage(ctx context.Context, clientMsgID string) (err error) {
	for i := 0; i < 3; i++ {
		err = l.cache.LockMessageTypeKey(ctx, clientMsgID, GlOBALLOCK)
		if err != nil {
			time.Sleep(time.Millisecond * 100)
			continue
		} else {
			break
		}
	}
	return err

}
func (l *LockerMessage) UnLockMessageTypeKey(ctx context.Context, clientMsgID string, typeKey string) error {
	return l.cache.UnLockMessageTypeKey(ctx, clientMsgID, typeKey)
}
func (l *LockerMessage) UnLockGlobalMessage(ctx context.Context, clientMsgID string) error {
	return l.cache.UnLockMessageTypeKey(ctx, clientMsgID, GlOBALLOCK)
}
