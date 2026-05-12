// Copyright © 2026 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package mgo

import (
	"context"
	"errors"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewGroupPinnedMsgMongo(db *mongo.Database) (database.GroupPinnedMsg, error) {
	coll := db.Collection(database.GroupPinnedMsgName)
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "group_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &groupPinnedMsgMgo{coll: coll}, nil
}

type groupPinnedMsgMgo struct {
	coll *mongo.Collection
}

func (g *groupPinnedMsgMgo) get(ctx context.Context, groupID string) (*model.GroupPinnedMsg, error) {
	doc, err := mongoutil.FindOne[*model.GroupPinnedMsg](ctx, g.coll, bson.M{"group_id": groupID})
	if err != nil {
		if errs.ErrRecordNotFound.Is(err) || errors.Is(err, mongo.ErrNoDocuments) {
			return &model.GroupPinnedMsg{GroupID: groupID}, nil
		}
		return nil, err
	}
	return doc, nil
}

func (g *groupPinnedMsgMgo) Get(ctx context.Context, groupID string) ([]*model.GroupPinnedMessage, error) {
	doc, err := g.get(ctx, groupID)
	if err != nil {
		return nil, err
	}
	return doc.PinnedMsgs, nil
}

// Pin 置顶一条消息：
// - 若提供的 msg.PinID 为空，则自动生成 ObjectID().Hex()
// - 同 seq 的旧记录会被先移除避免重复
// - 新记录 push 到数组首位，自动滚动保留最近 GroupPinnedMsgMaxKeep 条
func (g *groupPinnedMsgMgo) Pin(ctx context.Context, groupID string, msg *model.GroupPinnedMessage) ([]*model.GroupPinnedMessage, error) {
	if msg == nil {
		return nil, errs.ErrArgs.WrapMsg("pin msg is nil")
	}
	if msg.PinID == "" {
		msg.PinID = primitive.NewObjectID().Hex()
	}
	msg.GroupID = groupID

	if _, err := mongoutil.UpdateOneResult(ctx, g.coll,
		bson.M{"group_id": groupID},
		bson.M{"$pull": bson.M{"pinned_msgs": bson.M{"seq": msg.Seq}}},
	); err != nil {
		return nil, err
	}
	filter := bson.M{"group_id": groupID}
	update := bson.M{
		"$push": bson.M{
			"pinned_msgs": bson.M{
				"$each":     bson.A{msg},
				"$position": 0,
				"$slice":    model.GroupPinnedMsgMaxKeep,
			},
		},
		"$setOnInsert": bson.M{"group_id": groupID},
	}
	opts := options.Update().SetUpsert(true)
	if _, err := g.coll.UpdateOne(ctx, filter, update, opts); err != nil {
		return nil, errs.Wrap(err)
	}
	return g.Get(ctx, groupID)
}

// Unpin 取消置顶：
// - pinID 非空时按 pinID 精确删除（推荐）
// - 否则按 seq 删除
// 返回更新后的置顶列表（可能为空数组）
func (g *groupPinnedMsgMgo) Unpin(ctx context.Context, groupID string, pinID string, seq int64) ([]*model.GroupPinnedMessage, error) {
	if pinID == "" && seq <= 0 {
		return nil, errs.ErrArgs.WrapMsg("either pinID or seq must be provided")
	}
	pull := bson.M{}
	if pinID != "" {
		pull["pin_id"] = pinID
	} else {
		pull["seq"] = seq
	}
	if _, err := mongoutil.UpdateOneResult(ctx, g.coll,
		bson.M{"group_id": groupID},
		bson.M{"$pull": bson.M{"pinned_msgs": pull}},
	); err != nil {
		return nil, err
	}
	return g.Get(ctx, groupID)
}
