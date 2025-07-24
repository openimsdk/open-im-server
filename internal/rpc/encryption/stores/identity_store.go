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
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model/signal"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
)

type IdentityStore struct {
	coll *mongo.Collection
}

func NewIdentityStore(db *mongo.Database) IdentityStoreInterface {
	coll := db.Collection(signal.SignalIdentityKeyCollection)
	// Create indexes
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "device_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{Keys: bson.D{{Key: "user_id", Value: 1}}},
	})
	if err != nil {
		log.ZWarn(context.Background(), "failed to create indexes for identity store", err)
	}
	
	return &IdentityStore{coll: coll}
}

func (s *IdentityStore) Create(ctx context.Context, identityKey *signal.SignalIdentityKey) error {
	return mongoutil.InsertOne(ctx, s.coll, identityKey)
}

func (s *IdentityStore) Update(ctx context.Context, userID string, deviceID int32, identityKey *signal.SignalIdentityKey) error {
	filter := bson.M{"user_id": userID, "device_id": deviceID}
	update := bson.M{"$set": identityKey}
	return mongoutil.UpdateOne(ctx, s.coll, filter, update, false)
}

func (s *IdentityStore) Get(ctx context.Context, userID string, deviceID int32) (*signal.SignalIdentityKey, error) {
	filter := bson.M{"user_id": userID, "device_id": deviceID}
	identityKey, err := mongoutil.FindOne[*signal.SignalIdentityKey](ctx, s.coll, filter)
	if err != nil {
		if errs.ErrRecordNotFound.Is(err) {
			return nil, fmt.Errorf("identity key not found for user %s device %d", userID, deviceID)
		}
		return nil, err
	}
	return identityKey, nil
}

func (s *IdentityStore) Delete(ctx context.Context, userID string, deviceID int32) error {
	filter := bson.M{"user_id": userID, "device_id": deviceID}
	return mongoutil.DeleteOne(ctx, s.coll, filter)
}

func (s *IdentityStore) GetByUserID(ctx context.Context, userID string) ([]*signal.SignalIdentityKey, error) {
	filter := bson.M{"user_id": userID}
	return mongoutil.Find[*signal.SignalIdentityKey](ctx, s.coll, filter)
}

func (s *IdentityStore) Exists(ctx context.Context, userID string, deviceID int32) (bool, error) {
	filter := bson.M{"user_id": userID, "device_id": deviceID}
	count, err := s.coll.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}