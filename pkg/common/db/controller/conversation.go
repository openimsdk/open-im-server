package controller

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	relationTb "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/tx"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
)

type ConversationDatabase interface {
	//UpdateUserConversationFiled 更新用户该会话的属性信息
	UpdateUsersConversationFiled(ctx context.Context, userIDs []string, conversationID string, args map[string]interface{}) error
	//CreateConversation 创建一批新的会话
	CreateConversation(ctx context.Context, conversations []*relationTb.ConversationModel) error
	//SyncPeerUserPrivateConversation 同步对端私聊会话内部保证事务操作
	SyncPeerUserPrivateConversationTx(ctx context.Context, conversation []*relationTb.ConversationModel) error
	//FindConversations 根据会话ID获取某个用户的多个会话
	FindConversations(ctx context.Context, ownerUserID string, conversationIDs []string) ([]*relationTb.ConversationModel, error)
	//FindRecvMsgNotNotifyUserIDs 获取超级大群开启免打扰的用户ID
	FindRecvMsgNotNotifyUserIDs(ctx context.Context, groupID string) ([]string, error)
	//GetUserAllConversation 获取一个用户在服务器上所有的会话
	GetUserAllConversation(ctx context.Context, ownerUserID string) ([]*relationTb.ConversationModel, error)
	//SetUserConversations 设置用户多个会话属性，如果会话不存在则创建，否则更新,内部保证原子性
	SetUserConversations(ctx context.Context, ownerUserID string, conversations []*relationTb.ConversationModel) error
	//SetUsersConversationFiledTx 设置多个用户会话关于某个字段的更新操作，如果会话不存在则创建，否则更新，内部保证事务操作
	SetUsersConversationFiledTx(ctx context.Context, userIDs []string, conversation *relationTb.ConversationModel, filedMap map[string]interface{}) error
	CreateGroupChatConversation(ctx context.Context, groupID string, userIDs []string) error
	GetConversationIDs(ctx context.Context, userID string) ([]string, error)
	GetUserConversationIDsHash(ctx context.Context, ownerUserID string) (hash uint64, err error)
	GetAllConversationIDs(ctx context.Context) ([]string, error)
	GetUserAllHasReadSeqs(ctx context.Context, ownerUserID string) (map[string]int64, error)
	GetConversationsByConversationID(ctx context.Context, conversationIDs []string) ([]*relationTb.ConversationModel, error)
}

func NewConversationDatabase(conversation relationTb.ConversationModelInterface, cache cache.ConversationCache, tx tx.Tx) ConversationDatabase {
	return &conversationDatabase{
		conversationDB: conversation,
		cache:          cache,
		tx:             tx,
	}
}

type conversationDatabase struct {
	conversationDB relationTb.ConversationModelInterface
	cache          cache.ConversationCache
	tx             tx.Tx
}

func (c *conversationDatabase) SetUsersConversationFiledTx(ctx context.Context, userIDs []string, conversation *relationTb.ConversationModel, filedMap map[string]interface{}) (err error) {
	cache := c.cache.NewCache()
	if err := c.tx.Transaction(func(tx any) error {
		conversationTx := c.conversationDB.NewTx(tx)
		haveUserIDs, err := conversationTx.FindUserID(ctx, userIDs, []string{conversation.ConversationID})
		if err != nil {
			return err
		}
		if len(haveUserIDs) > 0 {
			_, err = conversationTx.UpdateByMap(ctx, haveUserIDs, conversation.ConversationID, filedMap)
			if err != nil {
				return err
			}
			cache = cache.DelUsersConversation(conversation.ConversationID, haveUserIDs...)
			if _, ok := filedMap["has_read_seq"]; ok {
				for _, userID := range haveUserIDs {
					cache = cache.DelUserAllHasReadSeqs(userID, conversation.ConversationID)
				}
			}
		}
		NotUserIDs := utils.DifferenceString(haveUserIDs, userIDs)
		log.ZDebug(ctx, "SetUsersConversationFiledTx", "NotUserIDs", NotUserIDs, "haveUserIDs", haveUserIDs, "userIDs", userIDs)
		var conversations []*relationTb.ConversationModel
		for _, v := range NotUserIDs {
			temp := new(relationTb.ConversationModel)
			if err := utils.CopyStructFields(temp, conversation); err != nil {
				return err
			}
			temp.OwnerUserID = v
			conversations = append(conversations, temp)

		}
		if len(conversations) > 0 {
			err = conversationTx.Create(ctx, conversations)
			if err != nil {
				return err
			}
			cache = cache.DelConversationIDs(NotUserIDs...).DelUserConversationIDsHash(NotUserIDs...)
		}
		return nil
	}); err != nil {
		return err
	}
	return cache.ExecDel(ctx)
}

func (c *conversationDatabase) UpdateUsersConversationFiled(ctx context.Context, userIDs []string, conversationID string, args map[string]interface{}) error {
	_, err := c.conversationDB.UpdateByMap(ctx, userIDs, conversationID, args)
	if err != nil {
		return err
	}
	return c.cache.DelUsersConversation(conversationID, userIDs...).ExecDel(ctx)
}

func (c *conversationDatabase) CreateConversation(ctx context.Context, conversations []*relationTb.ConversationModel) error {
	if err := c.conversationDB.Create(ctx, conversations); err != nil {
		return err
	}
	var userIDs []string
	cache := c.cache.NewCache()
	for _, conversation := range conversations {
		cache = cache.DelConvsersations(conversation.OwnerUserID, conversation.ConversationID)
		userIDs = append(userIDs, conversation.OwnerUserID)
	}
	return cache.DelConversationIDs(userIDs...).DelUserConversationIDsHash(userIDs...).ExecDel(ctx)
}

func (c *conversationDatabase) SyncPeerUserPrivateConversationTx(ctx context.Context, conversations []*relationTb.ConversationModel) error {
	cache := c.cache.NewCache()
	if err := c.tx.Transaction(func(tx any) error {
		conversationTx := c.conversationDB.NewTx(tx)
		for _, conversation := range conversations {
			for _, v := range [][2]string{{conversation.OwnerUserID, conversation.UserID}, {conversation.UserID, conversation.OwnerUserID}} {
				haveUserIDs, err := conversationTx.FindUserID(ctx, []string{v[0]}, []string{conversation.ConversationID})
				if err != nil {
					return err
				}
				if len(haveUserIDs) > 0 {
					_, err := conversationTx.UpdateByMap(ctx, []string{v[0]}, conversation.ConversationID, map[string]interface{}{"is_private_chat": conversation.IsPrivateChat})
					if err != nil {
						return err
					}
					cache = cache.DelUsersConversation(conversation.ConversationID, v[0])
				} else {
					newConversation := *conversation
					newConversation.OwnerUserID = v[0]
					newConversation.UserID = v[1]
					newConversation.ConversationID = conversation.ConversationID
					newConversation.IsPrivateChat = conversation.IsPrivateChat
					if err := conversationTx.Create(ctx, []*relationTb.ConversationModel{&newConversation}); err != nil {
						return err
					}
					cache = cache.DelConversationIDs(v[0]).DelUserConversationIDsHash(v[0])
				}
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return c.cache.ExecDel(ctx)
}

func (c *conversationDatabase) FindConversations(ctx context.Context, ownerUserID string, conversationIDs []string) ([]*relationTb.ConversationModel, error) {
	return c.cache.GetConversations(ctx, ownerUserID, conversationIDs)
}

func (c *conversationDatabase) GetConversation(ctx context.Context, ownerUserID string, conversationID string) (*relationTb.ConversationModel, error) {
	return c.cache.GetConversation(ctx, ownerUserID, conversationID)
}

func (c *conversationDatabase) GetUserAllConversation(ctx context.Context, ownerUserID string) ([]*relationTb.ConversationModel, error) {
	return c.cache.GetUserAllConversations(ctx, ownerUserID)
}

func (c *conversationDatabase) SetUserConversations(ctx context.Context, ownerUserID string, conversations []*relationTb.ConversationModel) error {
	cache := c.cache.NewCache()
	if err := c.tx.Transaction(func(tx any) error {
		var conversationIDs []string
		for _, conversation := range conversations {
			conversationIDs = append(conversationIDs, conversation.ConversationID)
		}
		conversationTx := c.conversationDB.NewTx(tx)
		existConversations, err := conversationTx.Find(ctx, ownerUserID, conversationIDs)
		if err != nil {
			return err
		}
		if len(existConversations) > 0 {
			for _, conversation := range conversations {
				err = conversationTx.Update(ctx, conversation)
				if err != nil {
					return err
				}
			}
		}
		var existConversationIDs []string
		for _, conversation := range existConversations {
			existConversationIDs = append(existConversationIDs, conversation.ConversationID)
		}

		var notExistConversations []*relationTb.ConversationModel
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
			cache = cache.DelConversationIDs(ownerUserID).DelUserConversationIDsHash(ownerUserID)
		}
		cache = cache.DelConvsersations(ownerUserID, existConversationIDs...)
		return nil
	}); err != nil {
		return err
	}
	return cache.ExecDel(ctx)
}

func (c *conversationDatabase) FindRecvMsgNotNotifyUserIDs(ctx context.Context, groupID string) ([]string, error) {
	return c.cache.GetSuperGroupRecvMsgNotNotifyUserIDs(ctx, groupID)
}

func (c *conversationDatabase) CreateGroupChatConversation(ctx context.Context, groupID string, userIDs []string) error {
	cache := c.cache.NewCache()
	conversationID := utils.GetConversationIDBySessionType(constant.SuperGroupChatType, groupID)
	if err := c.tx.Transaction(func(tx any) error {
		existConversationUserIDs, err := c.conversationDB.FindUserID(ctx, userIDs, []string{conversationID})
		if err != nil {
			return err
		}
		notExistUserIDs := utils.DifferenceString(userIDs, existConversationUserIDs)
		var conversations []*relationTb.ConversationModel
		for _, v := range notExistUserIDs {
			conversation := relationTb.ConversationModel{ConversationType: constant.SuperGroupChatType, GroupID: groupID, OwnerUserID: v, ConversationID: conversationID}
			conversations = append(conversations, &conversation)
		}
		cache = cache.DelConversationIDs(notExistUserIDs...).DelUserConversationIDsHash(notExistUserIDs...)
		if len(conversations) > 0 {
			err = c.conversationDB.Create(ctx, conversations)
			if err != nil {
				return err
			}
		}
		_, err = c.conversationDB.UpdateByMap(ctx, existConversationUserIDs, conversationID, map[string]interface{}{"max_seq": 0})
		if err != nil {
			return err
		}
		for _, v := range existConversationUserIDs {
			cache = cache.DelConvsersations(v, conversationID)
		}
		return nil
	}); err != nil {
		return err
	}
	return cache.ExecDel(ctx)
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

func (c *conversationDatabase) GetUserAllHasReadSeqs(ctx context.Context, ownerUserID string) (map[string]int64, error) {
	return c.cache.GetUserAllHasReadSeqs(ctx, ownerUserID)
}

func (c *conversationDatabase) GetConversationsByConversationID(ctx context.Context, conversationIDs []string) ([]*relationTb.ConversationModel, error) {
	return c.conversationDB.GetConversationsByConversationID(ctx, conversationIDs)
}
