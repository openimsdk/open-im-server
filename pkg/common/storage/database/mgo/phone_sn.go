// Copyright © 2024 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.

package mgo

import (
	"context"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewPhoneSNMongo(db *mongo.Database) (database.PhoneSN, error) {
	coll := db.Collection(database.PhoneSNInfoName)
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "phone", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &phoneSNMgo{coll: coll}, nil
}

type phoneSNMgo struct {
	coll *mongo.Collection
}

func (p *phoneSNMgo) GetByPhone(ctx context.Context, phone string) (*model.PhoneSNInfo, error) {
	if phone == "" {
		return nil, nil
	}
	doc, err := mongoutil.FindOne[*model.PhoneSNInfo](ctx, p.coll, bson.M{"phone": phone})
	if err != nil {
		if errs.ErrRecordNotFound.Is(err) {
			return nil, nil
		}
		return nil, err
	}
	return doc, nil
}

func (p *phoneSNMgo) Upsert(ctx context.Context, phone string, userID int64, isSnd bool) error {
	if phone == "" {
		return errs.ErrArgs.WrapMsg("phone is empty")
	}
	now := time.Now().UnixMilli()
	filter := bson.M{"phone": phone}
	setDoc := bson.M{
		"is_snd":      isSnd,
		"user_id":     userID,
		"update_time": now,
	}
	update := bson.M{
		"$set":         setDoc,
		"$setOnInsert": bson.M{"phone": phone},
	}
	opts := options.Update().SetUpsert(true)
	_, err := p.coll.UpdateOne(ctx, filter, update, opts)
	return errs.Wrap(err)
}
