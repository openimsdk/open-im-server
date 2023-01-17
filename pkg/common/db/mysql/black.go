package mysql

import (
	"Open_IM/pkg/common/trace_log"
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
	black.DB = initModel(db, &black)
	return &black
}

func (b *Black) Create(ctx context.Context, blacks []*Black) (err error) {
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetFuncName(1), err, "blacks", blacks)
	}()
	return utils.Wrap(b.DB.Create(&blacks).Error, "")
}

func (b *Black) Delete(ctx context.Context, blacks []*Black) (err error) {
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetFuncName(1), err, "blacks", blacks)
	}()
	return utils.Wrap(GroupMemberDB.Delete(blacks).Error, "")
}

func (b *Black) UpdateByMap(ctx context.Context, ownerUserID, blockUserID string, args map[string]interface{}) (err error) {
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID, "blockUserID", blockUserID, "args", args)
	}()
	return utils.Wrap(b.DB.Where("block_user_id = ? and block_user_id = ?", ownerUserID, blockUserID).Updates(args).Error, "")
}

func (b *Black) Update(ctx context.Context, blacks []*Black) (err error) {
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetFuncName(1), err, "blacks", blacks)
	}()
	return utils.Wrap(b.DB.Updates(&blacks).Error, "")
}

func (b *Black) Find(ctx context.Context, blacks []*Black) (blackList []*Black, err error) {
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetFuncName(1), err, "blacks", blacks, "blackList", blackList)
	}()
	var where [][]interface{}
	for _, black := range blacks {
		where = append(where, []interface{}{black.OwnerUserID, black.BlockUserID})
	}
	return blackList, utils.Wrap(GroupMemberDB.Where("(owner_user_id, block_user_id) in ?", where).Find(&blackList).Error, "")
}

func (b *Black) Take(ctx context.Context, ownerUserID, blockUserID string) (black *Black, err error) {
	black = &Black{}
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID, "blockUserID", blockUserID, "black", *black)
	}()
	return black, utils.Wrap(b.DB.Where("owner_user_id = ? and block_user_id = ?", ownerUserID, blockUserID).Take(black).Error, "")
}

func (b *Black) FindByOwnerUserID(ctx context.Context, ownerUserID string) (blackList []*Black, err error) {
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetFuncName(1), err, "ownerUserID", ownerUserID, "blackList", blackList)
	}()
	return blackList, utils.Wrap(GroupMemberDB.Where("owner_user_id = ?", ownerUserID).Find(&blackList).Error, "")
}

func InsertInToUserBlackList(ctx context.Context, black Black) (err error) {
	defer trace_log.SetCtxDebug(ctx, utils.GetSelfFuncName(), err, "black", black)
	black.CreateTime = time.Now()
	err = BlackDB.Create(black).Error
	return err
}

func CheckBlack(ownerUserID, blockUserID string) error {
	var black Black
	return BlackDB.Where("owner_user_id=? and block_user_id=?", ownerUserID, blockUserID).Find(&black).Error
}

func RemoveBlackList(ownerUserID, blockUserID string) error {
	err := BlackDB.Where("owner_user_id=? and block_user_id=?", ownerUserID, blockUserID).Delete(Black{}).Error
	return utils.Wrap(err, "RemoveBlackList failed")
}

func GetBlackListByUserID(ownerUserID string) ([]Black, error) {
	var blackListUsersInfo []Black
	err := BlackDB.Where("owner_user_id=?", ownerUserID).Find(&blackListUsersInfo).Error
	if err != nil {
		return nil, err
	}
	return blackListUsersInfo, nil
}

func GetBlackIDListByUserID(ownerUserID string) ([]string, error) {
	var blackIDList []string
	err := b.db.Where("owner_user_id=?", ownerUserID).Pluck("block_user_id", &blackIDList).Error
	if err != nil {
		return nil, err
	}
	return blackIDList, nil
}
