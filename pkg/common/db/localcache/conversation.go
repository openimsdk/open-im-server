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

package localcache

import (
	"context"
	"sync"

	"github.com/OpenIMSDK/protocol/conversation"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
)

type ConversationLocalCache struct {
	lock                              sync.Mutex
	superGroupRecvMsgNotNotifyUserIDs map[string]Hash
	conversationIDs                   map[string]Hash
	client                            *rpcclient.ConversationRpcClient
}

type Hash struct {
	hash uint64
	ids  []string
}

func NewConversationLocalCache(client *rpcclient.ConversationRpcClient) *ConversationLocalCache {
	return &ConversationLocalCache{
		superGroupRecvMsgNotNotifyUserIDs: make(map[string]Hash),
		conversationIDs:                   make(map[string]Hash),
		client:                            client,
	}
}

func (g *ConversationLocalCache) GetRecvMsgNotNotifyUserIDs(ctx context.Context, groupID string) ([]string, error) {
	resp, err := g.client.Client.GetRecvMsgNotNotifyUserIDs(ctx, &conversation.GetRecvMsgNotNotifyUserIDsReq{
		GroupID: groupID,
	})
	if err != nil {
		return nil, err
	}
	return resp.UserIDs, nil
}

func (g *ConversationLocalCache) GetConversationIDs(ctx context.Context, userID string) ([]string, error) {
	resp, err := g.client.Client.GetUserConversationIDsHash(ctx, &conversation.GetUserConversationIDsHashReq{
		OwnerUserID: userID,
	})
	if err != nil {
		return nil, err
	}

	g.lock.Lock()
	hash, ok := g.conversationIDs[userID]
	g.lock.Unlock()

	if !ok || hash.hash != resp.Hash {
		conversationIDsResp, err := g.client.Client.GetConversationIDs(ctx, &conversation.GetConversationIDsReq{
			UserID: userID,
		})
		if err != nil {
			return nil, err
		}

		g.lock.Lock()
		defer g.lock.Unlock()
		g.conversationIDs[userID] = Hash{
			hash: resp.Hash,
			ids:  conversationIDsResp.ConversationIDs,
		}

		return conversationIDsResp.ConversationIDs, nil
	}

	return hash.ids, nil
}
