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
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewGroupMsgBurnRecordMongo 初始化 group_msg_burn_record 集合及索引。
func NewGroupMsgBurnRecordMongo(db *mongo.Database) (database.GroupMsgBurnRecord, error) {
	coll := db.Collection(database.GroupMsgBurnRecordName)
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "group_id", Value: 1},
				{Key: "seq", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{{Key: "burn_end_time", Value: 1}},
		},
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &groupMsgBurnRecordMgo{coll: coll}, nil
}

type groupMsgBurnRecordMgo struct {
	coll *mongo.Collection
}

// UpsertOnRead 对每条 seq 执行 upsert：
//   - 首次插入（$setOnInsert）写入 member_count、burn_end_time、create_time，read_count 初始化为 1。
//   - 已存在时仅对 read_count 执行 $inc 1。
func (m *groupMsgBurnRecordMgo) UpsertOnRead(ctx context.Context, groupID string, seqs []int64, memberCount int32, burnEndTimeMs int64) error {
	if len(seqs) == 0 {
		return nil
	}
	now := time.Now().UnixMilli()
	models := make([]mongo.WriteModel, 0, len(seqs))
	for _, seq := range seqs {
		filter := bson.M{
			"group_id": groupID,
			"seq":      seq,
		}
		update := bson.M{
			"$inc": bson.M{"read_count": int32(1)},
			"$setOnInsert": bson.M{
				"group_id":      groupID,
				"seq":           seq,
				"member_count":  memberCount,
				"burn_end_time": burnEndTimeMs,
				"create_time":   now,
			},
		}
		models = append(models,
			mongo.NewUpdateOneModel().
				SetFilter(filter).
				SetUpdate(update).
				SetUpsert(true),
		)
	}
	_, err := m.coll.BulkWrite(ctx, models, options.BulkWrite().SetOrdered(false))
	return errs.Wrap(err)
}

// FindExpired 查询 burn_end_time <= nowMs 且 read_count >= member_count 的记录，
// 按 group_id 聚合后返回每组的最大 seq 与所有 seq 列表。
func (m *groupMsgBurnRecordMgo) FindExpired(ctx context.Context, nowMs int64, limit int) ([]*database.ExpiredGroupBurn, error) {
	if limit <= 0 {
		return nil, nil
	}
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$match", Value: bson.M{
			"burn_end_time": bson.M{"$lte": nowMs},
			"$expr":         bson.M{"$gte": bson.A{"$read_count", "$member_count"}},
		}}},
		bson.D{{Key: "$group", Value: bson.M{
			"_id":     "$group_id",
			"max_seq": bson.M{"$max": "$seq"},
			"seqs":    bson.M{"$push": "$seq"},
		}}},
		bson.D{{Key: "$limit", Value: int64(limit)}},
	}

	type aggRow struct {
		GroupID string  `bson:"_id"`
		MaxSeq  int64   `bson:"max_seq"`
		Seqs    []int64 `bson:"seqs"`
	}
	rows, err := mongoutil.Aggregate[*aggRow](ctx, m.coll, pipeline)
	if err != nil {
		return nil, err
	}
	res := make([]*database.ExpiredGroupBurn, 0, len(rows))
	for _, r := range rows {
		res = append(res, &database.ExpiredGroupBurn{
			GroupID: r.GroupID,
			MaxSeq:  r.MaxSeq,
			Seqs:    r.Seqs,
		})
	}
	return res, nil
}

// DeleteByGroupSeqs 删除指定群下一批 seq 的记录。
func (m *groupMsgBurnRecordMgo) DeleteByGroupSeqs(ctx context.Context, groupID string, seqs []int64) error {
	if len(seqs) == 0 {
		return nil
	}
	filter := bson.M{
		"group_id": groupID,
		"seq":      bson.M{"$in": seqs},
	}
	_, err := m.coll.DeleteMany(ctx, filter)
	return errs.Wrap(err)
}
