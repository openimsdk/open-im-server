package controller

import (
	"context"
	"errors"
	"testing"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	rediscache "github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/redis"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type repairConversationStore struct {
	database.Conversation
	created  []*model.Conversation
	existing *model.Conversation
}

func (*repairConversationStore) FindUserID(context.Context, []string, []string) ([]string, error) {
	return nil, nil
}

func (s *repairConversationStore) Find(context.Context, string, []string) ([]*model.Conversation, error) {
	if s.existing == nil {
		return nil, nil
	}
	return []*model.Conversation{s.existing}, nil
}

func (s *repairConversationStore) Create(_ context.Context, conversations []*model.Conversation) error {
	s.created = conversations
	return nil
}

type repairConversationTx struct {
	commitErr error
	committed bool
}

func (tx *repairConversationTx) Transaction(ctx context.Context, fn func(context.Context) error) error {
	if err := fn(ctx); err != nil {
		return err
	}
	tx.committed = tx.commitErr == nil
	return tx.commitErr
}

type repairConversationDeleter struct {
	cache.BatchDeleter
	keys     []string
	onDelete func([]string)
}

func (d *repairConversationDeleter) Clone() cache.BatchDeleter {
	return &repairConversationDeleter{keys: append([]string(nil), d.keys...), onDelete: d.onDelete}
}
func (d *repairConversationDeleter) AddKeys(keys ...string) { d.keys = append(d.keys, keys...) }
func (d *repairConversationDeleter) ChainExecDel(context.Context) error {
	d.onDelete(d.keys)
	return nil
}

func TestCreateMissingConversationsInvalidatesOwnerKeysAfterCommit(t *testing.T) {
	for _, commitFails := range []bool{false, true} {
		name := "committed"
		if commitFails {
			name = "failed commit"
		}
		t.Run(name, func(t *testing.T) {
			transaction := &repairConversationTx{}
			if commitFails {
				transaction.commitErr = errors.New("commit failed")
			}
			var deleted []string
			var deletedAfterCommit bool
			conversationCache := &rediscache.ConversationRedisCache{BatchDeleter: &repairConversationDeleter{
				onDelete: func(keys []string) {
					deletedAfterCommit = transaction.committed
					deleted = append(deleted, keys...)
				},
			}}
			store := &repairConversationStore{}
			database := NewConversationDatabase(store, conversationCache, transaction)
			id := "si_alice_bob"
			owners := []string{"alice", "bob"}
			err := database.SetUsersConversationFieldTx(context.Background(), owners, &model.Conversation{ConversationID: id}, map[string]any{"recv_msg_opt": 0})
			if commitFails {
				require.ErrorIs(t, err, transaction.commitErr)
				require.Empty(t, deleted)
				return
			}
			require.NoError(t, err)
			require.Len(t, store.created, 2)
			for _, owner := range owners {
				assert.Contains(t, deleted, cachekey.GetConversationKey(owner, id))
				assert.NotContains(t, deleted, cachekey.GetConversationKey(id, owner))
			}
			require.True(t, deletedAfterCommit, "cache deletion must happen after the transaction commits")
		})
	}
}

func TestCreateConversationIfMissingRepairsStaleCacheWithoutRewritingOwner(t *testing.T) {
	existing := &model.Conversation{OwnerUserID: "bob", ConversationID: "si_alice_bob", IsPinned: true}
	store := &repairConversationStore{existing: existing}
	transaction := &repairConversationTx{}
	var deleted []string
	conversationCache := &rediscache.ConversationRedisCache{BatchDeleter: &repairConversationDeleter{
		onDelete: func(keys []string) {
			require.True(t, transaction.committed)
			deleted = append(deleted, keys...)
		},
	}}
	db := NewConversationDatabase(store, conversationCache, transaction)
	created, err := db.CreateConversationIfMissing(context.Background(), &model.Conversation{
		OwnerUserID:    existing.OwnerUserID,
		ConversationID: existing.ConversationID,
	})
	require.NoError(t, err)
	require.False(t, created)
	require.Empty(t, store.created)
	require.True(t, existing.IsPinned)
	require.Contains(t, deleted, cachekey.GetConversationKey(existing.OwnerUserID, existing.ConversationID))
	require.Contains(t, deleted, cachekey.GetPinnedConversationIDs(existing.OwnerUserID))
}
