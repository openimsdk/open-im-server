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

package mgo

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewSignalMongo(db *mongo.Database) (database.SignalDatabase, error) {
	invColl := db.Collection(database.SignalInvitationName)
	_, err := invColl.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "room_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "invitee_user_id_list", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "create_time", Value: -1}},
		},
		// Fix P1(TTL): expire_at 字段为 BSON Date，MongoDB 后台每 60s 扫描一次并自动删除过期文档。
		// 覆盖场景：被叫网络断开、主叫 App 被杀、任何异常中断导致没有 Cancel/Reject/HungUp 的情况。
		{
			Keys:    bson.D{{Key: "expire_at", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(0),
		},
	})
	if err != nil {
		return nil, err
	}

	recColl := db.Collection(database.SignalRecordName)
	_, err = recColl.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "sid", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "send_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "recv_id", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "create_time", Value: -1}},
		},
	})
	if err != nil {
		return nil, err
	}

	return &signalMgo{invColl: invColl, recColl: recColl}, nil
}

type signalMgo struct {
	invColl *mongo.Collection
	recColl *mongo.Collection
}

func (s *signalMgo) CreateInvitation(ctx context.Context, inv *model.SignalInvitation) error {
	return mongoutil.InsertMany(ctx, s.invColl, []*model.SignalInvitation{inv})
}

func (s *signalMgo) GetInvitationByRoomID(ctx context.Context, roomID string) (*model.SignalInvitation, error) {
	return mongoutil.FindOne[*model.SignalInvitation](ctx, s.invColl, bson.M{"room_id": roomID})
}

func (s *signalMgo) GetInvitationByInviteeUserID(ctx context.Context, userID string) (*model.SignalInvitation, error) {
	opts := options.FindOne().SetSort(bson.M{"create_time": -1})
	return mongoutil.FindOne[*model.SignalInvitation](ctx, s.invColl, bson.M{"invitee_user_id_list": userID}, opts)
}

func (s *signalMgo) DeleteInvitation(ctx context.Context, roomID string) error {
	return mongoutil.DeleteMany(ctx, s.invColl, bson.M{"room_id": roomID})
}

func (s *signalMgo) RemoveInvitee(ctx context.Context, roomID string, userID string) error {
	filter := bson.M{"room_id": roomID}
	update := bson.M{"$pull": bson.M{"invitee_user_id_list": userID}}
	if _, err := s.invColl.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	_, err := s.invColl.DeleteOne(ctx, bson.M{
		"room_id":              roomID,
		"invitee_user_id_list": bson.M{"$size": 0},
	})
	return err
}

func (s *signalMgo) SetConnectTime(ctx context.Context, roomID string, connectTimeMs int64) error {
	_, err := s.invColl.UpdateOne(ctx,
		bson.M{"room_id": roomID},
		bson.M{"$set": bson.M{"connect_time": connectTimeMs}},
	)
	return err
}

func (s *signalMgo) GetInvitationByGroupID(ctx context.Context, groupID string) (*model.SignalInvitation, error) {
	opts := options.FindOne().SetSort(bson.M{"create_time": -1})
	return mongoutil.FindOne[*model.SignalInvitation](ctx, s.invColl, bson.M{"group_id": groupID}, opts)
}

func (s *signalMgo) GetInvitationsByRoomIDs(ctx context.Context, roomIDs []string) ([]*model.SignalInvitation, error) {
	return mongoutil.Find[*model.SignalInvitation](ctx, s.invColl, bson.M{"room_id": bson.M{"$in": roomIDs}})
}

func (s *signalMgo) GetBusyUserIDs(ctx context.Context, userIDs []string) ([]string, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	filter := bson.M{
		"$or": bson.A{
			bson.M{"inviter_user_id": bson.M{"$in": userIDs}},
			bson.M{"invitee_user_id_list": bson.M{"$in": userIDs}},
		},
	}
	invitations, err := mongoutil.Find[*model.SignalInvitation](ctx, s.invColl, filter,
		options.Find().SetProjection(bson.M{"inviter_user_id": 1, "invitee_user_id_list": 1}),
	)
	if err != nil {
		return nil, err
	}
	requested := make(map[string]struct{}, len(userIDs))
	for _, uid := range userIDs {
		requested[uid] = struct{}{}
	}
	busySet := make(map[string]struct{})
	for _, inv := range invitations {
		if _, ok := requested[inv.InviterUserID]; ok {
			busySet[inv.InviterUserID] = struct{}{}
		}
		for _, uid := range inv.InviteeUserIDList {
			if _, ok := requested[uid]; ok {
				busySet[uid] = struct{}{}
			}
		}
	}
	busy := make([]string, 0, len(busySet))
	for uid := range busySet {
		busy = append(busy, uid)
	}
	return busy, nil
}

func (s *signalMgo) CreateRecord(ctx context.Context, record *model.SignalRecord) error {
	return mongoutil.InsertMany(ctx, s.recColl, []*model.SignalRecord{record})
}

func (s *signalMgo) SearchRecords(ctx context.Context, sendID, recvID string, sessionType int32, startTime, endTime int64, pagination pagination.Pagination) (int64, []*model.SignalRecord, error) {
	filter := bson.M{}
	if sendID != "" {
		filter["send_id"] = sendID
	}
	if recvID != "" {
		filter["recv_id"] = recvID
	}
	if sessionType != 0 {
		filter["session_type"] = sessionType
	}
	if startTime > 0 || endTime > 0 {
		timeFilter := bson.M{}
		if startTime > 0 {
			timeFilter["$gte"] = startTime
		}
		if endTime > 0 {
			timeFilter["$lte"] = endTime
		}
		filter["create_time"] = timeFilter
	}
	return mongoutil.FindPage[*model.SignalRecord](ctx, s.recColl, filter, pagination, options.Find().SetSort(bson.M{"create_time": -1}))
}

func (s *signalMgo) DeleteRecords(ctx context.Context, sIDs []string) error {
	return mongoutil.DeleteMany(ctx, s.recColl, bson.M{"sid": bson.M{"$in": sIDs}})
}
