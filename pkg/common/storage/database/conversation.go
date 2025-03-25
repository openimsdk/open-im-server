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

package database

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/db/pagination"
)

type Conversation interface {
	Create(ctx context.Context, conversations []*model.Conversation) (err error)
	UpdateByMap(ctx context.Context, userIDs []string, conversationID string, args map[string]any) (rows int64, err error)
	UpdateUserConversations(ctx context.Context, userID string, args map[string]any) ([]*model.Conversation, error)
	Update(ctx context.Context, conversation *model.Conversation) (err error)
	Find(ctx context.Context, ownerUserID string, conversationIDs []string) (conversations []*model.Conversation, err error)
	FindUserID(ctx context.Context, userIDs []string, conversationIDs []string) ([]string, error)
	FindUserIDAllConversationID(ctx context.Context, userID string) ([]string, error)
	FindUserIDAllNotNotifyConversationID(ctx context.Context, userID string) ([]string, error)
	FindUserIDAllPinnedConversationID(ctx context.Context, userID string) ([]string, error)
	Take(ctx context.Context, userID, conversationID string) (conversation *model.Conversation, err error)
	FindConversationID(ctx context.Context, userID string, conversationIDs []string) (existConversationID []string, err error)
	FindUserIDAllConversations(ctx context.Context, userID string) (conversations []*model.Conversation, err error)
	FindRecvMsgUserIDs(ctx context.Context, conversationID string, recvOpts []int) ([]string, error)
	GetUserRecvMsgOpt(ctx context.Context, ownerUserID, conversationID string) (opt int, err error)
	GetAllConversationIDs(ctx context.Context) ([]string, error)
	GetAllConversationIDsNumber(ctx context.Context) (int64, error)
	PageConversationIDs(ctx context.Context, pagination pagination.Pagination) (conversationIDs []string, err error)
	GetConversationsByConversationID(ctx context.Context, conversationIDs []string) ([]*model.Conversation, error)
	GetConversationIDsNeedDestruct(ctx context.Context) ([]*model.Conversation, error)
	GetConversationNotReceiveMessageUserIDs(ctx context.Context, conversationID string) ([]string, error)
	FindConversationUserVersion(ctx context.Context, userID string, version uint, limit int) (*model.VersionLog, error)
	FindRandConversation(ctx context.Context, ts int64, limit int) ([]*model.Conversation, error)
}
