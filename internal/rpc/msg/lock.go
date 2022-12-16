package msg

import (
	"Open_IM/pkg/common/db"
	"time"
)

type MessageLocker interface {
	LockMessageTypeKey(clientMsgID, typeKey string) (err error)
	UnLockMessageTypeKey(clientMsgID string, typeKey string) error
}
type LockerMessage struct{}

func NewLockerMessage() *LockerMessage {
	return &LockerMessage{}
}
func (l *LockerMessage) LockMessageTypeKey(clientMsgID, typeKey string) (err error) {
	for i := 0; i < 3; i++ {
		err = db.DB.LockMessageTypeKey(clientMsgID, typeKey)
		if err != nil {
			time.Sleep(time.Millisecond * 100)
			continue
		} else {
			break
		}
	}
	return err

}
func (l *LockerMessage) UnLockMessageTypeKey(clientMsgID string, typeKey string) error {
	return db.DB.UnLockMessageTypeKey(clientMsgID, typeKey)
}
