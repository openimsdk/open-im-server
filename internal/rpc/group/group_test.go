package group

import (
	"context"
	"errors"
	"testing"

	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/open-im-server/v3/pkg/common/webhook"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcli"
	pbconv "github.com/openimsdk/protocol/conversation"
	pbgroup "github.com/openimsdk/protocol/group"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/protocol/wrapperspb"
	"github.com/openimsdk/tools/mcontext"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type mockGroupDB struct {
	controller.GroupDatabase
}

func (m *mockGroupDB) TakeGroup(ctx context.Context, groupID string) (*model.Group, error) {
	return &model.Group{GroupID: groupID, GroupName: "testGroup"}, nil
}

func (m *mockGroupDB) FindGroupMemberNum(ctx context.Context, groupID string) (uint32, error) {
	return 2, nil
}

func (m *mockGroupDB) TakeGroupOwner(ctx context.Context, groupID string) (*model.GroupMember, error) {
	return &model.GroupMember{GroupID: groupID, UserID: "alice", Nickname: "alice", FaceURL: "http://example.com/face.png"}, nil
}

func (m *mockGroupDB) TakeGroupMember(ctx context.Context, groupID string, userID string) (*model.GroupMember, error) {
	return &model.GroupMember{GroupID: groupID, UserID: userID, Nickname: userID, FaceURL: "http://example.com/face.png", RoleLevel: 100}, nil
}

func (m *mockGroupDB) UpdateGroup(ctx context.Context, groupID string, data map[string]any) error {
	return nil
}

func (m *mockGroupDB) FindGroupMemberUserID(ctx context.Context, groupID string) ([]string, error) {
	return []string{"alice", "bob"}, nil
}

type mockConversationRpcClient struct {
	pbconv.ConversationClient
	setConversationsFunc func(ctx context.Context, in *pbconv.SetConversationsReq, opts ...grpc.CallOption) (*pbconv.SetConversationsResp, error)
}

func (m *mockConversationRpcClient) SetConversations(ctx context.Context, in *pbconv.SetConversationsReq, opts ...grpc.CallOption) (*pbconv.SetConversationsResp, error) {
	if m.setConversationsFunc != nil {
		return m.setConversationsFunc(ctx, in, opts...)
	}
	return &pbconv.SetConversationsResp{}, nil
}

func TestSetGroupInfo_SetConversationsError(t *testing.T) {
	expectedErr := errors.New("rpc conversation unavailable")
	mockConvRpc := &mockConversationRpcClient{
		setConversationsFunc: func(ctx context.Context, in *pbconv.SetConversationsReq, opts ...grpc.CallOption) (*pbconv.SetConversationsResp, error) {
			return nil, expectedErr
		},
	}

	gs := &groupServer{
		db:                 &mockGroupDB{},
		notification:       &NotificationSender{},
		conversationClient: &rpcli.ConversationClient{ConversationClient: mockConvRpc},
		adminUserIDs:       []string{"imAdmin"},
		config:             &Config{},
		webhookClient:      webhook.NewWebhookClient(""),
	}

	ctx := authverify.WithTempAdmin(context.Background())
	ctx = mcontext.WithOpUserIDContext(ctx, "alice")

	req := &pbgroup.SetGroupInfoReq{
		GroupInfoForSet: &sdkws.GroupInfoForSet{
			GroupID:      "test_group_1",
			Notification: "announcement text",
		},
	}

	resp, err := gs.SetGroupInfo(ctx, req)
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, expectedErr)
}

func TestSetGroupInfoEx_SetConversationsError(t *testing.T) {
	expectedErr := errors.New("rpc conversation unavailable")
	mockConvRpc := &mockConversationRpcClient{
		setConversationsFunc: func(ctx context.Context, in *pbconv.SetConversationsReq, opts ...grpc.CallOption) (*pbconv.SetConversationsResp, error) {
			return nil, expectedErr
		},
	}

	gs := &groupServer{
		db:                 &mockGroupDB{},
		notification:       &NotificationSender{},
		conversationClient: &rpcli.ConversationClient{ConversationClient: mockConvRpc},
		adminUserIDs:       []string{"imAdmin"},
		config:             &Config{},
		webhookClient:      webhook.NewWebhookClient(""),
	}

	ctx := authverify.WithTempAdmin(context.Background())
	ctx = mcontext.WithOpUserIDContext(ctx, "alice")

	req := &pbgroup.SetGroupInfoExReq{
		GroupID:      "test_group_1",
		Notification: &wrapperspb.StringValue{Value: "announcement text"},
	}

	resp, err := gs.SetGroupInfoEx(ctx, req)
	assert.Nil(t, resp)
	assert.ErrorIs(t, err, expectedErr)
}

func TestSetGroupInfo_SetConversationsAdminContext(t *testing.T) {
	var capturedOpUserID string
	mockConvRpc := &mockConversationRpcClient{
		setConversationsFunc: func(ctx context.Context, in *pbconv.SetConversationsReq, opts ...grpc.CallOption) (*pbconv.SetConversationsResp, error) {
			capturedOpUserID = mcontext.GetOpUserID(ctx)
			return nil, errors.New("check admin context")
		},
	}

	gs := &groupServer{
		db:                 &mockGroupDB{},
		notification:       &NotificationSender{},
		conversationClient: &rpcli.ConversationClient{ConversationClient: mockConvRpc},
		adminUserIDs:       []string{"imAdmin"},
		config:             &Config{},
		webhookClient:      webhook.NewWebhookClient(""),
	}

	ctx := authverify.WithTempAdmin(context.Background())
	ctx = mcontext.WithOpUserIDContext(ctx, "alice")

	req := &pbgroup.SetGroupInfoReq{
		GroupInfoForSet: &sdkws.GroupInfoForSet{
			GroupID:      "test_group_1",
			Notification: "announcement text",
		},
	}

	_, _ = gs.SetGroupInfo(ctx, req)
	assert.Equal(t, "imAdmin", capturedOpUserID, "SetConversations must be called with Admin context")
}
