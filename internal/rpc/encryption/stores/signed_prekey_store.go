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

type SignedPreKeyStore struct {
	coll *mongo.Collection
}

func NewSignedPreKeyStore(db *mongo.Database) SignedPreKeyStoreInterface {
	coll := db.Collection(signal.SignalSignedPreKeyCollection)
	// Create indexes
	_, err := coll.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "user_id", Value: 1}, {Key: "device_id", Value: 1}, {Key: "key_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "device_id", Value: 1}}},
		{Keys: bson.D{{Key: "user_id", Value: 1}, {Key: "device_id", Value: 1}, {Key: "active", Value: 1}}},
		{Keys: bson.D{{Key: "active", Value: 1}, {Key: "created_time", Value: 1}}},
		{Keys: bson.D{{Key: "created_time", Value: 1}}},
	})
	if err != nil {
		log.ZWarn(context.Background(), "failed to create indexes for signed prekey store", err)
	}
	
	return &SignedPreKeyStore{coll: coll}
}

func (s *SignedPreKeyStore) Create(ctx context.Context, signedPrekey *signal.SignalSignedPreKey) error {
	return mongoutil.InsertOne(ctx, s.coll, signedPrekey)
}

func (s *SignedPreKeyStore) Update(ctx context.Context, userID string, deviceID int32, keyID uint32, signedPrekey *signal.SignalSignedPreKey) error {
	filter := bson.M{
		"user_id":   userID,
		"device_id": deviceID,
		"key_id":    keyID,
	}
	
	update := bson.M{"$set": signedPrekey}
	return mongoutil.UpdateOne(ctx, s.coll, filter, update, false)
}

func (s *SignedPreKeyStore) GetActive(ctx context.Context, userID string, deviceID int32) (*signal.SignalSignedPreKey, error) {
	filter := bson.M{
		"user_id":   userID,
		"device_id": deviceID,
		"active":    true,
	}
	
	// Get the most recent active signed prekey
	opts := options.FindOne().SetSort(bson.D{{Key: "created_time", Value: -1}})
	signedPrekey, err := mongoutil.FindOne[*signal.SignalSignedPreKey](ctx, s.coll, filter, opts)
	if err != nil {
		return nil, err
	}
	
	return signedPrekey, nil
}

func (s *SignedPreKeyStore) GetByKeyID(ctx context.Context, userID string, deviceID int32, keyID uint32) (*signal.SignalSignedPreKey, error) {
	filter := bson.M{
		"user_id":   userID,
		"device_id": deviceID,
		"key_id":    keyID,
	}
	
	signedPrekey, err := mongoutil.FindOne[*signal.SignalSignedPreKey](ctx, s.coll, filter)
	if err != nil {
		return nil, err
	}
	
	return signedPrekey, nil
}

func (s *SignedPreKeyStore) SetActive(ctx context.Context, userID string, deviceID int32, keyID uint32) error {
	session, err := s.coll.Database().Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)
	
	// Use transaction to ensure atomicity
	_, err = session.WithTransaction(ctx, func(sc mongo.SessionContext) (interface{}, error) {
		// First, deactivate all existing signed prekeys for this user/device
		deactivateFilter := bson.M{
			"user_id":   userID,
			"device_id": deviceID,
		}
		deactivateUpdate := bson.M{
			"$set": bson.M{"active": false},
		}
		_, err := s.coll.UpdateMany(sc, deactivateFilter, deactivateUpdate)
		if err != nil {
			return nil, err
		}
		
		// Then, activate the specified signed prekey
		activateFilter := bson.M{
			"user_id":   userID,
			"device_id": deviceID,
			"key_id":    keyID,
		}
		activateUpdate := bson.M{
			"$set": bson.M{"active": true},
		}
		return s.coll.UpdateOne(sc, activateFilter, activateUpdate)
	})
	
	return err
}

func (s *SignedPreKeyStore) Delete(ctx context.Context, userID string, deviceID int32, keyID uint32) error {
	filter := bson.M{
		"user_id":   userID,
		"device_id": deviceID,
		"key_id":    keyID,
	}
	
	return mongoutil.DeleteOne(ctx, s.coll, filter)
}

func (s *SignedPreKeyStore) GetAll(ctx context.Context, userID string, deviceID int32) ([]*signal.SignalSignedPreKey, error) {
	filter := bson.M{
		"user_id":   userID,
		"device_id": deviceID,
	}
	
	// Sort by creation time, newest first
	opts := options.Find().SetSort(bson.D{{Key: "created_time", Value: -1}})
	return mongoutil.Find[*signal.SignalSignedPreKey](ctx, s.coll, filter, opts)
}

func (s *SignedPreKeyStore) CleanupInactive(ctx context.Context, olderThan time.Duration) (int64, error) {
	cutoffTime := time.Now().Add(-olderThan)
	filter := bson.M{
		"active":       false,
		"created_time": bson.M{"$lt": cutoffTime},
	}
	
	result, err := s.coll.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}
	
	return result.DeletedCount, nil
}

func (s *SignedPreKeyStore) Exists(ctx context.Context, userID string, deviceID int32) (bool, error) {
	filter := bson.M{
		"user_id":   userID,
		"device_id": deviceID,
		"active":    true,
	}
	
	count, err := s.coll.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	
	return count > 0, nil
}