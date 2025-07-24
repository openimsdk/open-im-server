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

package stores

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model/signal"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/log"
)

type PreKeyStore struct {
	coll *mongo.Collection
}

func NewPreKeyStore(db *mongo.Database) PreKeyStoreInterface {
	coll := db.Collection(signal.SignalPreKeyCollection)
	// Create indexes
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "device_id", Value: 1}, {Key: "key_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "device_id", Value: 1}}},
		{Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "device_id", Value: 1}, {Key: "used", Value: 1}}},
		{Keys: bson.D{{Key: "created_time", Value: 1}}},
		{Keys: bson.D{{Key: "used_time", Value: 1}}},
	})
	if err != nil {
		log.ZWarn(context.Background(), "failed to create indexes for prekey store", err)
	}
	
	return &PreKeyStore{coll: coll}
}

func (s *PreKeyStore) Create(ctx context.Context, prekey *signal.SignalPreKey) error {
	return mongoutil.InsertOne(ctx, s.coll, prekey)
}

func (s *PreKeyStore) CreateBatch(ctx context.Context, prekeys []*signal.SignalPreKey) error {
	if len(prekeys) == 0 {
		return nil
	}
	return mongoutil.InsertMany(ctx, s.coll, prekeys)
}

func (s *PreKeyStore) GetAvailable(ctx context.Context, userID string, deviceID int32) (*signal.SignalPreKey, error) {
	filter := bson.M{
		"user_id":   userID,
		"device_id": deviceID,
		"used":      false,
	}
	
	// Get one available prekey, sorted by creation time (FIFO)
	opts := options.FindOne().SetSort(bson.D{{Key: "created_time", Value: 1}})
	prekey, err := mongoutil.FindOne[*signal.SignalPreKey](ctx, s.coll, filter, opts)
	if err != nil {
		return nil, err
	}
	
	return prekey, nil
}

func (s *PreKeyStore) MarkUsed(ctx context.Context, userID string, deviceID int32, keyID uint32) error {
	filter := bson.M{
		"user_id":   userID,
		"device_id": deviceID,
		"key_id":    keyID,
	}
	
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"used":      true,
			"used_time": now,
		},
	}
	
	return mongoutil.UpdateOne(ctx, s.coll, filter, update, false)
}

func (s *PreKeyStore) Delete(ctx context.Context, userID string, deviceID int32, keyID uint32) error {
	filter := bson.M{
		"user_id":   userID,
		"device_id": deviceID,
		"key_id":    keyID,
	}
	
	return mongoutil.DeleteOne(ctx, s.coll, filter)
}

func (s *PreKeyStore) DeleteAllByUserDevice(ctx context.Context, userID string, deviceID int32) error {
	filter := bson.M{
		"user_id":   userID,
		"device_id": deviceID,
	}
	
	return mongoutil.DeleteMany(ctx, s.coll, filter)
}

func (s *PreKeyStore) CountAvailable(ctx context.Context, userID string, deviceID int32) (int64, error) {
	filter := bson.M{
		"user_id":   userID,
		"device_id": deviceID,
		"used":      false,
	}
	
	return s.coll.CountDocuments(ctx, filter)
}

func (s *PreKeyStore) GetByKeyID(ctx context.Context, userID string, deviceID int32, keyID uint32) (*signal.SignalPreKey, error) {
	filter := bson.M{
		"user_id":   userID,
		"device_id": deviceID,
		"key_id":    keyID,
	}
	
	prekey, err := mongoutil.FindOne[*signal.SignalPreKey](ctx, s.coll, filter)
	if err != nil {
		return nil, err
	}
	
	return prekey, nil
}

func (s *PreKeyStore) CleanupUsed(ctx context.Context, olderThan time.Duration) (int64, error) {
	cutoffTime := time.Now().Add(-olderThan)
	filter := bson.M{
		"used":      true,
		"used_time": bson.M{"$lt": cutoffTime},
	}
	
	result, err := s.coll.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	
	return result.DeletedCount, nil
}