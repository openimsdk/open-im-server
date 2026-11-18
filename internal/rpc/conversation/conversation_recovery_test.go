package conversation

import (
	"context"
	"errors"
	"testing"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/protocol/constant"
	pbconv "github.com/openimsdk/protocol/conversation"
	"github.com/stretchr/testify/require"
)

type partialConversationDatabase struct {
	controller.ConversationDatabase
	existing map[string]*model.Conversation
	attempts []string
	failure  error
}

func (d *partialConversationDatabase) FindConversations(_ context.Context, owner string, _ []string) ([]*model.Conversation, error) {
	if conversation := d.existing[owner]; conversation != nil {
		return []*model.Conversation{conversation}, nil
	}
	return nil, nil
}

func (d *partialConversationDatabase) CreateConversationIfMissing(_ context.Context, conversation *model.Conversation) (bool, error) {
	owner := conversation.OwnerUserID
	d.attempts = append(d.attempts, owner)
	if owner == "bob" && d.failure != nil {
		return false, d.failure
	}
	d.existing[owner] = conversation
	return true, nil
}

func TestCreateSingleChatConversationsRepairsMissingOwner(t *testing.T) {
	ctx := context.Background()
	db := &partialConversationDatabase{existing: make(map[string]*model.Conversation), failure: errors.New("could not create receiver conversation")}
	server := &conversationServer{config: &Config{}, conversationDatabase: db}
	req := &pbconv.CreateSingleChatConversationsReq{SendID: "alice", RecvID: "bob", ConversationID: "si_alice_bob", ConversationType: constant.SingleChatType}
	_, err := server.CreateSingleChatConversations(ctx, req)
	require.ErrorIs(t, err, db.failure)
	require.Equal(t, []string{"alice", "bob"}, db.attempts)
	sender := db.existing["alice"]
	require.NotNil(t, sender)
	sender.IsPinned = true
	require.Nil(t, db.existing["bob"])

	db.failure = nil
	_, err = server.CreateSingleChatConversations(ctx, req)
	require.NoError(t, err)
	require.Equal(t, []string{"alice", "bob", "bob"}, db.attempts)
	require.Same(t, sender, db.existing["alice"])
	require.True(t, db.existing["alice"].IsPinned)
	require.Equal(t, "alice", db.existing["bob"].UserID)
	require.False(t, db.existing["bob"].IsPinned)
}
