// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package msg

import (
	"Open_IM/pkg/common/db"
	"time"
)

const GlOBLLOCK = "GLOBAL_LOCK"

type MessageLocker interface {
	LockMessageTypeKey(clientMsgID, typeKey string) (err error)
	UnLockMessageTypeKey(clientMsgID string, typeKey string) error
	LockGlobalMessage(clientMsgID string) (err error)
	UnLockGlobalMessage(clientMsgID string) (err error)
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
func (l *LockerMessage) LockGlobalMessage(clientMsgID string) (err error) {
	for i := 0; i < 3; i++ {
		err = db.DB.LockMessageTypeKey(clientMsgID, GlOBLLOCK)
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
func (l *LockerMessage) UnLockGlobalMessage(clientMsgID string) error {
	return db.DB.UnLockMessageTypeKey(clientMsgID, GlOBLLOCK)
}
