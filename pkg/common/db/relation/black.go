package relation

import (
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/utils"
	"context"
	"gorm.io/gorm"
	"time"
)

type Black struct {
	OwnerUserID    string    `gorm:"column:owner_user_id;primary_key;size:64"`
	BlockUserID    string    `gorm:"column:block_user_id;primary_key;size:64"`
	CreateTime     time.Time `gorm:"column:create_time"`
	AddSource      int32     `gorm:"column:add_source"`
	OperatorUserID string    `gorm:"column:operator_user_id;size:64"`
	Ex             string    `gorm:"column:ex;size:1024"`
	DB             *gorm.DB  `gorm:"-"`
}

func NewBlack(db *gorm.DB) *Black {
	var black Black
	black.DB = db.Model(&Black{})
	return &black
}

func (b *Black) Create(ctx context.Context, blacks []*Black) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "blacks", blacks)
	}()
	return utils.Wrap(b.DB.Create(&blacks).Error, "")
}

func (b *Black) Delete(ctx context.Context, blacks []*Black) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "blacks", blacks)
	}()
	return utils.Wrap(b.DB.Delete(blacks).Error, "")
}

func (b *Black) UpdateByMap(ctx context.Context, ownerUserID, blockUserID string, args map[string]interface{}) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID, "blockUserID", blockUserID, "args", args)
	}()
	return utils.Wrap(b.DB.Where("block_user_id = ? and block_user_id = ?", ownerUserID, blockUserID).Updates(args).Error, "")
}

func (b *Black) Update(ctx context.Context, blacks []*Black) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "blacks", blacks)
	}()
	return utils.Wrap(b.DB.Updates(&blacks).Error, "")
}

func (b *Black) Find(ctx context.Context, blacks []*Black) (blackList []*Black, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "blacks", blacks, "blackList", blackList)
	}()
	var where [][]interface{}
	for _, black := range blacks {
		where = append(where, []interface{}{black.OwnerUserID, black.BlockUserID})
	}
	return blackList, utils.Wrap(b.DB.Where("(owner_user_id, block_user_id) in ?", where).Find(&blackList).Error, "")
}

func (b *Black) GetBlackIDs(ctx context.Context, ownerUserID string) (userIDs []string, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID, "userIDs", userIDs)
	}()
	err = utils.Wrap(b.DB.Where("owner_user_id = ?", ownerUserID).Pluck("block_user_id", &userIDs).Error, "")
	return userIDs, err
}

func (b *Black) Take(ctx context.Context, ownerUserID, blockUserID string) (black *Black, err error) {
	black = &Black{}
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID, "blockUserID", blockUserID, "black", *black)
	}()
	return black, utils.Wrap(b.DB.Where("owner_user_id = ? and block_user_id = ?", ownerUserID, blockUserID).Take(black).Error, "")
}

func (b *Black) FindOwnerBlacks(ctx context.Context, ownerUserID string, pageNumber, showNumber int32) (blacks []*Black, total int64, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID, "blacks", blacks)
	}()
	err = b.DB.Model(b).Count(&total).Error
	if err != nil {
		return nil, 0, utils.Wrap(err, "")
	}
	err = utils.Wrap(b.DB.Limit(int(showNumber)).Offset(int(pageNumber*showNumber)).Find(&blacks).Error, "")
	return
}
