// Copyright © 2024 OpenIM. All rights reserved.
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

package model

import "time"

// Call record status values for CallRecord.Status.
const (
	CallStatusAnswered     int32 = 1 // 已接听
	CallStatusNotConnected int32 = 2 // 未接通
)

// Call record direction values for the querying user's perspective.
const (
	CallDirectionOutgoing int32 = 1 // 主叫（发起方）
	CallDirectionIncoming int32 = 2 // 被叫（接收方）
)

// SignalInvitation stores an ongoing or pending signal invitation, keyed by roomID.
// It is created when a call is initiated and can be queried when the callee starts the app.
type SignalInvitation struct {
	RoomID             string   `bson:"room_id"`
	InviterUserID      string   `bson:"inviter_user_id"`
	InviteeUserIDList  []string `bson:"invitee_user_id_list"`
	CustomData         string   `bson:"custom_data"`
	GroupID            string   `bson:"group_id"`
	Timeout            int32    `bson:"timeout"`
	MediaType          string   `bson:"media_type"`
	PlatformID         int32    `bson:"platform_id"`
	SessionType        int32    `bson:"session_type"`
	InitiateTime       int64    `bson:"initiate_time"`
	// ConnectTime is the Unix ms timestamp when a callee accepted the call (0 until answered).
	ConnectTime        int64    `bson:"connect_time"`
	BusyLineUserIDList []string `bson:"busy_line_user_id_list"`
	OfflinePushTitle   string   `bson:"offline_push_title"`
	OfflinePushDesc    string   `bson:"offline_push_desc"`
	OfflinePushEx      string   `bson:"offline_push_ex"`
	CreateTime         int64    `bson:"create_time"`
	// ExpireAt 是 MongoDB BSON Date 类型，供 TTL 索引自动清理过期邀请（无人响应/异常中断场景）。
	// 值 = 创建时间 + Timeout + 30s 缓冲，由 invitationToModel 负责填充。
	ExpireAt time.Time `bson:"expire_at"`
}

// CallRecord stores a completed call event (answered or not connected) for call history.
type CallRecord struct {
	SID                 string   `bson:"sid"`
	RoomID              string   `bson:"room_id"`
	Status              int32    `bson:"status"`                // CallStatusAnswered / CallStatusNotConnected
	Duration            int64    `bson:"duration"`              // total duration in seconds (initiate→end); kept for backward compat
	DialDuration        int64    `bson:"dial_duration"`         // 拨打时长: initiate→connect (answered) or initiate→end (not connected), seconds
	CallDuration        int64    `bson:"call_duration"`         // 通话时长: connect→end for answered calls; 0 if not connected, seconds
	CreateTime          int64    `bson:"create_time"`           // Unix ms, when the call was initiated
	MediaType           string   `bson:"media_type"`            // "audio" or "video"
	SessionType         int32    `bson:"session_type"`
	InviterUserID       string   `bson:"inviter_user_id"`
	InviterUserNickname string   `bson:"inviter_user_nickname"`
	InviterUserFaceURL  string   `bson:"inviter_user_face_url"`
	InviteeUserIDList   []string `bson:"invitee_user_id_list"`  // all invitees (for search by participant)
	GroupID             string   `bson:"group_id"`
	GroupName           string   `bson:"group_name"`
}

// SignalRecord stores a completed call record used for history queries.
type SignalRecord struct {
	SID                  string   `bson:"sid"`
	RoomID               string   `bson:"room_id"`
	FileName             string   `bson:"file_name"`
	MediaType            string   `bson:"media_type"`
	SessionType          int32    `bson:"session_type"`
	InviterUserID        string   `bson:"inviter_user_id"`
	InviterUserNickname  string   `bson:"inviter_user_nickname"`
	GroupID              string   `bson:"group_id"`
	GroupName            string   `bson:"group_name"`
	InviterUserIDList    []string `bson:"inviter_user_id_list"`
	SendID               string   `bson:"send_id"`
	RecvID               string   `bson:"recv_id"`
	CreateTime           int64    `bson:"create_time"`
	EndTime              int64    `bson:"end_time"`
	FileSize             string   `bson:"file_size"`
	FileURL              string   `bson:"file_url"`
}
