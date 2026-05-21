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

package controller

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/db/pagination"
)

// RtcDatabase defines the business-level operations on RTC signal storage.
type RtcDatabase interface {
	CreateInvitation(ctx context.Context, inv *model.SignalInvitation) error
	GetInvitationByRoomID(ctx context.Context, roomID string) (*model.SignalInvitation, error)
	GetInvitationByInviteeUserID(ctx context.Context, userID string) (*model.SignalInvitation, error)
	DeleteInvitation(ctx context.Context, roomID string) error
	RemoveInvitee(ctx context.Context, roomID string, userID string) error
	SetConnectTime(ctx context.Context, roomID string, connectTimeMs int64) error
	GetInvitationByGroupID(ctx context.Context, groupID string) (*model.SignalInvitation, error)
	GetInvitationsByRoomIDs(ctx context.Context, roomIDs []string) ([]*model.SignalInvitation, error)
	// GetBusyUserIDs returns the subset of userIDs that are currently in an active call.
	GetBusyUserIDs(ctx context.Context, userIDs []string) ([]string, error)

	CreateRecord(ctx context.Context, record *model.SignalRecord) error
	SearchRecords(ctx context.Context, sendID, recvID string, sessionType int32, startTime, endTime int64, pagination pagination.Pagination) (int64, []*model.SignalRecord, error)
	DeleteRecords(ctx context.Context, sIDs []string) error

	// Call record operations (通话记录).
	CreateCallRecord(ctx context.Context, record *model.CallRecord) error
	SearchCallRecords(ctx context.Context, userID string, status int32, startTime, endTime int64, keyword string, pg pagination.Pagination) (int64, []*model.CallRecord, error)
	DeleteCallRecords(ctx context.Context, sids []string) error
}

type rtcDatabase struct {
	db         database.SignalDatabase
	callRecord database.CallRecordDatabase
}

func NewRtcDatabase(db database.SignalDatabase, callRecord database.CallRecordDatabase) RtcDatabase {
	return &rtcDatabase{db: db, callRecord: callRecord}
}

func (r *rtcDatabase) CreateInvitation(ctx context.Context, inv *model.SignalInvitation) error {
	return r.db.CreateInvitation(ctx, inv)
}

func (r *rtcDatabase) GetInvitationByRoomID(ctx context.Context, roomID string) (*model.SignalInvitation, error) {
	return r.db.GetInvitationByRoomID(ctx, roomID)
}

func (r *rtcDatabase) GetInvitationByInviteeUserID(ctx context.Context, userID string) (*model.SignalInvitation, error) {
	return r.db.GetInvitationByInviteeUserID(ctx, userID)
}

func (r *rtcDatabase) DeleteInvitation(ctx context.Context, roomID string) error {
	return r.db.DeleteInvitation(ctx, roomID)
}

func (r *rtcDatabase) RemoveInvitee(ctx context.Context, roomID string, userID string) error {
	return r.db.RemoveInvitee(ctx, roomID, userID)
}

func (r *rtcDatabase) SetConnectTime(ctx context.Context, roomID string, connectTimeMs int64) error {
	return r.db.SetConnectTime(ctx, roomID, connectTimeMs)
}

func (r *rtcDatabase) GetInvitationByGroupID(ctx context.Context, groupID string) (*model.SignalInvitation, error) {
	return r.db.GetInvitationByGroupID(ctx, groupID)
}

func (r *rtcDatabase) GetInvitationsByRoomIDs(ctx context.Context, roomIDs []string) ([]*model.SignalInvitation, error) {
	return r.db.GetInvitationsByRoomIDs(ctx, roomIDs)
}

func (r *rtcDatabase) GetBusyUserIDs(ctx context.Context, userIDs []string) ([]string, error) {
	return r.db.GetBusyUserIDs(ctx, userIDs)
}

func (r *rtcDatabase) CreateRecord(ctx context.Context, record *model.SignalRecord) error {
	return r.db.CreateRecord(ctx, record)
}

func (r *rtcDatabase) SearchRecords(ctx context.Context, sendID, recvID string, sessionType int32, startTime, endTime int64, pg pagination.Pagination) (int64, []*model.SignalRecord, error) {
	return r.db.SearchRecords(ctx, sendID, recvID, sessionType, startTime, endTime, pg)
}

func (r *rtcDatabase) DeleteRecords(ctx context.Context, sIDs []string) error {
	return r.db.DeleteRecords(ctx, sIDs)
}

func (r *rtcDatabase) CreateCallRecord(ctx context.Context, record *model.CallRecord) error {
	return r.callRecord.CreateCallRecord(ctx, record)
}

func (r *rtcDatabase) SearchCallRecords(ctx context.Context, userID string, status int32, startTime, endTime int64, keyword string, pg pagination.Pagination) (int64, []*model.CallRecord, error) {
	return r.callRecord.SearchCallRecords(ctx, userID, status, startTime, endTime, keyword, pg)
}

func (r *rtcDatabase) DeleteCallRecords(ctx context.Context, sids []string) error {
	return r.callRecord.DeleteCallRecords(ctx, sids)
}
