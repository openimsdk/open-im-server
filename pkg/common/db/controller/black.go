package controller

import (
	"Open_IM/pkg/common/db/relation"
	"context"
	"gorm.io/gorm"
)

type BlackInterface interface {
	// Create 增加黑名单
	Create(ctx context.Context, blacks []*relation.Black) (err error)
	// Delete 删除黑名单
	Delete(ctx context.Context, blacks []*relation.Black) (err error)
	// FindOwnerBlacks 获取黑名单列表
	FindOwnerBlacks(ctx context.Context, ownerUserID string, pageNumber, showNumber int32) (blackList []*relation.Black, err error)
	// CheckIn 检查user2是否在user1的黑名单列表中(inUser1Blacks==true) 检查user1是否在user2的黑名单列表中(inUser2Blacks==true)
	CheckIn(ctx context.Context, ownerUserID, blackUserID string) (inUser1Blacks bool, inUser2Blacks bool, err error)
}

type BlackController struct {
	database BlackDatabaseInterface
}

func NewBlackController(db *gorm.DB) *BlackController {
	return &BlackController{database: NewBlackDatabase(db)}
}

// Create 增加黑名单
func (b *BlackController) Create(ctx context.Context, blacks []*relation.Black) (err error) {
}

// Delete 删除黑名单
func (b *BlackController) Delete(ctx context.Context, blacks []*relation.Black) (err error) {
}

// FindOwnerBlacks 获取黑名单列表
func (b *BlackController) FindOwnerBlacks(ctx context.Context, ownerUserID string, pageNumber, showNumber int32) (blackList []*relation.Black, err error) {
}

// CheckIn 检查user2是否在user1的黑名单列表中(inUser1Blacks==true) 检查user1是否在user2的黑名单列表中(inUser2Blacks==true)
func (b *BlackController) CheckIn(ctx context.Context, ownerUserID, blackUserID string) (inUser1Blacks bool, inUser2Blacks bool, err error) {
}

type BlackDatabaseInterface interface {
	// Create 增加黑名单
	Create(ctx context.Context, blacks []*relation.Black) (err error)
	// Delete 删除黑名单
	Delete(ctx context.Context, blacks []*relation.Black) (err error)
	// FindOwnerBlacks 获取黑名单列表
	FindOwnerBlacks(ctx context.Context, ownerUserID string, pageNumber, showNumber int32) (blackList []*relation.Black, err error)
	// CheckIn 检查user2是否在user1的黑名单列表中(inUser1Blacks==true) 检查user1是否在user2的黑名单列表中(inUser2Blacks==true)
	CheckIn(ctx context.Context, ownerUserID, blackUserID string) (inUser1Blacks bool, inUser2Blacks bool, err error)
}
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
}

// Delete 删除黑名单
func (b *BlackDatabase) Delete(ctx context.Context, blacks []*relation.Black) (err error) {
}

// FindOwnerBlacks 获取黑名单列表
func (b *BlackDatabase) FindOwnerBlacks(ctx context.Context, ownerUserID string, pageNumber, showNumber int32) (blackList []*relation.Black, err error) {
}

// CheckIn 检查user2是否在user1的黑名单列表中(inUser1Blacks==true) 检查user1是否在user2的黑名单列表中(inUser2Blacks==true)
func (b *BlackDatabase) CheckIn(ctx context.Context, ownerUserID, blackUserID string) (inUser1Blacks bool, inUser2Blacks bool, err error) {
}
