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

package controller

import (
	"context"
	"time"

	"github.com/OpenIMSDK/tools/pagination"

	"github.com/openimsdk/open-im-server/v3/pkg/msgprocessor"

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/tx"
	"github.com/OpenIMSDK/tools/utils"

	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	relationtb "github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
)

type ConversationDatabase interface {
	// UpdateUserConversationFiled 更新用户该会话的属性信息
	UpdateUsersConversationFiled(ctx context.Context, userIDs []string, conversationID string, args map[string]any) error
	// CreateConversation 创建一批新的会话
	CreateConversation(ctx context.Context, conversations []*relationtb.ConversationModel) error
	// SyncPeerUserPrivateConversation 同步对端私聊会话内部保证事务操作
	SyncPeerUserPrivateConversationTx(ctx context.Context, conversation []*relationtb.ConversationModel) error
	// FindConversations 根据会话ID获取某个用户的多个会话
	FindConversations(ctx context.Context, ownerUserID string, conversationIDs []string) ([]*relationtb.ConversationModel, error)
	// FindRecvMsgNotNotifyUserIDs 获取超级大群开启免打扰的用户ID
	//FindRecvMsgNotNotifyUserIDs(ctx context.Context, groupID string) ([]string, error)
	// GetUserAllConversation 获取一个用户在服务器上所有的会话
	GetUserAllConversation(ctx context.Context, ownerUserID string) ([]*relationtb.ConversationModel, error)
	// SetUserConversations 设置用户多个会话属性，如果会话不存在则创建，否则更新,内部保证原子性
	SetUserConversations(ctx context.Context, ownerUserID string, conversations []*relationtb.ConversationModel) error
	// SetUsersConversationFiledTx 设置多个用户会话关于某个字段的更新操作，如果会话不存在则创建，否则更新，内部保证事务操作
	SetUsersConversationFiledTx(ctx context.Context, userIDs []string, conversation *relationtb.ConversationModel, filedMap map[string]any) error
	CreateGroupChatConversation(ctx context.Context, groupID string, userIDs []string) error
	GetConversationIDs(ctx context.Context, userID string) ([]string, error)
	GetUserConversationIDsHash(ctx context.Context, ownerUserID string) (hash uint64, err error)
	GetAllConversationIDs(ctx context.Context) ([]string, error)
	GetAllConversationIDsNumber(ctx context.Context) (int64, error)
	PageConversationIDs(ctx context.Context, pagination pagination.Pagination) (conversationIDs []string, err error)
	//GetUserAllHasReadSeqs(ctx context.Context, ownerUserID string) (map[string]int64, error)
	GetConversationsByConversationID(ctx context.Context, conversationIDs []string) ([]*relationtb.ConversationModel, error)
	GetConversationIDsNeedDestruct(ctx context.Context) ([]*relationtb.ConversationModel, error)
	GetConversationNotReceiveMessageUserIDs(ctx context.Context, conversationID string) ([]string, error)
}

func NewConversationDatabase(conversation relationtb.ConversationModelInterface, cache cache.ConversationCache, tx tx.CtxTx) ConversationDatabase {
	return &conversationDatabase{
		conversationDB: conversation,
		cache:          cache,
		tx:             tx,
	}
}

type conversationDatabase struct {
	conversationDB relationtb.ConversationModelInterface
	cache          cache.ConversationCache
	tx             tx.CtxTx
}

func (c *conversationDatabase) SetUsersConversationFiledTx(ctx context.Context, userIDs []string, conversation *relationtb.ConversationModel, filedMap map[string]any) (err error) {
	return c.tx.Transaction(ctx, func(ctx context.Context) error {
		cache := c.cache.NewCache()
		if conversation.GroupID != "" {
			cache = cache.DelSuperGroupRecvMsgNotNotifyUserIDs(conversation.GroupID).DelSuperGroupRecvMsgNotNotifyUserIDsHash(conversation.GroupID)
		}
		haveUserIDs, err := c.conversationDB.FindUserID(ctx, userIDs, []string{conversation.ConversationID})
		if err != nil {
			return err
		}
		if len(haveUserIDs) > 0 {
			_, err = c.conversationDB.UpdateByMap(ctx, haveUserIDs, conversation.ConversationID, filedMap)
			if err != nil {
				return err
			}
			cache = cache.DelUsersConversation(conversation.ConversationID, haveUserIDs...)
			if _, ok := filedMap["has_read_seq"]; ok {
				for _, userID := range haveUserIDs {
					cache = cache.DelUserAllHasReadSeqs(userID, conversation.ConversationID)
				}
			}
			if _, ok := filedMap["recv_msg_opt"]; ok {
				cache = cache.DelConversationNotReceiveMessageUserIDs(conversation.ConversationID)
			}
		}
		NotUserIDs := utils.DifferenceString(haveUserIDs, userIDs)
		log.ZDebug(ctx, "SetUsersConversationFiledTx", "NotUserIDs", NotUserIDs, "haveUserIDs", haveUserIDs, "userIDs", userIDs)
		var conversations []*relationtb.ConversationModel
		now := time.Now()
		for _, v := range NotUserIDs {
			temp := new(relationtb.ConversationModel)
			if err := utils.CopyStructFields(temp, conversation); err != nil {
				return err
			}
			temp.OwnerUserID = v
			temp.CreateTime = now
			conversations = append(conversations, temp)
		}
		if len(conversations) > 0 {
			err = c.conversationDB.Create(ctx, conversations)
			if err != nil {
				return err
			}
			cache = cache.DelConversationIDs(NotUserIDs...).DelUserConversationIDsHash(NotUserIDs...).DelConversations(conversation.ConversationID, NotUserIDs...)
		}
		return cache.ExecDel(ctx)
	})
}

func (c *conversationDatabase) UpdateUsersConversationFiled(ctx context.Context, userIDs []string, conversationID string, args map[string]any) error {
	_, err := c.conversationDB.UpdateByMap(ctx, userIDs, conversationID, args)
	if err != nil {
		return err
	}
	cache := c.cache.NewCache()
	cache = cache.DelUsersConversation(conversationID, userIDs...)
	if _, ok := args["recv_msg_opt"]; ok {
		cache = cache.DelConversationNotReceiveMessageUserIDs(conversationID)
	}
	return cache.ExecDel(ctx)
}

func (c *conversationDatabase) CreateConversation(ctx context.Context, conversations []*relationtb.ConversationModel) error {
	if err := c.conversationDB.Create(ctx, conversations); err != nil {
		return err
	}
	var userIDs []string
	cache := c.cache.NewCache()
	for _, conversation := range conversations {
		cache = cache.DelConversations(conversation.OwnerUserID, conversation.ConversationID)
		cache = cache.DelConversationNotReceiveMessageUserIDs(conversation.ConversationID)
		userIDs = append(userIDs, conversation.OwnerUserID)
	}
	return cache.DelConversationIDs(userIDs...).DelUserConversationIDsHash(userIDs...).ExecDel(ctx)
}

func (c *conversationDatabase) SyncPeerUserPrivateConversationTx(ctx context.Context, conversations []*relationtb.ConversationModel) error {
	return c.tx.Transaction(ctx, func(ctx context.Context) error {
		cache := c.cache.NewCache()
		for _, conversation := range conversations {
			for _, v := range [][2]string{{conversation.OwnerUserID, conversation.UserID}, {conversation.UserID, conversation.OwnerUserID}} {
				ownerUserID := v[0]
				userID := v[1]
				haveUserIDs, err := c.conversationDB.FindUserID(ctx, []string{ownerUserID}, []string{conversation.ConversationID})
				if err != nil {
					return err
				}
				if len(haveUserIDs) > 0 {
					_, err := c.conversationDB.UpdateByMap(ctx, []string{ownerUserID}, conversation.ConversationID, map[string]any{"is_private_chat": conversation.IsPrivateChat})
					if err != nil {
						return err
					}
					cache = cache.DelUsersConversation(conversation.ConversationID, ownerUserID)
				} else {
					newConversation := *conversation
					newConversation.OwnerUserID = ownerUserID
					newConversation.UserID = userID
					newConversation.ConversationID = conversation.ConversationID
					newConversation.IsPrivateChat = conversation.IsPrivateChat
					if err := c.conversationDB.Create(ctx, []*relationtb.ConversationModel{&newConversation}); err != nil {
						return err
					}
					cache = cache.DelConversationIDs(ownerUserID).DelUserConversationIDsHash(ownerUserID)
				}
			}
		}
		return cache.ExecDel(ctx)
	})
}

func (c *conversationDatabase) FindConversations(ctx context.Context, ownerUserID string, conversationIDs []string) ([]*relationtb.ConversationModel, error) {
	return c.cache.GetConversations(ctx, ownerUserID, conversationIDs)
}

func (c *conversationDatabase) GetConversation(ctx context.Context, ownerUserID string, conversationID string) (*relationtb.ConversationModel, error) {
	return c.cache.GetConversation(ctx, ownerUserID, conversationID)
}

func (c *conversationDatabase) GetUserAllConversation(ctx context.Context, ownerUserID string) ([]*relationtb.ConversationModel, error) {
	return c.cache.GetUserAllConversations(ctx, ownerUserID)
}

func (c *conversationDatabase) SetUserConversations(ctx context.Context, ownerUserID string, conversations []*relationtb.ConversationModel) error {
	return c.tx.Transaction(ctx, func(ctx context.Context) error {
		cache := c.cache.NewCache()
		groupIDs := utils.Distinct(utils.Filter(conversations, func(e *relationtb.ConversationModel) (string, bool) {
			return e.GroupID, e.GroupID != ""
		}))
		for _, groupID := range groupIDs {
			cache = cache.DelSuperGroupRecvMsgNotNotifyUserIDs(groupID).DelSuperGroupRecvMsgNotNotifyUserIDsHash(groupID)
		}
		var conversationIDs []string
		for _, conversation := range conversations {
			conversationIDs = append(conversationIDs, conversation.ConversationID)
			cache = cache.DelConversations(conversation.OwnerUserID, conversation.ConversationID)
		}
		existConversations, err := c.conversationDB.Find(ctx, ownerUserID, conversationIDs)
		if err != nil {
			return err
		}
		if len(existConversations) > 0 {
			for _, conversation := range conversations {
				err = c.conversationDB.Update(ctx, conversation)
				if err != nil {
					return err
				}
			}
		}
		var existConversationIDs []string
		for _, conversation := range existConversations {
			existConversationIDs = append(existConversationIDs, conversation.ConversationID)
		}

		var notExistConversations []*relationtb.ConversationModel
		for _, conversation := range conversations {
			if !utils.IsContain(conversation.ConversationID, existConversationIDs) {
				notExistConversations = append(notExistConversations, conversation)
			}
		}
		if len(notExistConversations) > 0 {
			err = c.conversationDB.Create(ctx, notExistConversations)
			if err != nil {
				return err
			}
			cache = cache.DelConversationIDs(ownerUserID).
				DelUserConversationIDsHash(ownerUserID).
				DelConversationNotReceiveMessageUserIDs(utils.Slice(notExistConversations, func(e *relationtb.ConversationModel) string { return e.ConversationID })...)
		}
		return cache.ExecDel(ctx)
	})
}

//func (c *conversationDatabase) FindRecvMsgNotNotifyUserIDs(ctx context.Context, groupID string) ([]string, error) {
//	return c.cache.GetSuperGroupRecvMsgNotNotifyUserIDs(ctx, groupID)
//}

func (c *conversationDatabase) CreateGroupChatConversation(ctx context.Context, groupID string, userIDs []string) error {
	return c.tx.Transaction(ctx, func(ctx context.Context) error {
		cache := c.cache.NewCache()
		conversationID := msgprocessor.GetConversationIDBySessionType(constant.SuperGroupChatType, groupID)
		existConversationUserIDs, err := c.conversationDB.FindUserID(ctx, userIDs, []string{conversationID})
		if err != nil {
			return err
		}
		notExistUserIDs := utils.DifferenceString(userIDs, existConversationUserIDs)
		var conversations []*relationtb.ConversationModel
		for _, v := range notExistUserIDs {
			conversation := relationtb.ConversationModel{ConversationType: constant.SuperGroupChatType, GroupID: groupID, OwnerUserID: v, ConversationID: conversationID}
			conversations = append(conversations, &conversation)
			cache = cache.DelConversations(v, conversationID).DelConversationNotReceiveMessageUserIDs(conversationID)
		}
		cache = cache.DelConversationIDs(notExistUserIDs...).DelUserConversationIDsHash(notExistUserIDs...)
		if len(conversations) > 0 {
			err = c.conversationDB.Create(ctx, conversations)
			if err != nil {
				return err
			}
		}
		_, err = c.conversationDB.UpdateByMap(ctx, existConversationUserIDs, conversationID, map[string]any{"max_seq": 0})
		if err != nil {
			return err
		}
		for _, v := range existConversationUserIDs {
			cache = cache.DelConversations(v, conversationID)
		}
		return cache.ExecDel(ctx)
	})
}

func (c *conversationDatabase) GetConversationIDs(ctx context.Context, userID string) ([]string, error) {
	return c.cache.GetUserConversationIDs(ctx, userID)
}

func (c *conversationDatabase) GetUserConversationIDsHash(ctx context.Context, ownerUserID string) (hash uint64, err error) {
	return c.cache.GetUserConversationIDsHash(ctx, ownerUserID)
}

func (c *conversationDatabase) GetAllConversationIDs(ctx context.Context) ([]string, error) {
	return c.conversationDB.GetAllConversationIDs(ctx)
}

func (c *conversationDatabase) GetAllConversationIDsNumber(ctx context.Context) (int64, error) {
	return c.conversationDB.GetAllConversationIDsNumber(ctx)
}

func (c *conversationDatabase) PageConversationIDs(ctx context.Context, pagination pagination.Pagination) ([]string, error) {
	return c.conversationDB.PageConversationIDs(ctx, pagination)
}

func (c *conversationDatabase) GetConversationsByConversationID(ctx context.Context, conversationIDs []string) ([]*relationtb.ConversationModel, error) {
	return c.conversationDB.GetConversationsByConversationID(ctx, conversationIDs)
}

func (c *conversationDatabase) GetConversationIDsNeedDestruct(ctx context.Context) ([]*relationtb.ConversationModel, error) {
	return c.conversationDB.GetConversationIDsNeedDestruct(ctx)
}

func (c *conversationDatabase) GetConversationNotReceiveMessageUserIDs(ctx context.Context, conversationID string) ([]string, error) {
	return c.cache.GetConversationNotReceiveMessageUserIDs(ctx, conversationID)
}
