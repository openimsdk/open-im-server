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

package mgo

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"

	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// prefixes and suffixes.
const (
	SubscriptionPrefix = "subscription_prefix"
	SubscribedPrefix   = "subscribed_prefix"
)

// MaximumSubscription Maximum number of subscriptions.
const (
	MaximumSubscription = 3000
)

func NewUserMongoDriver(database *mongo.Database) database.SubscribeUser {
	return &UserMongoDriver{
		userCollection: database.Collection(model.SubscribeUserTableName),
	}
}

type UserMongoDriver struct {
	userCollection *mongo.Collection
}

// AddSubscriptionList Subscriber's handling of thresholds.
func (u *UserMongoDriver) AddSubscriptionList(ctx context.Context, userID string, userIDList []string) error {
	// Check the number of lists in the key.
	pipeline := mongo.Pipeline{
		{{"$match", bson.D{{"user_id", SubscriptionPrefix + userID}}}},
		{{"$project", bson.D{{"count", bson.D{{"$size", "$user_id_list"}}}}}},
	}
	// perform aggregate operations
	cursor, err := u.userCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return errs.Wrap(err)
	}
	defer cursor.Close(ctx)
	var cnt struct {
		Count int `bson:"count"`
	}
	// iterate over aggregated results
	for cursor.Next(ctx) {
		err = cursor.Decode(&cnt)
		if err != nil {
			return errs.Wrap(err)
		}
	}
	var newUserIDList []string
	// If the threshold is exceeded, pop out the previous MaximumSubscription - len(userIDList) and insert it.
	if cnt.Count+len(userIDList) > MaximumSubscription {
		newUserIDList, err = u.GetAllSubscribeList(ctx, userID)
		if err != nil {
			return err
		}
		newUserIDList = newUserIDList[MaximumSubscription-len(userIDList):]
		_, err = u.userCollection.UpdateOne(
			ctx,
			bson.M{"user_id": SubscriptionPrefix + userID},
			bson.M{"$set": bson.M{"user_id_list": newUserIDList}},
		)
		if err != nil {
			return err
		}
		// Another way to subscribe to N before pop,Delete after testing
		/*for i := 1; i <= MaximumSubscription-len(userIDList); i++ {
			_, err := u.userCollection.UpdateOne(
				ctx,
				bson.M{"user_id": SubscriptionPrefix + userID},
				bson.M{SubscriptionPrefix + userID: bson.M{"$pop": -1}},
			)
			if err != nil {
				return err
			}
		}*/
	}
	upsert := true
	opts := &options.UpdateOptions{
		Upsert: &upsert,
	}
	_, err = u.userCollection.UpdateOne(
		ctx,
		bson.M{"user_id": SubscriptionPrefix + userID},
		bson.M{"$addToSet": bson.M{"user_id_list": bson.M{"$each": userIDList}}},
		opts,
	)
	if err != nil {
		return errs.Wrap(err)
	}
	for _, user := range userIDList {
		_, err = u.userCollection.UpdateOne(
			ctx,
			bson.M{"user_id": SubscribedPrefix + user},
			bson.M{"$addToSet": bson.M{"user_id_list": userID}},
			opts,
		)
		if err != nil {
			return errs.WrapMsg(err, "transaction failed")
		}
	}
	return nil
}

// UnsubscriptionList Handling of unsubscribe.
func (u *UserMongoDriver) UnsubscriptionList(ctx context.Context, userID string, userIDList []string) error {
	_, err := u.userCollection.UpdateOne(
		ctx,
		bson.M{"user_id": SubscriptionPrefix + userID},
		bson.M{"$pull": bson.M{"user_id_list": bson.M{"$in": userIDList}}},
	)
	if err != nil {
		return errs.Wrap(err)
	}
	err = u.RemoveSubscribedListFromUser(ctx, userID, userIDList)
	if err != nil {
		return errs.Wrap(err)
	}
	return nil
}

// RemoveSubscribedListFromUser Among the unsubscribed users, delete the user from the subscribed list.
func (u *UserMongoDriver) RemoveSubscribedListFromUser(ctx context.Context, userID string, userIDList []string) error {
	var err error
	for _, userIDTemp := range userIDList {
		_, err = u.userCollection.UpdateOne(
			ctx,
			bson.M{"user_id": SubscribedPrefix + userIDTemp},
			bson.M{"$pull": bson.M{"user_id_list": userID}},
		)
	}
	return errs.Wrap(err)
}

// GetAllSubscribeList Get all users subscribed by this user.
func (u *UserMongoDriver) GetAllSubscribeList(ctx context.Context, userID string) (userIDList []string, err error) {
	var user model.SubscribeUser
	cursor := u.userCollection.FindOne(
		ctx,
		bson.M{"user_id": SubscriptionPrefix + userID})
	err = cursor.Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []string{}, nil
		} else {
			return nil, errs.Wrap(err)
		}
	}
	return user.UserIDList, nil
}

// GetSubscribedList Get the user subscribed by those users.
func (u *UserMongoDriver) GetSubscribedList(ctx context.Context, userID string) (userIDList []string, err error) {
	var user model.SubscribeUser
	cursor := u.userCollection.FindOne(
		ctx,
		bson.M{"user_id": SubscribedPrefix + userID})
	err = cursor.Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return []string{}, nil
		} else {
			return nil, errs.Wrap(err)
		}
	}
	return user.UserIDList, nil
}
