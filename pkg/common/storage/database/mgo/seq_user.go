package mgo

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/tools/db/mongoutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewSeqUserMongo(db *mongo.Database) (database.SeqUser, error) {
	coll := db.Collection(database.SeqConversationName)
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
			{Key: "conversation_id", Value: 1},
		},
	})
	if err != nil {
		return nil, err
	}
	return &seqUserMongo{coll: coll}, nil
}

type seqUserMongo struct {
	coll *mongo.Collection
}

func (s *seqUserMongo) setSeq(ctx context.Context, userID string, conversationID string, seq int64, field string) error {
	filter := map[string]any{
		"user_id":         userID,
		"conversation_id": conversationID,
	}
	update := map[string]any{
		"$set": map[string]any{"field": int64(0)},
	}
	opt := options.Update().SetUpsert(true)
	return mongoutil.UpdateOne(ctx, s.coll, filter, update, false, opt)
}

func (s *seqUserMongo) GetMaxSeq(ctx context.Context, userID string, conversationID string) (int64, error) {

	//TODO implement me
	panic("implement me")
}

func (s *seqUserMongo) SetMaxSeq(ctx context.Context, userID string, conversationID string, seq int64) error {
	//TODO implement me
	panic("implement me")
}

func (s *seqUserMongo) GetMinSeq(ctx context.Context, userID string, conversationID string) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (s *seqUserMongo) SetMinSeq(ctx context.Context, userID string, conversationID string, seq int64) error {
	//TODO implement me
	panic("implement me")
}

func (s *seqUserMongo) GetReadSeq(ctx context.Context, userID string, conversationID string) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (s *seqUserMongo) SetReadSeq(ctx context.Context, userID string, conversationID string, seq int64) error {
	//TODO implement me
	panic("implement me")
}
