package controller

import (
	"Open_IM/pkg/common/db/table/relation"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"gorm.io/gorm"
)

type BlackDatabase interface {
	// Create 增加黑名单
	Create(ctx context.Context, blacks []*relation.BlackModel) (err error)
	// Delete 删除黑名单
	Delete(ctx context.Context, blacks []*relation.BlackModel) (err error)
	// FindOwnerBlacks 获取黑名单列表
	FindOwnerBlacks(ctx context.Context, ownerUserID string, pageNumber, showNumber int32) (blacks []*relation.BlackModel, total int64, err error)
	// CheckIn 检查user2是否在user1的黑名单列表中(inUser1Blacks==true) 检查user1是否在user2的黑名单列表中(inUser2Blacks==true)
	CheckIn(ctx context.Context, userID1, userID2 string) (inUser1Blacks bool, inUser2Blacks bool, err error)
}

type blackDatabase struct {
	black relation.BlackModelInterface
}

func NewBlackDatabase(black relation.BlackModelInterface) BlackDatabase {
	return &blackDatabase{black}
}

// Create 增加黑名单
func (b *blackDatabase) Create(ctx context.Context, blacks []*relation.BlackModel) (err error) {
	return b.black.Create(ctx, blacks)
}

// Delete 删除黑名单
func (b *blackDatabase) Delete(ctx context.Context, blacks []*relation.BlackModel) (err error) {
	return b.black.Delete(ctx, blacks)
}

// FindOwnerBlacks 获取黑名单列表
func (b *blackDatabase) FindOwnerBlacks(ctx context.Context, ownerUserID string, pageNumber, showNumber int32) (blacks []*relation.BlackModel, total int64, err error) {
	return b.black.FindOwnerBlacks(ctx, ownerUserID, pageNumber, showNumber)
}

// CheckIn 检查user2是否在user1的黑名单列表中(inUser1Blacks==true) 检查user1是否在user2的黑名单列表中(inUser2Blacks==true)
func (b *blackDatabase) CheckIn(ctx context.Context, userID1, userID2 string) (inUser1Blacks bool, inUser2Blacks bool, err error) {
	_, err = b.black.Take(ctx, userID1, userID2)
	if err != nil {
		if errors.Unwrap(err) != gorm.ErrRecordNotFound {
			return
		}
		inUser1Blacks = false
	} else {
		inUser1Blacks = true
	}

	inUser2Blacks = true
	_, err = b.black.Take(ctx, userID2, userID1)
	if err != nil {
		if utils.Unwrap(err) != gorm.ErrRecordNotFound {
			return
		}
		inUser2Blacks = false
	} else {
		inUser2Blacks = true
	}
	return
}
