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

package table

import (
	"context"
	"time"

	"github.com/OpenIMSDK/tools/pagination"
)

type MeetingInfo struct {
	RoomID      string    `bson:"room_id"`
	MeetingName string    `bson:"meeting_name"`
	HostUserID  string    `bson:"host_user_id"`
	Status      int64     `bson:"status"`
	StartTime   time.Time `bson:"start_time"`
	EndTime     time.Time `bson:"end_time"`
	CreateTime  time.Time `bson:"create_time"`
	Ex          string    `bson:"ex"`
}

type MeetingInterface interface {
	Find(ctx context.Context, roomIDs []string) ([]*MeetingInfo, error)
	CreateMeetingInfo(ctx context.Context, meetingInfo *MeetingInfo) error
	UpdateMeetingInfo(ctx context.Context, roomID string, update map[string]any) error
	GetUnCompleteMeetingIDList(ctx context.Context, roomIDs []string) ([]string, error)
	Delete(ctx context.Context, roomIDs []string) error
	GetMeetingRecords(ctx context.Context, hostUserID string, startTime, endTime time.Time, pagination pagination.Pagination) (int64, []*MeetingInfo, error)
}

type MeetingInvitationInfo struct {
	RoomID     string    `bson:"room_id"`
	UserID     string    `bson:"user_id"`
	CreateTime time.Time `bson:"create_time"`
}

type MeetingInvitationInterface interface {
	FindUserIDs(ctx context.Context, roomID string) ([]string, error)
	CreateMeetingInvitationInfo(ctx context.Context, roomID string, inviteeUserIDs []string) error
	GetUserInvitedMeetingIDs(ctx context.Context, userID string) (meetingIDs []string, err error)
	Delete(ctx context.Context, roomIDs []string) error
	GetMeetingRecords(ctx context.Context, joinedUserID string, startTime, endTime time.Time, pagination pagination.Pagination) (int64, []string, error)
}

type MeetingVideoRecord struct {
	RoomID     string    `bson:"room_id"`
	FileURL    string    `bson:"file_url"`
	CreateTime time.Time `bson:"create_time"`
}

type MeetingRecordInterface interface {
	CreateMeetingVideoRecord(ctx context.Context, meetingVideoRecord *MeetingVideoRecord) error
}
