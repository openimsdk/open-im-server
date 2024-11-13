package mgo

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/errs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func NewStreamMsgMongo(db *mongo.Database) (*StreamMsgMongo, error) {
	coll := db.Collection(database.StreamMsgName)
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "client_msg_id", Value: 1},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	return &StreamMsgMongo{coll: coll}, nil
}

type StreamMsgMongo struct {
	coll *mongo.Collection
}

func (m *StreamMsgMongo) CreateStreamMsg(ctx context.Context, val *model.StreamMsg) error {
	if val.Packets == nil {
		val.Packets = []string{}
	}
	return mongoutil.InsertMany(ctx, m.coll, []*model.StreamMsg{val})
}

func (m *StreamMsgMongo) AppendStreamMsg(ctx context.Context, clientMsgID string, startIndex int, packets []string, end bool, deadlineTime time.Time) error {
	update := bson.M{
		"$set": bson.M{
			"end":           end,
			"deadline_time": deadlineTime,
		},
	}
	if len(packets) > 0 {
		update["$push"] = bson.M{
			"packets": bson.M{
				"$each":     packets,
				"$position": startIndex,
			},
		}
	}
	return mongoutil.UpdateOne(ctx, m.coll, bson.M{"client_msg_id": clientMsgID, "end": false}, update, true)
}

func (m *StreamMsgMongo) GetStreamMsg(ctx context.Context, clientMsgID string) (*model.StreamMsg, error) {
	return mongoutil.FindOne[*model.StreamMsg](ctx, m.coll, bson.M{"client_msg_id": clientMsgID})
}
