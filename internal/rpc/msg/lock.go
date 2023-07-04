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
