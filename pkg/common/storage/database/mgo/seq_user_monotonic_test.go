package mgo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestReadSeqMongoAtomicMax(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("out of order receipts use max without a preceding read", func(mt *mtest.T) {
		s := &seqUserMongo{coll: mt.Coll}
		for _, seq := range []int64{346, 200, 346} {
			mt.AddMockResponses(mtest.CreateSuccessResponse(bson.E{Key: "n", Value: 1}, bson.E{Key: "nModified", Value: 1}))
			require.NoError(mt, s.SetUserReadSeq(context.Background(), "group", "user", seq))
		}
		events := mt.GetAllStartedEvents()
		require.Len(mt, events, 3)
		for i, event := range events {
			require.Equal(mt, "update", event.CommandName)
			updates, err := event.Command.Lookup("updates").Array().Values()
			require.NoError(mt, err)
			update := updates[0].Document()
			require.Equal(mt, "group", update.Lookup("q", "conversation_id").StringValue())
			require.Equal(mt, "user", update.Lookup("q", "user_id").StringValue())
			require.Equal(mt, []int64{346, 200, 346}[i], update.Lookup("u", "$max", "read_seq").Int64())
			_, err = update.LookupErr("u", "$setOnInsert", "read_seq")
			require.Error(mt, err, "insertion defaults must not conflict with $max")
			require.True(mt, update.Lookup("upsert").Boolean())
		}
	})
	mt.Run("concurrent upsert retries the existing unique row", func(mt *mtest.T) {
		mt.AddMockResponses(
			mtest.CreateWriteErrorsResponse(mtest.WriteError{Index: 0, Code: 11000, Message: "duplicate user/conversation"}),
			mtest.CreateSuccessResponse(bson.E{Key: "n", Value: 1}, bson.E{Key: "nModified", Value: 1}),
		)
		s := &seqUserMongo{coll: mt.Coll}
		require.NoError(mt, s.SetUserReadSeq(context.Background(), "group", "user", 346))
		events := mt.GetAllStartedEvents()
		require.Len(mt, events, 2)
		updates, err := events[1].Command.Lookup("updates").Array().Values()
		require.NoError(mt, err)
		update := updates[0].Document()
		upsert, err := update.LookupErr("upsert")
		require.True(mt, err != nil || !upsert.Boolean())
		require.EqualValues(mt, 346, update.Lookup("u", "$max", "read_seq").Int64())
	})
	mt.Run("database failures propagate", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{Code: 91, Message: "unavailable"}))
		s := &seqUserMongo{coll: mt.Coll}
		require.Error(mt, s.SetUserReadSeq(context.Background(), "group", "user", 346))
	})
	mt.Run("negative sequences do not issue a command", func(mt *mtest.T) {
		s := &seqUserMongo{coll: mt.Coll}
		require.Error(mt, s.SetUserReadSeq(context.Background(), "group", "user", -1))
		require.Empty(mt, mt.GetAllStartedEvents())
	})
}

func TestReadSeqMongoUniqueIndex(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("one record per user and conversation", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateSuccessResponse())
		_, err := NewSeqUserMongo(mt.DB)
		require.NoError(mt, err)
		events := mt.GetAllStartedEvents()
		require.Len(mt, events, 1)
		require.Equal(mt, "createIndexes", events[0].CommandName)
		indexes, err := events[0].Command.Lookup("indexes").Array().Values()
		require.NoError(mt, err)
		index := indexes[0].Document()
		require.True(mt, index.Lookup("unique").Boolean())
		require.Equal(mt, "user_conversation_unique", index.Lookup("name").StringValue())
		require.EqualValues(mt, 1, index.Lookup("key", "user_id").Int32())
		require.EqualValues(mt, 1, index.Lookup("key", "conversation_id").Int32())
	})
	mt.Run("existing duplicates fail explicitly", func(mt *mtest.T) {
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{Code: 11000, Message: "duplicate rows"}))
		_, err := NewSeqUserMongo(mt.DB)
		require.ErrorContains(mt, err, "check duplicate user_id/conversation_id rows before rollout")
	})
}
