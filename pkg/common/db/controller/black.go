package controller

import (
	"Open_IM/pkg/common/db/table/relation"
	"context"
	"errors"
	"gorm.io/gorm"
)

type BlackInterface interface {
	// Create 增加黑名单
	Create(ctx context.Context, blacks []*relation.BlackModel) (err error)
	// Delete 删除黑名单
	Delete(ctx context.Context, blacks []*relation.BlackModel) (err error)
	// FindOwnerBlacks 获取黑名单列表
	FindOwnerBlacks(ctx context.Context, ownerUserID string, pageNumber, showNumber int32) (blacks []*relation.BlackModel, total int64, err error)
	// CheckIn 检查user2是否在user1的黑名单列表中(inUser1Blacks==true) 检查user1是否在user2的黑名单列表中(inUser2Blacks==true)
	CheckIn(ctx context.Context, userID1, userID2 string) (inUser1Blacks bool, inUser2Blacks bool, err error)
}

type BlackController struct {
	database BlackDatabaseInterface
}

func NewBlackController(db *gorm.DB) *BlackController {
	return &BlackController{database: NewBlackDatabase(db)}
}

// Create 增加黑名单
func (b *BlackController) Create(ctx context.Context, blacks []*relation.Black) (err error) {
	return b.database.Create(ctx, blacks)
}

// Delete 删除黑名单
func (b *BlackController) Delete(ctx context.Context, blacks []*relation.Black) (err error) {
	return b.database.Delete(ctx, blacks)
}

// FindOwnerBlacks 获取黑名单列表
func (b *BlackController) FindOwnerBlacks(ctx context.Context, ownerUserID string, pageNumber, showNumber int32) (blackList []*relation.Black, total int64, err error) {
	return b.database.FindOwnerBlacks(ctx, ownerUserID, pageNumber, showNumber)
}

// CheckIn 检查user2是否在user1的黑名单列表中(inUser1Blacks==true) 检查user1是否在user2的黑名单列表中(inUser2Blacks==true)
func (b *BlackController) CheckIn(ctx context.Context, userID1, userID2 string) (inUser1Blacks bool, inUser2Blacks bool, err error) {
	return b.database.CheckIn(ctx, userID1, userID2)
}

type BlackDatabaseInterface interface {
	// Create 增加黑名单
	Create(ctx context.Context, blacks []*relation.Black) (err error)
	// Delete 删除黑名单
	Delete(ctx context.Context, blacks []*relation.Black) (err error)
	// FindOwnerBlacks 获取黑名单列表
	FindOwnerBlacks(ctx context.Context, ownerUserID string, pageNumber, showNumber int32) (blacks []*relation.Black, total int64, err error)
	// CheckIn 检查user2是否在user1的黑名单列表中(inUser1Blacks==true) 检查user1是否在user2的黑名单列表中(inUser2Blacks==true)
	CheckIn(ctx context.Context, userID1, userID2 string) (inUser1Blacks bool, inUser2Blacks bool, err error)
}

type BlackDatabase struct {
	sqlDB *relation.Black
}

func NewBlackDatabase(db *gorm.DB) *BlackDatabase {
	sqlDB := relation.NewBlack(db)
	database := &BlackDatabase{
		sqlDB: sqlDB,
	}
	return database
}

// Create 增加黑名单
func (b *BlackDatabase) Create(ctx context.Context, blacks []*relation.Black) (err error) {
	return b.sqlDB.Create(ctx, blacks)
}

// Delete 删除黑名单
func (b *BlackDatabase) Delete(ctx context.Context, blacks []*relation.Black) (err error) {
	return b.sqlDB.Delete(ctx, blacks)
}

// FindOwnerBlacks 获取黑名单列表
func (b *BlackDatabase) FindOwnerBlacks(ctx context.Context, ownerUserID string, pageNumber, showNumber int32) (blacks []*relation.Black, total int64, err error) {
	return b.sqlDB.FindOwnerBlacks(ctx, ownerUserID, pageNumber, showNumber)
}

// CheckIn 检查user2是否在user1的黑名单列表中(inUser1Blacks==true) 检查user1是否在user2的黑名单列表中(inUser2Blacks==true)
func (b *BlackDatabase) CheckIn(ctx context.Context, userID1, userID2 string) (inUser1Blacks bool, inUser2Blacks bool, err error) {
	_, err = b.sqlDB.Take(ctx, userID1, userID2)
	if err != nil {
		if errors.Unwrap(err) != gorm.ErrRecordNotFound {
			return
		}
		inUser1Blacks = false
	} else {
		inUser1Blacks = true
	}

	inUser2Blacks = true
	_, err = b.sqlDB.Take(ctx, userID2, userID1)
	if err != nil {
		if errors.Unwrap(err) != gorm.ErrRecordNotFound {
			return
		}
		inUser2Blacks = false
	} else {
		inUser2Blacks = true
	}
	return
}
