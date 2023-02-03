package relation

import (
	"Open_IM/pkg/common/db/table"
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/utils"
	"context"
	"gorm.io/gorm"
)

type Conversation interface {
	TableName() string
	Create(ctx context.Context, conversations []*table.ConversationModel) (err error)
	Delete(ctx context.Context, groupIDs []string) (err error)
	UpdateByMap(ctx context.Context, groupID string, args map[string]interface{}) (err error)
	Update(ctx context.Context, groups []*table.ConversationModel) (err error)
	Find(ctx context.Context, groupIDs []string) (groups []*table.ConversationModel, err error)
	Take(ctx context.Context, groupID string) (group *table.ConversationModel, err error)
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

func (c *ConversationGorm) Create(ctx context.Context, conversations []*table.ConversationModel) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "conversations", conversations)
	}()
	return utils.Wrap(getDBConn(g.DB, tx).Create(&conversations).Error, "")
}

func (c *ConversationGorm) Delete(ctx context.Context, groupIDs []string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupIDs", groupIDs)
	}()
	return utils.Wrap(getDBConn(g.DB, tx).Where("group_id in (?)", groupIDs).Delete(&table.ConversationModel{}).Error, "")
}

func (c *ConversationGorm) UpdateByMap(ctx context.Context, groupID string, args map[string]interface{}) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "args", args)
	}()
	return utils.Wrap(getDBConn(g.DB, tx).Where("group_id = ?", groupID).Model(g).Updates(args).Error, "")
}

func (c *ConversationGorm) Update(ctx context.Context, groups []*table.ConversationModel) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groups", groups)
	}()
	return utils.Wrap(getDBConn(g.DB, tx).Updates(&groups).Error, "")
}

func (c *ConversationGorm) Find(ctx context.Context, groupIDs []string) (groups []*table.ConversationModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupIDs", groupIDs, "groups", groups)
	}()
	return groups, utils.Wrap(getDBConn(g.DB, tx).Where("group_id in (?)", groupIDs).Find(&groups).Error, "")
}

func (c *ConversationGorm) Take(ctx context.Context, groupID string) (group *table.ConversationModel, err error) {
	group = &Group{}
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "group", *group)
	}()
	return group, utils.Wrap(getDBConn(g.DB, tx).Where("group_id = ?", groupID).Take(group).Error, "")
}
