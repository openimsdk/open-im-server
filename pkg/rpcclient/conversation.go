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

package rpcclient

import (
	"context"
	"fmt"

	pbconversation "github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/system/program"
	"google.golang.org/grpc"
)

type Conversation struct {
	Client pbconversation.ConversationClient
	conn   grpc.ClientConnInterface
	discov discovery.SvcDiscoveryRegistry
}

func NewConversation(discov discovery.SvcDiscoveryRegistry, rpcRegisterName string) *Conversation {
	conn, err := discov.GetConn(context.Background(), rpcRegisterName)
	if err != nil {
		program.ExitWithError(err)
	}
	client := pbconversation.NewConversationClient(conn)
	return &Conversation{discov: discov, conn: conn, Client: client}
}

type ConversationRpcClient Conversation

func NewConversationRpcClient(discov discovery.SvcDiscoveryRegistry, rpcRegisterName string) ConversationRpcClient {
	return ConversationRpcClient(*NewConversation(discov, rpcRegisterName))
}

func (c *ConversationRpcClient) GetSingleConversationRecvMsgOpt(ctx context.Context, userID, conversationID string) (int32, error) {
	var req pbconversation.GetConversationReq
	req.OwnerUserID = userID
	req.ConversationID = conversationID
	conversation, err := c.Client.GetConversation(ctx, &req)
	if err != nil {
		return 0, err
	}
	return conversation.GetConversation().RecvMsgOpt, err
}

func (c *ConversationRpcClient) SingleChatFirstCreateConversation(ctx context.Context, recvID, sendID,
	conversationID string, conversationType int32) error {
	_, err := c.Client.CreateSingleChatConversations(ctx,
		&pbconversation.CreateSingleChatConversationsReq{
			RecvID: recvID, SendID: sendID, ConversationID: conversationID,
			ConversationType: conversationType,
		})
	return err
}

func (c *ConversationRpcClient) GroupChatFirstCreateConversation(ctx context.Context, groupID string, userIDs []string) error {
	_, err := c.Client.CreateGroupChatConversations(ctx, &pbconversation.CreateGroupChatConversationsReq{UserIDs: userIDs, GroupID: groupID})
	return err
}

func (c *ConversationRpcClient) SetConversationMaxSeq(ctx context.Context, ownerUserIDs []string, conversationID string, maxSeq int64) error {
	_, err := c.Client.SetConversationMaxSeq(ctx, &pbconversation.SetConversationMaxSeqReq{OwnerUserID: ownerUserIDs, ConversationID: conversationID, MaxSeq: maxSeq})
	return err
}

func (c *ConversationRpcClient) SetConversationMinSeq(ctx context.Context, ownerUserIDs []string, conversationID string, minSeq int64) error {
	_, err := c.Client.SetConversationMinSeq(ctx, &pbconversation.SetConversationMinSeqReq{OwnerUserID: ownerUserIDs, ConversationID: conversationID, MinSeq: minSeq})
	return err
}

func (c *ConversationRpcClient) SetConversations(ctx context.Context, userIDs []string, conversation *pbconversation.ConversationReq) error {
	_, err := c.Client.SetConversations(ctx, &pbconversation.SetConversationsReq{UserIDs: userIDs, Conversation: conversation})
	return err
}

func (c *ConversationRpcClient) UpdateConversation(ctx context.Context, conversation *pbconversation.UpdateConversationReq) error {
	_, err := c.Client.UpdateConversation(ctx, conversation)
	return err
}

func (c *ConversationRpcClient) GetConversationIDs(ctx context.Context, ownerUserID string) ([]string, error) {
	resp, err := c.Client.GetConversationIDs(ctx, &pbconversation.GetConversationIDsReq{UserID: ownerUserID})
	if err != nil {
		return nil, err
	}
	return resp.ConversationIDs, nil
}

func (c *ConversationRpcClient) GetConversation(ctx context.Context, ownerUserID, conversationID string) (*pbconversation.Conversation, error) {
	resp, err := c.Client.GetConversation(ctx, &pbconversation.GetConversationReq{OwnerUserID: ownerUserID, ConversationID: conversationID})
	if err != nil {
		return nil, err
	}
	return resp.Conversation, nil
}

func (c *ConversationRpcClient) GetConversationsByConversationID(ctx context.Context, conversationIDs []string) ([]*pbconversation.Conversation, error) {
	if len(conversationIDs) == 0 {
		return nil, nil
	}
	resp, err := c.Client.GetConversationsByConversationID(ctx, &pbconversation.GetConversationsByConversationIDReq{ConversationIDs: conversationIDs})
	if err != nil {
		return nil, err
	}
	if len(resp.Conversations) == 0 {
		return nil, errs.ErrRecordNotFound.WrapMsg(fmt.Sprintf("conversationIDs: %v not found", conversationIDs))
	}
	return resp.Conversations, nil
}

func (c *ConversationRpcClient) GetConversationOfflinePushUserIDs(ctx context.Context, conversationID string, userIDs []string) ([]string, error) {
	resp, err := c.Client.GetConversationOfflinePushUserIDs(ctx, &pbconversation.GetConversationOfflinePushUserIDsReq{ConversationID: conversationID, UserIDs: userIDs})
	if err != nil {
		return nil, err
	}
	return resp.UserIDs, nil
}

func (c *ConversationRpcClient) GetConversations(ctx context.Context, ownerUserID string, conversationIDs []string) ([]*pbconversation.Conversation, error) {
	if len(conversationIDs) == 0 {
		return nil, nil
	}
	resp, err := c.Client.GetConversations(
		ctx,
		&pbconversation.GetConversationsReq{OwnerUserID: ownerUserID, ConversationIDs: conversationIDs},
	)
	if err != nil {
		return nil, err
	}
	return resp.Conversations, nil
}

func (c *ConversationRpcClient) GetConversationNotReceiveMessageUserIDs(ctx context.Context, conversationID string) ([]string, error) {
	resp, err := c.Client.GetConversationNotReceiveMessageUserIDs(ctx, &pbconversation.GetConversationNotReceiveMessageUserIDsReq{ConversationID: conversationID})
	if err != nil {
		return nil, err
	}
	return resp.UserIDs, nil
}

func (c *ConversationRpcClient) GetConversationsNeedClearMsg(ctx context.Context) ([]*pbconversation.Conversation, error) {
	resp, err := c.Client.GetConversationsNeedClearMsg(ctx, &pbconversation.GetConversationsNeedClearMsgReq{})
	if err != nil {
		return nil, err
	}
	return resp.Conversations, nil
}
