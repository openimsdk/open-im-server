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

package locker

import (
	"context"
	"reflect"
	"testing"

	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
)

func TestNewLockerMessage(t *testing.T) {
	type args struct {
		cache cache.MsgModel
	}
	tests := []struct {
		name string
		args args
		want *LockerMessage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLockerMessage(tt.args.cache); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLockerMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLockerMessage_LockMessageTypeKey(t *testing.T) {
	type fields struct {
		cache cache.MsgModel
	}
	type args struct {
		ctx         context.Context
		clientMsgID string
		typeKey     string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LockerMessage{
				cache: tt.fields.cache,
			}
			if err := l.LockMessageTypeKey(tt.args.ctx, tt.args.clientMsgID, tt.args.typeKey); (err != nil) != tt.wantErr {
				t.Errorf("LockerMessage.LockMessageTypeKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLockerMessage_LockGlobalMessage(t *testing.T) {
	type fields struct {
		cache cache.MsgModel
	}
	type args struct {
		ctx         context.Context
		clientMsgID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LockerMessage{
				cache: tt.fields.cache,
			}
			if err := l.LockGlobalMessage(tt.args.ctx, tt.args.clientMsgID); (err != nil) != tt.wantErr {
				t.Errorf("LockerMessage.LockGlobalMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLockerMessage_UnLockMessageTypeKey(t *testing.T) {
	type fields struct {
		cache cache.MsgModel
	}
	type args struct {
		ctx         context.Context
		clientMsgID string
		typeKey     string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LockerMessage{
				cache: tt.fields.cache,
			}
			if err := l.UnLockMessageTypeKey(tt.args.ctx, tt.args.clientMsgID, tt.args.typeKey); (err != nil) != tt.wantErr {
				t.Errorf("LockerMessage.UnLockMessageTypeKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLockerMessage_UnLockGlobalMessage(t *testing.T) {
	type fields struct {
		cache cache.MsgModel
	}
	type args struct {
		ctx         context.Context
		clientMsgID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &LockerMessage{
				cache: tt.fields.cache,
			}
			if err := l.UnLockGlobalMessage(tt.args.ctx, tt.args.clientMsgID); (err != nil) != tt.wantErr {
				t.Errorf("LockerMessage.UnLockGlobalMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
