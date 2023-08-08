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

	"google.golang.org/grpc"

	pbConversation "github.com/OpenIMSDK/protocol/conversation"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/OpenIMSDK/tools/errs"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
)

type Conversation struct {
	Client pbConversation.ConversationClient
	conn   grpc.ClientConnInterface
	discov discoveryregistry.SvcDiscoveryRegistry
}

func NewConversation(discov discoveryregistry.SvcDiscoveryRegistry) *Conversation {
	conn, err := discov.GetConn(context.Background(), config.Config.RpcRegisterName.OpenImConversationName)
	if err != nil {
		panic(err)
	}
	client := pbConversation.NewConversationClient(conn)
	return &Conversation{discov: discov, conn: conn, Client: client}
}

type ConversationRpcClient Conversation

func NewConversationRpcClient(discov discoveryregistry.SvcDiscoveryRegistry) ConversationRpcClient {
	return ConversationRpcClient(*NewConversation(discov))
}

func (c *ConversationRpcClient) GetSingleConversationRecvMsgOpt(ctx context.Context, userID, conversationID string) (int32, error) {
	var req pbConversation.GetConversationReq
	req.OwnerUserID = userID
	req.ConversationID = conversationID
	conversation, err := c.Client.GetConversation(ctx, &req)
	if err != nil {
		return 0, err
	}
	return conversation.GetConversation().RecvMsgOpt, err
}

func (c *ConversationRpcClient) SingleChatFirstCreateConversation(ctx context.Context, recvID, sendID string) error {
	_, err := c.Client.CreateSingleChatConversations(ctx, &pbConversation.CreateSingleChatConversationsReq{RecvID: recvID, SendID: sendID})
	return err
}

func (c *ConversationRpcClient) GroupChatFirstCreateConversation(ctx context.Context, groupID string, userIDs []string) error {
	_, err := c.Client.CreateGroupChatConversations(ctx, &pbConversation.CreateGroupChatConversationsReq{UserIDs: userIDs, GroupID: groupID})
	return err
}

func (c *ConversationRpcClient) SetConversationMaxSeq(ctx context.Context, ownerUserIDs []string, conversationID string, maxSeq int64) error {
	_, err := c.Client.SetConversationMaxSeq(ctx, &pbConversation.SetConversationMaxSeqReq{OwnerUserID: ownerUserIDs, ConversationID: conversationID, MaxSeq: maxSeq})
	return err
}

func (c *ConversationRpcClient) SetConversations(ctx context.Context, userIDs []string, conversation *pbConversation.ConversationReq) error {
	_, err := c.Client.SetConversations(ctx, &pbConversation.SetConversationsReq{UserIDs: userIDs, Conversation: conversation})
	return err
}

func (c *ConversationRpcClient) GetConversationIDs(ctx context.Context, ownerUserID string) ([]string, error) {
	resp, err := c.Client.GetConversationIDs(ctx, &pbConversation.GetConversationIDsReq{UserID: ownerUserID})
	if err != nil {
		return nil, err
	}
	return resp.ConversationIDs, nil
}

func (c *ConversationRpcClient) GetConversation(ctx context.Context, ownerUserID, conversationID string) (*pbConversation.Conversation, error) {
	resp, err := c.Client.GetConversation(ctx, &pbConversation.GetConversationReq{OwnerUserID: ownerUserID, ConversationID: conversationID})
	if err != nil {
		return nil, err
	}
	return resp.Conversation, nil
}

func (c *ConversationRpcClient) GetConversationsByConversationID(ctx context.Context, conversationIDs []string) ([]*pbConversation.Conversation, error) {
	if len(conversationIDs) == 0 {
		return nil, nil
	}
	resp, err := c.Client.GetConversationsByConversationID(ctx, &pbConversation.GetConversationsByConversationIDReq{ConversationIDs: conversationIDs})
	if err != nil {
		return nil, err
	}
	if len(resp.Conversations) == 0 {
		return nil, errs.ErrRecordNotFound.Wrap(fmt.Sprintf("conversationIDs: %v not found", conversationIDs))
	}
	return resp.Conversations, nil
}

func (c *ConversationRpcClient) GetConversations(
	ctx context.Context,
	ownerUserID string,
	conversationIDs []string,
) ([]*pbConversation.Conversation, error) {
	if len(conversationIDs) == 0 {
		return nil, nil
	}
	resp, err := c.Client.GetConversations(
		ctx,
		&pbConversation.GetConversationsReq{OwnerUserID: ownerUserID, ConversationIDs: conversationIDs},
	)
	if err != nil {
		return nil, err
	}
	return resp.Conversations, nil
}
