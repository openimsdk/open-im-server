package mgo

import (
	"context"
	"errors"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/db/mongoutil"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (s *seqConversationMongo) MallocSeq(ctx context.Context, conversationID string, size int64) ([]int64, error) {
	first, err := s.Malloc(ctx, conversationID, size)
	if err != nil {
		return nil, err
	}
	seqs := make([]int64, 0, size)
	for i := int64(0); i < size; i++ {
		seqs = append(seqs, first+i+1)
	}
	return seqs, nil
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
	return mongoutil.UpdateOne(ctx, s.coll, bson.M{"conversation_id": conversationID}, bson.M{"$set": bson.M{"min_seq": seq}}, false)
}

func (s *seqConversationMongo) GetConversation(ctx context.Context, conversationID string) (*model.SeqConversation, error) {
	return mongoutil.FindOne[*model.SeqConversation](ctx, s.coll, bson.M{"conversation_id": conversationID})
}
