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
	"github.com/OpenIMSDK/tools/errs"
	"time"

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/tools/mgoutil"
	"github.com/OpenIMSDK/tools/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
)

func NewConversationMongo(db *mongo.Database) (*ConversationMgo, error) {
	coll := db.Collection("conversation")
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "owner_user_id", Value: 1},
			{Key: "conversation_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &ConversationMgo{coll: coll}, nil
}

type ConversationMgo struct {
	coll *mongo.Collection
}

func (c *ConversationMgo) Create(ctx context.Context, conversations []*relation.ConversationModel) (err error) {
	return mgoutil.InsertMany(ctx, c.coll, conversations)
}

func (c *ConversationMgo) Delete(ctx context.Context, groupIDs []string) (err error) {
	return mgoutil.DeleteMany(ctx, c.coll, bson.M{"group_id": bson.M{"$in": groupIDs}})
}

func (c *ConversationMgo) UpdateByMap(ctx context.Context, userIDs []string, conversationID string, args map[string]any) (rows int64, err error) {
	res, err := mgoutil.UpdateMany(ctx, c.coll, bson.M{"owner_user_id": bson.M{"$in": userIDs}, "conversation_id": conversationID}, bson.M{"$set": args})
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, nil
}

func (c *ConversationMgo) Update(ctx context.Context, conversation *relation.ConversationModel) (err error) {
	return mgoutil.UpdateOne(ctx, c.coll, bson.M{"owner_user_id": conversation.OwnerUserID, "conversation_id": conversation.ConversationID}, bson.M{"$set": conversation}, true)
}

func (c *ConversationMgo) Find(ctx context.Context, ownerUserID string, conversationIDs []string) (conversations []*relation.ConversationModel, err error) {
	return mgoutil.Find[*relation.ConversationModel](ctx, c.coll, bson.M{"owner_user_id": ownerUserID, "conversation_id": bson.M{"$in": conversationIDs}})
}

func (c *ConversationMgo) FindUserID(ctx context.Context, userIDs []string, conversationIDs []string) ([]string, error) {
	return mgoutil.Find[string](
		ctx,
		c.coll,
		bson.M{"owner_user_id": bson.M{"$in": userIDs}, "conversation_id": bson.M{"$in": conversationIDs}},
		options.Find().SetProjection(bson.M{"_id": 0, "owner_user_id": 1}),
	)
}

func (c *ConversationMgo) FindUserIDAllConversationID(ctx context.Context, userID string) ([]string, error) {
	return mgoutil.Find[string](ctx, c.coll, bson.M{"owner_user_id": userID}, options.Find().SetProjection(bson.M{"_id": 0, "conversation_id": 1}))
}

func (c *ConversationMgo) Take(ctx context.Context, userID, conversationID string) (conversation *relation.ConversationModel, err error) {
	return mgoutil.FindOne[*relation.ConversationModel](ctx, c.coll, bson.M{"owner_user_id": userID, "conversation_id": conversationID})
}

func (c *ConversationMgo) FindConversationID(ctx context.Context, userID string, conversationIDs []string) (existConversationID []string, err error) {
	return mgoutil.Find[string](ctx, c.coll, bson.M{"owner_user_id": userID, "conversation_id": bson.M{"$in": conversationIDs}}, options.Find().SetProjection(bson.M{"_id": 0, "conversation_id": 1}))
}

func (c *ConversationMgo) FindUserIDAllConversations(ctx context.Context, userID string) (conversations []*relation.ConversationModel, err error) {
	return mgoutil.Find[*relation.ConversationModel](ctx, c.coll, bson.M{"owner_user_id": userID})
}

func (c *ConversationMgo) FindRecvMsgNotNotifyUserIDs(ctx context.Context, groupID string) ([]string, error) {
	return mgoutil.Find[string](ctx, c.coll, bson.M{"group_id": groupID, "recv_msg_opt": constant.ReceiveNotNotifyMessage}, options.Find().SetProjection(bson.M{"_id": 0, "owner_user_id": 1}))
}

func (c *ConversationMgo) GetUserRecvMsgOpt(ctx context.Context, ownerUserID, conversationID string) (opt int, err error) {
	return mgoutil.FindOne[int](ctx, c.coll, bson.M{"owner_user_id": ownerUserID, "conversation_id": conversationID}, options.FindOne().SetProjection(bson.M{"recv_msg_opt": 1}))
}

func (c *ConversationMgo) GetAllConversationIDs(ctx context.Context) ([]string, error) {
	return mgoutil.Aggregate[string](ctx, c.coll, []bson.M{
		{"$group": bson.M{"_id": "$conversation_id"}},
		{"$project": bson.M{"_id": 0, "conversation_id": "$_id"}},
	})
}

func (c *ConversationMgo) GetAllConversationIDsNumber(ctx context.Context) (int64, error) {
	counts, err := mgoutil.Aggregate[int64](ctx, c.coll, []bson.M{
		{"$group": bson.M{"_id": "$conversation_id"}},
		{"$group": bson.M{"_id": nil, "count": bson.M{"$sum": 1}}},
		{"$project": bson.M{"_id": 0}},
	})
	if err != nil {
		return 0, err
	}
	if len(counts) == 0 {
		return 0, nil
	}
	return counts[0], nil
}

func (c *ConversationMgo) PageConversationIDs(ctx context.Context, pagination pagination.Pagination) (conversationIDs []string, err error) {
	return mgoutil.FindPageOnly[string](ctx, c.coll, bson.M{}, pagination, options.Find().SetProjection(bson.M{"conversation_id": 1}))
}

func (c *ConversationMgo) GetConversationsByConversationID(ctx context.Context, conversationIDs []string) ([]*relation.ConversationModel, error) {
	return mgoutil.Find[*relation.ConversationModel](ctx, c.coll, bson.M{"conversation_id": bson.M{"$in": conversationIDs}})
}

func (c *ConversationMgo) GetConversationIDsNeedDestruct(ctx context.Context) ([]*relation.ConversationModel, error) {
	//"is_msg_destruct = 1 && msg_destruct_time != 0 && (UNIX_TIMESTAMP(NOW()) > (msg_destruct_time + UNIX_TIMESTAMP(latest_msg_destruct_time)) || latest_msg_destruct_time is NULL)"
	return mgoutil.Find[*relation.ConversationModel](ctx, c.coll, bson.M{
		"is_msg_destruct":   1,
		"msg_destruct_time": bson.M{"$ne": 0},
		"$or": []bson.M{
			{
				"$expr": bson.M{
					"$gt": []any{
						time.Now(),
						bson.M{"$add": []any{"$msg_destruct_time", "$latest_msg_destruct_time"}},
					},
				},
			},
			{
				"latest_msg_destruct_time": nil,
			},
		},
	})
}

func (c *ConversationMgo) GetConversationNotReceiveMessageUserIDs(ctx context.Context, conversationID string) ([]string, error) {
	return mgoutil.Find[string](
		ctx,
		c.coll,
		bson.M{"conversation_id": conversationID, "recv_msg_opt": bson.M{"$ne": constant.ReceiveMessage}},
		options.Find().SetProjection(bson.M{"_id": 0, "owner_user_id": 1}),
	)
}
