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
	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMsgBurnDeadlineMongo(db *mongo.Database) (database.MsgBurnDeadline, error) {
	coll := db.Collection(database.MsgBurnDeadlineName)
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "user_id", Value: 1},
				{Key: "conversation_id", Value: 1},
				{Key: "seq", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "deadline_ms", Value: 1}},
		},
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &msgBurnDeadlineMgo{coll: coll}, nil
}

type msgBurnDeadlineMgo struct {
	coll *mongo.Collection
}

func (m *msgBurnDeadlineMgo) UpsertIfAbsent(ctx context.Context, items []*model.MsgBurnDeadline) error {
	if len(items) == 0 {
		return nil
	}
	models := make([]mongo.WriteModel, 0, len(items))
	for _, item := range items {
		filter := bson.M{
			"user_id":         item.UserID,
			"conversation_id": item.ConversationID,
			"seq":             item.Seq,
		}
		setOnInsert := bson.M{
			"user_id":         item.UserID,
			"conversation_id": item.ConversationID,
			"seq":             item.Seq,
			"peer_id":         item.PeerID,
			"deadline_ms":     item.DeadlineMs,
			"create_time":     item.CreateTime,
		}
		models = append(models,
			mongo.NewUpdateOneModel().
				SetFilter(filter).
				SetUpdate(bson.M{"$setOnInsert": setOnInsert}).
				SetUpsert(true),
		)
	}
	_, err := m.coll.BulkWrite(ctx, models, options.BulkWrite().SetOrdered(false))
	return errs.Wrap(err)
}

func (m *msgBurnDeadlineMgo) FindExpiredGroups(ctx context.Context, nowMs int64, limit int) ([]*database.ExpiredBurnGroup, error) {
	if limit <= 0 {
		return nil, nil
	}
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{"deadline_ms": bson.M{"$lte": nowMs}}}},
		bson.D{{Key: "$group", Value: bson.M{
			"_id": bson.M{
				"user_id":         "$user_id",
				"conversation_id": "$conversation_id",
			},
			"peer_id": bson.M{"$first": "$peer_id"},
			"max_seq": bson.M{"$max": "$seq"},
			"seqs":    bson.M{"$push": "$seq"},
		}}},
		bson.D{{Key: "$limit", Value: int64(limit)}},
	}
	type aggRow struct {
		ID struct {
			UserID         string `bson:"user_id"`
			ConversationID string `bson:"conversation_id"`
		} `bson:"_id"`
		PeerID string  `bson:"peer_id"`
		MaxSeq int64   `bson:"max_seq"`
		Seqs   []int64 `bson:"seqs"`
	}
	rows, err := mongoutil.Aggregate[*aggRow](ctx, m.coll, pipeline)
	if err != nil {
		return nil, err
	}
	res := make([]*database.ExpiredBurnGroup, 0, len(rows))
	for _, r := range rows {
		res = append(res, &database.ExpiredBurnGroup{
			UserID:         r.ID.UserID,
			ConversationID: r.ID.ConversationID,
			PeerID:         r.PeerID,
			MaxSeq:         r.MaxSeq,
			Seqs:           r.Seqs,
		})
	}
	return res, nil
}

func (m *msgBurnDeadlineMgo) DeleteByUserConversationSeqs(ctx context.Context, userID, conversationID string, seqs []int64) error {
	if len(seqs) == 0 {
		return nil
	}
	filter := bson.M{
		"user_id":         userID,
		"conversation_id": conversationID,
		"seq":             bson.M{"$in": seqs},
	}
	_, err := m.coll.DeleteMany(ctx, filter)
	return errs.Wrap(err)
}
