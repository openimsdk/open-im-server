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

package conversation

import (
	"context"
	"errors"
	"testing"

	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/controller"
	relationtb "github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcli"
	"github.com/openimsdk/protocol/constant"
	pbconversation "github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/wrapperspb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

type mockGroupClient struct {
	group.GroupClient
	getGroupMemberUserIDsFunc func(ctx context.Context, in *group.GetGroupMemberUserIDsReq, opts ...grpc.CallOption) (*group.GetGroupMemberUserIDsResp, error)
}

func (m *mockGroupClient) GetGroupMemberUserIDs(ctx context.Context, in *group.GetGroupMemberUserIDsReq, opts ...grpc.CallOption) (*group.GetGroupMemberUserIDsResp, error) {
	if m.getGroupMemberUserIDsFunc != nil {
		return m.getGroupMemberUserIDsFunc(ctx, in, opts...)
	}
	return nil, nil
}

type mockConversationDatabase struct {
	controller.ConversationDatabase
	findConversationsFunc           func(ctx context.Context, ownerUserID string, conversationIDs []string) ([]*relationtb.Conversation, error)
	setUsersConversationFieldTxFunc func(ctx context.Context, userIDs []string, conversation *relationtb.Conversation, fieldMap map[string]any) error
}

func (m *mockConversationDatabase) FindConversations(ctx context.Context, ownerUserID string, conversationIDs []string) ([]*relationtb.Conversation, error) {
	if m.findConversationsFunc != nil {
		return m.findConversationsFunc(ctx, ownerUserID, conversationIDs)
	}
	return nil, nil
}

func (m *mockConversationDatabase) SetUsersConversationFieldTx(ctx context.Context, userIDs []string, conversation *relationtb.Conversation, fieldMap map[string]any) error {
	if m.setUsersConversationFieldTxFunc != nil {
		return m.setUsersConversationFieldTxFunc(ctx, userIDs, conversation, fieldMap)
	}
	return nil
}

func TestSetConversations_ReadGroupChatType_FilterNonMembers(t *testing.T) {
	ctx := authverify.WithTempAdmin(context.Background())
	groupID := "test_group_1"

	var updatedUserIDs []string
	var foundUserIDs []string

	mockDB := &mockConversationDatabase{
		findConversationsFunc: func(ctx context.Context, ownerUserID string, conversationIDs []string) ([]*relationtb.Conversation, error) {
			foundUserIDs = append(foundUserIDs, ownerUserID)
			return nil, nil // Return empty, so user goes to needUpdateUsersList
		},
		setUsersConversationFieldTxFunc: func(ctx context.Context, userIDs []string, conversation *relationtb.Conversation, fieldMap map[string]any) error {
			updatedUserIDs = append(updatedUserIDs, userIDs...)
			return nil
		},
	}

	mockGrp := &mockGroupClient{
		getGroupMemberUserIDsFunc: func(ctx context.Context, in *group.GetGroupMemberUserIDsReq, opts ...grpc.CallOption) (*group.GetGroupMemberUserIDsResp, error) {
			assert.Equal(t, groupID, in.GroupID)
			return &group.GetGroupMemberUserIDsResp{
				UserIDs: []string{"userA", "userB"},
			}, nil
		},
	}

	server := &conversationServer{
		conversationDatabase: mockDB,
		groupClient: &rpcli.GroupClient{
			GroupClient: mockGrp,
		},
	}

	req := &pbconversation.SetConversationsReq{
		UserIDs: []string{"userA", "userX"}, // userX is not in group
		Conversation: &pbconversation.ConversationReq{
			ConversationID:   "sg_" + groupID,
			ConversationType: constant.ReadGroupChatType,
			GroupID:          groupID,
			GroupAtType:      &wrapperspb.Int32Value{Value: constant.AtMe},
		},
	}

	resp, err := server.SetConversations(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// userX must be filtered out, only userA should be queried and updated
	assert.Equal(t, []string{"userA"}, foundUserIDs)
	assert.Equal(t, []string{"userA"}, updatedUserIDs)
}

func TestSetConversations_ReadGroupChatType_AllNonMembers(t *testing.T) {
	ctx := authverify.WithTempAdmin(context.Background())
	groupID := "test_group_1"

	dbCalled := false
	mockDB := &mockConversationDatabase{
		findConversationsFunc: func(ctx context.Context, ownerUserID string, conversationIDs []string) ([]*relationtb.Conversation, error) {
			dbCalled = true
			return nil, nil
		},
		setUsersConversationFieldTxFunc: func(ctx context.Context, userIDs []string, conversation *relationtb.Conversation, fieldMap map[string]any) error {
			dbCalled = true
			return nil
		},
	}

	mockGrp := &mockGroupClient{
		getGroupMemberUserIDsFunc: func(ctx context.Context, in *group.GetGroupMemberUserIDsReq, opts ...grpc.CallOption) (*group.GetGroupMemberUserIDsResp, error) {
			return &group.GetGroupMemberUserIDsResp{
				UserIDs: []string{"userA", "userB"},
			}, nil
		},
	}

	server := &conversationServer{
		conversationDatabase: mockDB,
		groupClient: &rpcli.GroupClient{
			GroupClient: mockGrp,
		},
	}

	req := &pbconversation.SetConversationsReq{
		UserIDs: []string{"userX", "userY"}, // Neither is a member
		Conversation: &pbconversation.ConversationReq{
			ConversationID:   "sg_" + groupID,
			ConversationType: constant.ReadGroupChatType,
			GroupID:          groupID,
			GroupAtType:      &wrapperspb.Int32Value{Value: constant.AtMe},
		},
	}

	resp, err := server.SetConversations(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Database should NOT be touched at all
	assert.False(t, dbCalled)
}

func TestSetConversations_ReadGroupChatType_ExtractGroupIDFromConversationID(t *testing.T) {
	ctx := authverify.WithTempAdmin(context.Background())
	groupID := "test_group_extract"

	var queriedGroupID string
	mockGrp := &mockGroupClient{
		getGroupMemberUserIDsFunc: func(ctx context.Context, in *group.GetGroupMemberUserIDsReq, opts ...grpc.CallOption) (*group.GetGroupMemberUserIDsResp, error) {
			queriedGroupID = in.GroupID
			return &group.GetGroupMemberUserIDsResp{
				UserIDs: []string{"userA"},
			}, nil
		},
	}

	var updatedUserIDs []string
	mockDB := &mockConversationDatabase{
		findConversationsFunc: func(ctx context.Context, ownerUserID string, conversationIDs []string) ([]*relationtb.Conversation, error) {
			return nil, nil
		},
		setUsersConversationFieldTxFunc: func(ctx context.Context, userIDs []string, conversation *relationtb.Conversation, fieldMap map[string]any) error {
			updatedUserIDs = append(updatedUserIDs, userIDs...)
			return nil
		},
	}

	server := &conversationServer{
		conversationDatabase: mockDB,
		groupClient: &rpcli.GroupClient{
			GroupClient: mockGrp,
		},
	}

	// GroupID is left empty, but ConversationID is "sg_test_group_extract"
	req := &pbconversation.SetConversationsReq{
		UserIDs: []string{"userA", "userX"},
		Conversation: &pbconversation.ConversationReq{
			ConversationID:   "sg_" + groupID,
			ConversationType: constant.ReadGroupChatType,
			GroupAtType:      &wrapperspb.Int32Value{Value: constant.AtMe},
		},
	}

	resp, err := server.SetConversations(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, groupID, queriedGroupID)
	assert.Equal(t, []string{"userA"}, updatedUserIDs)
}

func TestSetConversations_SingleChat_NoGroupMemberCheck(t *testing.T) {
	ctx := authverify.WithTempAdmin(context.Background())

	groupClientCalled := false
	mockGrp := &mockGroupClient{
		getGroupMemberUserIDsFunc: func(ctx context.Context, in *group.GetGroupMemberUserIDsReq, opts ...grpc.CallOption) (*group.GetGroupMemberUserIDsResp, error) {
			groupClientCalled = true
			return nil, nil
		},
	}

	var updatedUserIDs []string
	mockDB := &mockConversationDatabase{
		findConversationsFunc: func(ctx context.Context, ownerUserID string, conversationIDs []string) ([]*relationtb.Conversation, error) {
			return nil, nil
		},
		setUsersConversationFieldTxFunc: func(ctx context.Context, userIDs []string, conversation *relationtb.Conversation, fieldMap map[string]any) error {
			updatedUserIDs = append(updatedUserIDs, userIDs...)
			return nil
		},
	}

	server := &conversationServer{
		conversationDatabase: mockDB,
		groupClient: &rpcli.GroupClient{
			GroupClient: mockGrp,
		},
	}

	req := &pbconversation.SetConversationsReq{
		UserIDs: []string{"userA"},
		Conversation: &pbconversation.ConversationReq{
			ConversationID:   "si_userA_userB",
			ConversationType: constant.SingleChatType,
			RecvMsgOpt:       &wrapperspb.Int32Value{Value: constant.ReceiveMessage},
		},
	}

	resp, err := server.SetConversations(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.False(t, groupClientCalled)
	assert.Equal(t, []string{"userA"}, updatedUserIDs)
}

func TestSetConversations_NilConversation(t *testing.T) {
	ctx := authverify.WithTempAdmin(context.Background())
	server := &conversationServer{}

	req := &pbconversation.SetConversationsReq{
		UserIDs:      []string{"userA"},
		Conversation: nil,
	}

	resp, err := server.SetConversations(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestSetConversations_EmptyUserIDs(t *testing.T) {
	ctx := authverify.WithTempAdmin(context.Background())
	server := &conversationServer{}

	req := &pbconversation.SetConversationsReq{
		UserIDs: []string{},
		Conversation: &pbconversation.ConversationReq{
			ConversationID:   "sg_group1",
			ConversationType: constant.ReadGroupChatType,
		},
	}

	resp, err := server.SetConversations(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestSetConversations_GetGroupMemberUserIDsError(t *testing.T) {
	ctx := authverify.WithTempAdmin(context.Background())
	groupID := "test_group_err"

	mockGrp := &mockGroupClient{
		getGroupMemberUserIDsFunc: func(ctx context.Context, in *group.GetGroupMemberUserIDsReq, opts ...grpc.CallOption) (*group.GetGroupMemberUserIDsResp, error) {
			return nil, errors.New("rpc error: failed to fetch members")
		},
	}

	server := &conversationServer{
		groupClient: &rpcli.GroupClient{
			GroupClient: mockGrp,
		},
	}

	req := &pbconversation.SetConversationsReq{
		UserIDs: []string{"userA"},
		Conversation: &pbconversation.ConversationReq{
			ConversationID:   "sg_" + groupID,
			ConversationType: constant.ReadGroupChatType,
			GroupID:          groupID,
			GroupAtType:      &wrapperspb.Int32Value{Value: constant.AtMe},
		},
	}

	resp, err := server.SetConversations(ctx, req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}
