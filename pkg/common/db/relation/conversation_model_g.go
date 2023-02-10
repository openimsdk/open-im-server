package relation

import (
	"Open_IM/pkg/common/db/table/relation"
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/utils"
	"context"
	"gorm.io/gorm"
)

type Conversation interface {
	TableName() string
	Create(ctx context.Context, conversations []*relation.ConversationModel, tx ...any) (err error)
	Delete(ctx context.Context, groupIDs []string) (err error)
	UpdateByMap(ctx context.Context, userIDList []string, conversationID string, args map[string]interface{}, tx ...any) (err error)
	Update(ctx context.Context, conversations []*relation.ConversationModel, tx ...any) (err error)
	Find(ctx context.Context, ownerUserID string, conversationIDs []string, tx ...any) (conversations []*relation.ConversationModel, err error)
	FindUserID(ctx context.Context, userIDList []string, conversationID string, tx ...any) ([]string, error)
	FindUserIDAllConversationID(ctx context.Context, userID string, tx ...any) ([]string, error)
	Take(ctx context.Context, userID, conversationID string, tx ...any) (conversation *relation.ConversationModel, err error)
	FindConversationID(ctx context.Context, userID string, conversationIDList []string, tx ...any) (existConversationID []string, err error)
}
type ConversationGorm struct {
	DB *gorm.DB
}

func (c *ConversationGorm) TableName() string {
	panic("implement me")
}

func NewConversationGorm(DB *gorm.DB) Conversation {
	return &ConversationGorm{DB: DB}
}

func (c *ConversationGorm) Create(ctx context.Context, conversations []*relation.ConversationModel, tx ...any) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "conversations", conversations)
	}()
	return utils.Wrap(getDBConn(c.DB, tx).Create(&conversations).Error, "")
}

func (c *ConversationGorm) Delete(ctx context.Context, groupIDs []string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupIDs", groupIDs)
	}()
	return utils.Wrap(getDBConn(g.DB, tx).Where("group_id in (?)", groupIDs).Delete(&relation.ConversationModel{}).Error, "")
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
	var newConversations []relation.ConversationModel
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID, "groups", conversations)
	}()
	err = utils.Wrap(getDBConn(c.DB, tx).Where("owner_user_id=? and conversation_id IN (?)", ownerUserID, conversationIDs).Find(&newConversations).Error, "")
	for _, v := range newConversations {
		v1 := v
		conversations = append(conversations, &v1)
	}
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
