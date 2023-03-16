package relation

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"gorm.io/gorm"
)

type Conversation interface {
	Create(ctx context.Context, conversations []*relation.ConversationModel) (err error)
	Delete(ctx context.Context, groupIDs []string) (err error)
	UpdateByMap(ctx context.Context, userIDList []string, conversationID string, args map[string]interface{}) (err error)
	Update(ctx context.Context, conversations []*relation.ConversationModel) (err error)
	Find(ctx context.Context, ownerUserID string, conversationIDs []string) (conversations []*relation.ConversationModel, err error)
	FindUserID(ctx context.Context, userIDList []string, conversationID string) ([]string, error)
	FindUserIDAllConversationID(ctx context.Context, userID string) ([]string, error)
	Take(ctx context.Context, userID, conversationID string) (conversation *relation.ConversationModel, err error)
	FindConversationID(ctx context.Context, userID string, conversationIDList []string) (existConversationID []string, err error)
	FindRecvMsgNotNotifyUserIDs(ctx context.Context, groupID string) ([]string, error)
	NewTx(tx any) Conversation
}
type ConversationGorm struct {
	DB *gorm.DB
}

func NewConversationGorm(DB *gorm.DB) Conversation {
	return &ConversationGorm{DB: DB}
}

func (c *ConversationGorm) NewTx(tx any) Conversation {
	return &ConversationGorm{DB: tx.(*gorm.DB)}
}

func (c *ConversationGorm) Create(ctx context.Context, conversations []*relation.ConversationModel) (err error) {
	return utils.Wrap(c.DB.Create(&conversations).Error, "")
}

func (c *ConversationGorm) Delete(ctx context.Context, groupIDs []string) (err error) {
	return utils.Wrap(c.DB.Where("group_id in (?)", groupIDs).Delete(&relation.ConversationModel{}).Error, "")
}

func (c *ConversationGorm) UpdateByMap(ctx context.Context, userIDList []string, conversationID string, args map[string]interface{}) (err error) {
	return utils.Wrap(c.DB.Model(&relation.ConversationModel{}).Where("owner_user_id IN (?) and  conversation_id=?", userIDList, conversationID).Updates(args).Error, "")
}

func (c *ConversationGorm) Update(ctx context.Context, conversations []*relation.ConversationModel) (err error) {
	return utils.Wrap(c.DB.Updates(&conversations).Error, "")
}

func (c *ConversationGorm) Find(ctx context.Context, ownerUserID string, conversationIDs []string) (conversations []*relation.ConversationModel, err error) {
	err = utils.Wrap(c.DB.Where("owner_user_id=? and conversation_id IN (?)", ownerUserID, conversationIDs).Find(&conversations).Error, "")
	return conversations, err
}

func (c *ConversationGorm) Take(ctx context.Context, userID, conversationID string) (conversation *relation.ConversationModel, err error) {
	cc := &relation.ConversationModel{}
	return cc, utils.Wrap(c.DB.Where("conversation_id = ? And owner_user_id = ?", conversationID, userID).Take(cc).Error, "")
}

func (c *ConversationGorm) FindUserID(ctx context.Context, userIDList []string, conversationID string) (existUserID []string, err error) {
	return existUserID, utils.Wrap(c.DB.Where(" owner_user_id IN (?) and conversation_id=?", userIDList, conversationID).Pluck("owner_user_id", &existUserID).Error, "")
}

func (c *ConversationGorm) FindConversationID(ctx context.Context, userID string, conversationIDList []string) (existConversationID []string, err error) {
	return existConversationID, utils.Wrap(c.DB.Where(" conversation_id IN (?) and owner_user_id=?", conversationIDList, userID).Pluck("conversation_id", &existConversationID).Error, "")
}

func (c *ConversationGorm) FindUserIDAllConversationID(ctx context.Context, userID string) (conversationIDList []string, err error) {
	return conversationIDList, utils.Wrap(c.DB.Model(&relation.ConversationModel{}).Where("owner_user_id=?", userID).Pluck("conversation_id", &conversationIDList).Error, "")
}

func (c *ConversationGorm) FindRecvMsgNotNotifyUserIDs(ctx context.Context, groupID string) (userIDs []string, err error) {
	return userIDs, utils.Wrap(c.DB.Model(&relation.ConversationModel{}).Where("group_id = ? and recv_msg_opt = ?", groupID, constant.ReceiveNotNotifyMessage).Pluck("user_id", &userIDs).Error, "")
}
