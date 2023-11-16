package newmgo

import (
	"context"
	"github.com/OpenIMSDK/protocol/constant"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/newmgo/mgotool"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
	"github.com/openimsdk/open-im-server/v3/pkg/common/pagination"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func NewConversationMongo(db *mongo.Database) (*ConversationMgo, error) {
	return &ConversationMgo{
		coll: db.Collection("conversation"),
	}, nil
}

type ConversationMgo struct {
	coll *mongo.Collection
}

func (c *ConversationMgo) Create(ctx context.Context, conversations []*relation.ConversationModel) (err error) {
	return mgotool.InsertMany(ctx, c.coll, conversations)
}

func (c *ConversationMgo) Delete(ctx context.Context, groupIDs []string) (err error) {
	return mgotool.DeleteMany(ctx, c.coll, bson.M{"group_id": bson.M{"$in": groupIDs}})
}

func (c *ConversationMgo) UpdateByMap(ctx context.Context, userIDs []string, conversationID string, args map[string]any) (rows int64, err error) {
	res, err := mgotool.UpdateMany(ctx, c.coll, bson.M{"owner_user_id": bson.M{"$in": userIDs}, "conversation_id": conversationID}, bson.M{"$set": args})
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, nil
}

func (c *ConversationMgo) Update(ctx context.Context, conversation *relation.ConversationModel) (err error) {
	return mgotool.UpdateOne(ctx, c.coll, bson.M{"owner_user_id": conversation.OwnerUserID, "conversation_id": conversation.ConversationID}, bson.M{"$set": conversation}, true)
}

func (c *ConversationMgo) Find(ctx context.Context, ownerUserID string, conversationIDs []string) (conversations []*relation.ConversationModel, err error) {
	return mgotool.Find[*relation.ConversationModel](ctx, c.coll, bson.M{"owner_user_id": ownerUserID, "conversation_id": bson.M{"$in": conversationIDs}})
}

func (c *ConversationMgo) FindUserID(ctx context.Context, userIDs []string, conversationIDs []string) ([]string, error) {
	return mgotool.Find[string](ctx, c.coll, bson.M{"owner_user_id": bson.M{"$in": userIDs}, "conversation_id": bson.M{"$in": conversationIDs}}, options.Find().SetProjection(bson.M{"owner_user_id": 1}))
}

func (c *ConversationMgo) FindUserIDAllConversationID(ctx context.Context, userID string) ([]string, error) {
	return mgotool.Find[string](ctx, c.coll, bson.M{"owner_user_id": userID}, options.Find().SetProjection(bson.M{"conversation_id": 1}))
}

func (c *ConversationMgo) Take(ctx context.Context, userID, conversationID string) (conversation *relation.ConversationModel, err error) {
	return mgotool.FindOne[*relation.ConversationModel](ctx, c.coll, bson.M{"owner_user_id": userID, "conversation_id": conversationID})
}

func (c *ConversationMgo) FindConversationID(ctx context.Context, userID string, conversationIDs []string) (existConversationID []string, err error) {
	return mgotool.Find[string](ctx, c.coll, bson.M{"owner_user_id": userID, "conversation_id": bson.M{"$in": conversationIDs}}, options.Find().SetProjection(bson.M{"conversation_id": 1}))
}

func (c *ConversationMgo) FindUserIDAllConversations(ctx context.Context, userID string) (conversations []*relation.ConversationModel, err error) {
	return mgotool.Find[*relation.ConversationModel](ctx, c.coll, bson.M{"owner_user_id": userID})
}

func (c *ConversationMgo) FindRecvMsgNotNotifyUserIDs(ctx context.Context, groupID string) ([]string, error) {
	return mgotool.Find[string](ctx, c.coll, bson.M{"group_id": groupID, "recv_msg_opt": constant.ReceiveNotNotifyMessage}, options.Find().SetProjection(bson.M{"owner_user_id": 1}))
}

func (c *ConversationMgo) GetUserRecvMsgOpt(ctx context.Context, ownerUserID, conversationID string) (opt int, err error) {
	return mgotool.FindOne[int](ctx, c.coll, bson.M{"owner_user_id": ownerUserID, "conversation_id": conversationID}, options.FindOne().SetProjection(bson.M{"recv_msg_opt": 1}))
}

func (c *ConversationMgo) GetAllConversationIDs(ctx context.Context) ([]string, error) {
	return mgotool.Aggregate[string](ctx, c.coll, []bson.M{
		{"$group": bson.M{"_id": "$conversation_id"}},
		{"$project": bson.M{"_id": 0, "conversation_id": "$_id"}},
	})
}

func (c *ConversationMgo) GetAllConversationIDsNumber(ctx context.Context) (int64, error) {
	counts, err := mgotool.Aggregate[int64](ctx, c.coll, []bson.M{
		{"$group": bson.M{"_id": "$conversation_id"}},
		{"$project": bson.M{"_id": 0, "conversation_id": "$_id"}},
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
	return mgotool.FindPageOnly[string](ctx, c.coll, bson.M{}, pagination, options.Find().SetProjection(bson.M{"conversation_id": 1}))
}

func (c *ConversationMgo) GetConversationsByConversationID(ctx context.Context, conversationIDs []string) ([]*relation.ConversationModel, error) {
	return mgotool.Find[*relation.ConversationModel](ctx, c.coll, bson.M{"conversation_id": bson.M{"$in": conversationIDs}})
}

func (c *ConversationMgo) GetConversationIDsNeedDestruct(ctx context.Context) ([]*relation.ConversationModel, error) {
	//"is_msg_destruct = 1 && msg_destruct_time != 0 && (UNIX_TIMESTAMP(NOW()) > (msg_destruct_time + UNIX_TIMESTAMP(latest_msg_destruct_time)) || latest_msg_destruct_time is NULL)"
	return mgotool.Find[*relation.ConversationModel](ctx, c.coll, bson.M{
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
	return mgotool.Find[string](ctx, c.coll, bson.M{"conversation_id": conversationID, "recv_msg_opt": bson.M{"$ne": constant.ReceiveMessage}}, options.Find().SetProjection(bson.M{"owner_user_id": 1}))
}

func (c *ConversationMgo) NewTx(tx any) relation.ConversationModelInterface {
	//TODO implement me
	panic("implement me")
}
