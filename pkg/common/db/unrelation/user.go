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

package unrelation

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/unrelation"
	"github.com/OpenIMSDK/tools/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

//  prefixes and suffixes.
const (
	SubscriptionPrefix = "subscription_prefix"
	SubscribedPrefix   = "subscribed_prefix"
)

// MaximumSubscription Maximum number of subscriptions.
const (
	MaximumSubscription = 3000
)

func NewUserMongoDriver(database *mongo.Database) unrelation.UserModelInterface {
	return &UserMongoDriver{
		userCollection: database.Collection(unrelation.SubscribeUser),
	}
}

type UserMongoDriver struct {
	userCollection *mongo.Collection
}

// AddSubscriptionList Subscriber's handling of thresholds.
func (u *UserMongoDriver) AddSubscriptionList(ctx context.Context, userID string, userIDList []string) error {
	// Check the number of lists in the key.
	filter := bson.M{SubscriptionPrefix + userID: bson.M{"$size": 1}}
	result, err := u.userCollection.Find(context.Background(), filter)
	if err != nil {
		return err
	}
	var newUserIDList []string
	for result.Next(context.Background()) {
		err := result.Decode(&newUserIDList)
		if err != nil {
			log.Fatal(err)
		}
	}
	// If the threshold is exceeded, pop out the previous MaximumSubscription - len(userIDList) and insert it.
	if len(newUserIDList)+len(userIDList) > MaximumSubscription {
		newUserIDList = newUserIDList[MaximumSubscription-len(userIDList):]
		_, err := u.userCollection.UpdateOne(
			ctx,
			bson.M{"user_id": SubscriptionPrefix + userID},
			bson.M{"$set": bson.M{"user_id_list": newUserIDList}},
		)
		if err != nil {
			return err
		}
		//for i := 1; i <= MaximumSubscription-len(userIDList); i++ {
		//	_, err := u.userCollection.UpdateOne(
		//		ctx,
		//		bson.M{"user_id": SubscriptionPrefix + userID},
		//		bson.M{SubscriptionPrefix + userID: bson.M{"$pop": -1}},
		//	)
		//	if err != nil {
		//		return err
		//	}
		//}
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
		return err
	}
	for _, user := range userIDList {
		_, err = u.userCollection.UpdateOne(
			ctx,
			bson.M{"user_id": SubscribedPrefix + user},
			bson.M{"$addToSet": bson.M{"user_id_list": userID}},
			opts,
		)
		if err != nil {
			return utils.Wrap(err, "transaction failed")
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
		return err
	}
	err = u.RemoveSubscribedListFromUser(ctx, userID, userIDList)
	if err != nil {
		return err
	}
	return nil
}

// RemoveSubscribedListFromUser Among the unsubscribed users, delete the user from the subscribed list.
func (u *UserMongoDriver) RemoveSubscribedListFromUser(ctx context.Context, userID string, userIDList []string) error {
	var newUserIDList []string
	for _, value := range userIDList {
		newUserIDList = append(newUserIDList, SubscribedPrefix+value)
	}
	_, err := u.userCollection.UpdateOne(
		ctx,
		bson.M{"user_id": bson.M{"$in": newUserIDList}},
		bson.M{"$pull": bson.M{"user_id_list": userID}},
	)
	return utils.Wrap(err, "")
}
