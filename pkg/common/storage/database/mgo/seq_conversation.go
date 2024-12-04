package mgo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/openimsdk/tools/db/mongoutil"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

func NewSeqConversationMongo(db *mongo.Database) (database.SeqConversation, error) {
	coll := db.Collection(database.SeqConversationName)
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "conversation_id", Value: 1},
		},
	})
	if err != nil {
		return nil, err
	}
	return &seqConversationMongo{coll: coll}, nil
}

type seqConversationMongo struct {
	coll *mongo.Collection
}

func (s *seqConversationMongo) setSeq(ctx context.Context, conversationID string, seq int64, field string) error {
	filter := map[string]any{
		"conversation_id": conversationID,
	}
	insert := bson.M{
		"conversation_id": conversationID,
		"min_seq":         0,
		"max_seq":         0,
	}
	delete(insert, field)
	update := map[string]any{
		"$set": bson.M{
			field: seq,
		},
		"$setOnInsert": insert,
	}
	opt := options.Update().SetUpsert(true)
	return mongoutil.UpdateOne(ctx, s.coll, filter, update, false, opt)
}

func (s *seqConversationMongo) Malloc(ctx context.Context, conversationID string, size int64) (int64, error) {
	if size < 0 {
		return 0, errors.New("size must be greater than 0")
	}
	if size == 0 {
		return s.GetMaxSeq(ctx, conversationID)
	}
	filter := map[string]any{"conversation_id": conversationID}
	update := map[string]any{
		"$inc": map[string]any{"max_seq": size},
		"$set": map[string]any{"min_seq": int64(0)},
	}
	opt := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After).SetProjection(map[string]any{"_id": 0, "max_seq": 1})
	lastSeq, err := mongoutil.FindOneAndUpdate[int64](ctx, s.coll, filter, update, opt)
	if err != nil {
		return 0, err
	}
	return lastSeq - size, nil
}

func (s *seqConversationMongo) SetMaxSeq(ctx context.Context, conversationID string, seq int64) error {
	return s.setSeq(ctx, conversationID, seq, "max_seq")
}

func (s *seqConversationMongo) GetMaxSeq(ctx context.Context, conversationID string) (int64, error) {
	seq, err := mongoutil.FindOne[int64](ctx, s.coll, bson.M{"conversation_id": conversationID}, options.FindOne().SetProjection(map[string]any{"_id": 0, "max_seq": 1}))
	if err == nil {
		return seq, nil
	} else if IsNotFound(err) {
		return 0, nil
	} else {
		return 0, err
	}
}

func (s *seqConversationMongo) GetMinSeq(ctx context.Context, conversationID string) (int64, error) {
	seq, err := mongoutil.FindOne[int64](ctx, s.coll, bson.M{"conversation_id": conversationID}, options.FindOne().SetProjection(map[string]any{"_id": 0, "min_seq": 1}))
	if err == nil {
		return seq, nil
	} else if IsNotFound(err) {
		return 0, nil
	} else {
		return 0, err
	}
}

func (s *seqConversationMongo) SetMinSeq(ctx context.Context, conversationID string, seq int64) error {
	return s.setSeq(ctx, conversationID, seq, "min_seq")
}

func (s *seqConversationMongo) GetConversation(ctx context.Context, conversationID string) (*model.SeqConversation, error) {
	return mongoutil.FindOne[*model.SeqConversation](ctx, s.coll, bson.M{"conversation_id": conversationID})
}
