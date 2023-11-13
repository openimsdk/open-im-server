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

package msgprocessor

import (
	"testing"

	"github.com/OpenIMSDK/protocol/sdkws"
	"google.golang.org/protobuf/proto"
)

func TestGetNotificationConversationIDByMsg(t *testing.T) {
	type args struct {
		msg *sdkws.MsgData
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetNotificationConversationIDByMsg(tt.args.msg); got != tt.want {
				t.Errorf("GetNotificationConversationIDByMsg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetChatConversationIDByMsg(t *testing.T) {
	type args struct {
		msg *sdkws.MsgData
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetChatConversationIDByMsg(tt.args.msg); got != tt.want {
				t.Errorf("GetChatConversationIDByMsg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenConversationUniqueKey(t *testing.T) {
	type args struct {
		msg *sdkws.MsgData
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenConversationUniqueKey(tt.args.msg); got != tt.want {
				t.Errorf("GenConversationUniqueKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetConversationIDByMsg(t *testing.T) {
	type args struct {
		msg *sdkws.MsgData
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetConversationIDByMsg(tt.args.msg); got != tt.want {
				t.Errorf("GetConversationIDByMsg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetConversationIDBySessionType(t *testing.T) {
	type args struct {
		sessionType int
		ids         []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetConversationIDBySessionType(tt.args.sessionType, tt.args.ids...); got != tt.want {
				t.Errorf("GetConversationIDBySessionType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNotificationConversationIDByConversationID(t *testing.T) {
	type args struct {
		conversationID string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetNotificationConversationIDByConversationID(tt.args.conversationID); got != tt.want {
				t.Errorf("GetNotificationConversationIDByConversationID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetNotificationConversationID(t *testing.T) {
	type args struct {
		sessionType int
		ids         []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetNotificationConversationID(tt.args.sessionType, tt.args.ids...); got != tt.want {
				t.Errorf("GetNotificationConversationID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNotification(t *testing.T) {
	type args struct {
		conversationID string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotification(tt.args.conversationID); got != tt.want {
				t.Errorf("IsNotification() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNotificationByMsg(t *testing.T) {
	type args struct {
		msg *sdkws.MsgData
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNotificationByMsg(tt.args.msg); got != tt.want {
				t.Errorf("IsNotificationByMsg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseConversationID(t *testing.T) {
	type args struct {
		msg *sdkws.MsgData
	}
	tests := []struct {
		name               string
		args               args
		wantIsNotification bool
		wantConversationID string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIsNotification, gotConversationID := ParseConversationID(tt.args.msg)
			if gotIsNotification != tt.wantIsNotification {
				t.Errorf("ParseConversationID() gotIsNotification = %v, want %v", gotIsNotification, tt.wantIsNotification)
			}
			if gotConversationID != tt.wantConversationID {
				t.Errorf("ParseConversationID() gotConversationID = %v, want %v", gotConversationID, tt.wantConversationID)
			}
		})
	}
}

func TestMsgBySeq_Len(t *testing.T) {
	tests := []struct {
		name string
		s    MsgBySeq
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Len(); got != tt.want {
				t.Errorf("MsgBySeq.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgBySeq_Less(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		s    MsgBySeq
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Less(tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("MsgBySeq.Less() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMsgBySeq_Swap(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		s    MsgBySeq
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Swap(tt.args.i, tt.args.j)
		})
	}
}

func TestPb2String(t *testing.T) {
	type args struct {
		pb proto.Message
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Pb2String(tt.args.pb)
			if (err != nil) != tt.wantErr {
				t.Errorf("Pb2String() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Pb2String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString2Pb(t *testing.T) {
	type args struct {
		s  string
		pb proto.Message
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := String2Pb(tt.args.s, tt.args.pb); (err != nil) != tt.wantErr {
				t.Errorf("String2Pb() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
