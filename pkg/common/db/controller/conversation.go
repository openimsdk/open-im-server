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
	SyncPeerUserPrivateConversationTx(ctx context.Context, conversation *relationTb.ConversationModel) error
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
}

func NewConversationDatabase(conversation relationTb.ConversationModelInterface, cache cache.ConversationCache, tx tx.Tx) ConversationDatabase {
	return &ConversationDataBase{
		conversationDB: conversation,
		cache:          cache,
		tx:             tx,
	}
}

type ConversationDataBase struct {
	conversationDB relationTb.ConversationModelInterface
	cache          cache.ConversationCache
	tx             tx.Tx
}

func (c *ConversationDataBase) SetUsersConversationFiledTx(ctx context.Context, userIDs []string, conversation *relationTb.ConversationModel, filedMap map[string]interface{}) error {
	return c.tx.Transaction(func(tx any) error {
		conversationTx := c.conversationDB.NewTx(tx)
		haveUserIDs, err := conversationTx.FindUserID(ctx, userIDs, []string{conversation.ConversationID})
		if err != nil {
			return err
		}
		cache := c.cache.NewCache()
		if len(haveUserIDs) > 0 {
			_, err = conversationTx.UpdateByMap(ctx, haveUserIDs, conversation.ConversationID, filedMap)
			if err != nil {
				return err
			}
			cache = cache.DelUsersConversation(conversation.ConversationID, haveUserIDs...)
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
			cache = cache.DelConversationIDs(NotUserIDs)
		}
		// clear cache
		log.ZDebug(ctx, "SetUsersConversationFiledTx", "cache", cache.GetPreDelKeys(), "addr", &cache)
		return cache.ExecDel(ctx)
	})
}

func (c *ConversationDataBase) UpdateUsersConversationFiled(ctx context.Context, userIDs []string, conversationID string, args map[string]interface{}) error {
	_, err := c.conversationDB.UpdateByMap(ctx, userIDs, conversationID, args)
	if err != nil {
		return err
	}
	return c.cache.DelUsersConversation(conversationID, userIDs...).ExecDel(ctx)
}

func (c *ConversationDataBase) CreateConversation(ctx context.Context, conversations []*relationTb.ConversationModel) error {
	return c.tx.Transaction(func(tx any) error {
		if err := c.conversationDB.NewTx(tx).Create(ctx, conversations); err != nil {
			return err
		}
		return nil
	})
}

func (c *ConversationDataBase) SyncPeerUserPrivateConversationTx(ctx context.Context, conversation *relationTb.ConversationModel) error {
	return c.tx.Transaction(func(tx any) error {
		conversationTx := c.conversationDB.NewTx(tx)
		cache := c.cache.NewCache()
		for _, v := range [][3]string{{conversation.OwnerUserID, conversation.ConversationID, conversation.UserID}, {conversation.UserID, utils.GetConversationIDBySessionType(conversation.OwnerUserID, constant.SingleChatType), conversation.OwnerUserID}} {
			haveUserIDs, err := conversationTx.FindUserID(ctx, []string{v[0]}, []string{v[1]})
			if err != nil {
				return err
			}
			if len(haveUserIDs) > 0 {
				_, err := conversationTx.UpdateByMap(ctx, []string{v[0]}, v[1], map[string]interface{}{"is_private_chat": conversation.IsPrivateChat})
				if err != nil {
					return err
				}
				cache = cache.DelUsersConversation(v[1], v[0])
			} else {
				newConversation := *conversation
				newConversation.OwnerUserID = v[0]
				newConversation.UserID = v[2]
				newConversation.ConversationID = v[1]
				newConversation.IsPrivateChat = conversation.IsPrivateChat
				if err := conversationTx.Create(ctx, []*relationTb.ConversationModel{&newConversation}); err != nil {
					return err
				}
				cache = cache.DelConversationIDs([]string{v[0]})
			}
		}
		return c.cache.ExecDel(ctx)
	})
}

func (c *ConversationDataBase) FindConversations(ctx context.Context, ownerUserID string, conversationIDs []string) ([]*relationTb.ConversationModel, error) {
	return c.cache.GetConversations(ctx, ownerUserID, conversationIDs)
}

func (c *ConversationDataBase) GetConversation(ctx context.Context, ownerUserID string, conversationID string) (*relationTb.ConversationModel, error) {
	return c.cache.GetConversation(ctx, ownerUserID, conversationID)
}

func (c *ConversationDataBase) GetUserAllConversation(ctx context.Context, ownerUserID string) ([]*relationTb.ConversationModel, error) {
	return c.cache.GetUserAllConversations(ctx, ownerUserID)
}

func (c *ConversationDataBase) SetUserConversations(ctx context.Context, ownerUserID string, conversations []*relationTb.ConversationModel) error {
	return c.tx.Transaction(func(tx any) error {
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
		}
		cache := c.cache.NewCache()
		if len(notExistConversations) > 0 {
			cache = cache.DelConversationIDs([]string{ownerUserID})
		}
		return cache.DelConvsersations(ownerUserID, existConversationIDs).ExecDel(ctx)
	})
}

func (c *ConversationDataBase) FindRecvMsgNotNotifyUserIDs(ctx context.Context, groupID string) ([]string, error) {
	return c.cache.GetSuperGroupRecvMsgNotNotifyUserIDs(ctx, groupID)
}
