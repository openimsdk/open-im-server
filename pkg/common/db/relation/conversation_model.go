package relation

import (
	"Open_IM/pkg/common/db/table/relation"
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/utils"
	"context"
	"gorm.io/gorm"
)

type Conversation interface {
	Create(ctx context.Context, conversations []*relation.ConversationModel, tx ...any) (err error)
	Delete(ctx context.Context, groupIDs []string, tx ...any) (err error)
	UpdateByMap(ctx context.Context, userIDList []string, conversationID string, args map[string]interface{}, tx ...any) (err error)
	Update(ctx context.Context, conversations []*relation.ConversationModel, tx ...any) (err error)
	Find(ctx context.Context, ownerUserID string, conversationIDs []string, tx ...any) (conversations []*relation.ConversationModel, err error)
	FindUserID(ctx context.Context, userIDList []string, conversationID string, tx ...any) ([]string, error)
	FindUserIDAllConversationID(ctx context.Context, userID string, tx ...any) ([]string, error)
	Take(ctx context.Context, userID, conversationID string, tx ...any) (conversation *relation.ConversationModel, err error)
	FindConversationID(ctx context.Context, userID string, conversationIDList []string, tx ...any) (existConversationID []string, err error)
	Transaction(func(tx any) error) error
	NewTx(tx any) Conversation
}
type ConversationGorm struct {
	DB *gorm.DB
}

func NewConversationGorm(DB *gorm.DB) Conversation {
	return &ConversationGorm{DB: DB}
}
func (c *ConversationGorm) Transaction(fn func(tx any) error) error {
	return c.DB.Transaction(func(tx *gorm.DB) error {
		return fn(tx)
	})
}

func (c *ConversationGorm) NewTx(tx any) Conversation {
	return &ConversationGorm{DB: tx.(*gorm.DB)}
}

func (c *ConversationGorm) Create(ctx context.Context, conversations []*relation.ConversationModel, tx ...any) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "conversations", conversations)
	}()
	return utils.Wrap(getDBConn(c.DB, tx).Create(&conversations).Error, "")
}

func (c *ConversationGorm) Delete(ctx context.Context, groupIDs []string, tx ...any) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupIDs", groupIDs)
	}()
	return utils.Wrap(getDBConn(c.DB, tx).Where("group_id in (?)", groupIDs).Delete(&relation.ConversationModel{}).Error, "")
}

func (c *ConversationGorm) UpdateByMap(ctx context.Context, userIDList []string, conversationID string, args map[string]interface{}, tx ...any) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userIDList", userIDList, "conversationID", conversationID)
	}()
	return utils.Wrap(getDBConn(c.DB, tx).Model(&relation.ConversationModel{}).Where("owner_user_id IN (?) and  conversation_id=?", userIDList, conversationID).Updates(args).Error, "")
}

func (c *ConversationGorm) Update(ctx context.Context, conversations []*relation.ConversationModel, tx ...any) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "conversations", conversations)
	}()
	return utils.Wrap(getDBConn(c.DB, tx).Updates(&conversations).Error, "")
}

func (c *ConversationGorm) Find(ctx context.Context, ownerUserID string, conversationIDs []string, tx ...any) (conversations []*relation.ConversationModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID, "groups", conversations)
	}()
	err = utils.Wrap(getDBConn(c.DB, tx).Where("owner_user_id=? and conversation_id IN (?)", ownerUserID, conversationIDs).Find(&conversations).Error, "")
	return conversations, err
}

func (c *ConversationGorm) Take(ctx context.Context, userID, conversationID string, tx ...any) (conversation *relation.ConversationModel, err error) {
	cc := &relation.ConversationModel{}
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userID", userID, "conversation", *conversation)
	}()
	return cc, utils.Wrap(getDBConn(c.DB, tx).Where("conversation_id = ? And owner_user_id = ?", conversationID, userID).Take(cc).Error, "")
}
func (c *ConversationGorm) FindUserID(ctx context.Context, userIDList []string, conversationID string, tx ...any) (existUserID []string, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userID", userIDList, "existUserID", existUserID)
	}()
	return existUserID, utils.Wrap(getDBConn(c.DB, tx).Where(" owner_user_id IN (?) and conversation_id=?", userIDList, conversationID).Pluck("owner_user_id", &existUserID).Error, "")
}
func (c *ConversationGorm) FindConversationID(ctx context.Context, userID string, conversationIDList []string, tx ...any) (existConversationID []string, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userID", userID, "existConversationIDList", existConversationID)
	}()
	return existConversationID, utils.Wrap(getDBConn(c.DB, tx).Where(" conversation_id IN (?) and owner_user_id=?", conversationIDList, userID).Pluck("conversation_id", &existConversationID).Error, "")
}
func (c *ConversationGorm) FindUserIDAllConversationID(ctx context.Context, userID string, tx ...any) (conversationIDList []string, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userID", userID, "conversationIDList", conversationIDList)
	}()
	return conversationIDList, utils.Wrap(getDBConn(c.DB, tx).Model(&relation.ConversationModel{}).Where("owner_user_id=?", userID).Pluck("conversation_id", &conversationIDList).Error, "")
}
