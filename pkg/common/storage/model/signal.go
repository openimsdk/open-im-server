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
	BusyLineUserIDList []string `bson:"busy_line_user_id_list"`
	OfflinePushTitle   string   `bson:"offline_push_title"`
	OfflinePushDesc    string   `bson:"offline_push_desc"`
	OfflinePushEx      string   `bson:"offline_push_ex"`
	CreateTime         int64    `bson:"create_time"`
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
