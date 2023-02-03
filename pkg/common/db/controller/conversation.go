package controller

import (
	"Open_IM/pkg/common/db/cache"
	"Open_IM/pkg/common/db/relation"
	"Open_IM/pkg/common/db/table"
	"context"
)

type ConversationInterface interface {
	//GetUserIDExistConversation 获取拥有该会话的的用户ID列表
	GetUserIDExistConversation(ctx context.Context, userIDList []string, conversationID string) ([]string, error)
	//UpdateUserConversationFiled 更新用户该会话的属性信息
	UpdateUsersConversationFiled(ctx context.Context, UserIDList []string, conversationID string, args map[string]interface{}) error
	//CreateConversation 创建一批新的会话
	CreateConversation(ctx context.Context, conversations []*table.ConversationModel) error
	//SyncPeerUserPrivateConversation 同步对端私聊会话内部保证事务操作
	SyncPeerUserPrivateConversationTx(ctx context.Context, conversation *table.ConversationModel) error
	//FindConversations 根据会话ID获取某个用户的多个会话
	FindConversations(ctx context.Context, ownerUserID string, conversationID []string) ([]*table.ConversationModel, error)
	//GetUserAllConversation 获取一个用户在服务器上所有的会话
	GetUserAllConversation(ctx context.Context, ownerUserID string) ([]*table.ConversationModel, error)
	//SetUserConversations 设置用户多个会话属性，如果会话不存在则创建，否则更新,内部保证原子性
	SetUserConversations(ctx context.Context, ownerUserID string, conversations []*table.ConversationModel) error
}
type ConversationController struct {
	database ConversationDataBaseInterface
}

func NewConversationController(database ConversationDataBaseInterface) *ConversationController {
	return &ConversationController{database: database}
}

func (c *ConversationController) GetUserIDExistConversation(ctx context.Context, userIDList []string, conversationID string) ([]string, error) {
	return c.database.GetUserIDExistConversation(ctx, userIDList, conversationID)
}

func (c ConversationController) UpdateUsersConversationFiled(ctx context.Context, UserIDList []string, conversationID string, args map[string]interface{}) error {
	panic("implement me")
}

func (c ConversationController) CreateConversation(ctx context.Context, conversations []*table.ConversationModel) error {
	panic("implement me")
}

func (c ConversationController) SyncPeerUserPrivateConversationTx(ctx context.Context, conversation *table.ConversationModel) error {
	panic("implement me")
}

func (c ConversationController) FindConversations(ctx context.Context, ownerUserID string, conversationID []string) ([]*table.ConversationModel, error) {
	panic("implement me")
}

func (c ConversationController) GetUserAllConversation(ctx context.Context, ownerUserID string) ([]*table.ConversationModel, error) {
	panic("implement me")
}
func (c ConversationController) SetUserConversations(ctx context.Context, ownerUserID string, conversations []*table.ConversationModel) error {
	panic("implement me")
}

var _ ConversationInterface = (*ConversationController)(nil)

type ConversationDataBaseInterface interface {
	//GetUserIDExistConversation 获取拥有该会话的的用户ID列表
	GetUserIDExistConversation(ctx context.Context, userIDList []string, conversationID string) ([]string, error)
	//UpdateUserConversationFiled 更新用户该会话的属性信息
	UpdateUsersConversationFiled(ctx context.Context, UserIDList []string, conversationID string, args map[string]interface{}) error
	//CreateConversation 创建一批新的会话
	CreateConversation(ctx context.Context, conversations []*table.ConversationModel) error
	//SyncPeerUserPrivateConversation 同步对端私聊会话内部保证事务操作
	SyncPeerUserPrivateConversationTx(ctx context.Context, conversation *table.ConversationModel) error
	//FindConversations 根据会话ID获取某个用户的多个会话
	FindConversations(ctx context.Context, ownerUserID string, conversationID []string) ([]*table.ConversationModel, error)
	//GetUserAllConversation 获取一个用户在服务器上所有的会话
	GetUserAllConversation(ctx context.Context, ownerUserID string) ([]*table.ConversationModel, error)
	//SetUserConversations 设置用户多个会话属性，如果会话不存在则创建，否则更新,内部保证原子性
	SetUserConversations(ctx context.Context, ownerUserID string, conversations []*table.ConversationModel) error
}
type ConversationDataBase struct {
	db    relation.Conversation
	cache cache.ConversationCache
}

func (c ConversationDataBase) GetUserIDExistConversation(ctx context.Context, userIDList []string, conversationID string) ([]string, error) {
	panic("implement me")
}

func (c ConversationDataBase) UpdateUsersConversationFiled(ctx context.Context, UserIDList []string, conversationID string, args map[string]interface{}) error {
	panic("implement me")
}

func (c ConversationDataBase) CreateConversation(ctx context.Context, conversations []*table.ConversationModel) error {
	panic("implement me")
}

func (c ConversationDataBase) SyncPeerUserPrivateConversationTx(ctx context.Context, conversation *table.ConversationModel) error {
	panic("implement me")
}

func (c ConversationDataBase) FindConversations(ctx context.Context, ownerUserID string, conversationID []string) ([]*table.ConversationModel, error) {
	panic("implement me")
}

func (c ConversationDataBase) GetUserAllConversation(ctx context.Context, ownerUserID string) ([]*table.ConversationModel, error) {
	panic("implement me")
}

func (c ConversationDataBase) SetUserConversations(ctx context.Context, ownerUserID string, conversations []*table.ConversationModel) error {
	panic("implement me")
}

func NewConversationDataBase(db relation.Conversation, cache cache.ConversationCache) *ConversationDataBase {
	return &ConversationDataBase{db: db, cache: cache}
}

//func NewConversationController(db *gorm.DB, rdb redis.UniversalClient) ConversationInterface {
//	groupController := &ConversationController{database: newGroupDatabase(db, rdb, mgoClient)}
//	return groupController
//}
