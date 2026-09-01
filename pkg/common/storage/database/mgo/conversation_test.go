package mgo

import (
	"context"
	"testing"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

type conversationVersionLogStub struct{}

func (conversationVersionLogStub) IncrVersion(context.Context, string, []string, int32) error {
	return nil
}

func (conversationVersionLogStub) FindChangeLog(context.Context, string, uint, int) (*model.VersionLog, error) {
	return nil, nil
}

func (conversationVersionLogStub) BatchFindChangeLog(context.Context, []string, []uint, []int) ([]*model.VersionLog, error) {
	return nil, nil
}

func (conversationVersionLogStub) DeleteAfterUnchangedLog(context.Context, time.Time) error {
	return nil
}

func (conversationVersionLogStub) Delete(context.Context, string) error {
	return nil
}

func TestConversationMgoCreateIgnoresDuplicateKey(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("duplicate conversation is idempotent", func(mt *mtest.T) {
		mt.AddMockResponses(bson.D{
			{Key: "ok", Value: 1},
			{Key: "writeErrors", Value: bson.A{
				bson.D{{Key: "index", Value: 0}, {Key: "code", Value: 11000}, {Key: "errmsg", Value: "duplicate key"}},
			}},
		})

		conversationDB := &ConversationMgo{coll: mt.Coll, version: conversationVersionLogStub{}}
		err := conversationDB.Create(context.Background(), []*model.Conversation{{
			OwnerUserID:    "user-1",
			ConversationID: "group-1",
		}})
		if err != nil {
			mt.Fatalf("expected duplicate conversation creation to be idempotent, got %v", err)
		}
	})
}

func TestConversationMgoCreatePropagatesNonDuplicateError(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("non-duplicate insert errors are returned", func(mt *mtest.T) {
		mt.AddMockResponses(bson.D{
			{Key: "ok", Value: 0},
			{Key: "code", Value: 13},
			{Key: "errmsg", Value: "permission denied"},
		})

		conversationDB := &ConversationMgo{coll: mt.Coll, version: conversationVersionLogStub{}}
		err := conversationDB.Create(context.Background(), []*model.Conversation{{
			OwnerUserID:    "user-1",
			ConversationID: "group-1",
		}})
		if err == nil {
			mt.Fatal("expected non-duplicate insert error to be returned")
		}
	})
}
