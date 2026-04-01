package controller

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/db/pagination"
)

// UserGlobalBlackDatabase 全局黑名单业务接口
type UserGlobalBlackDatabase interface {
	// AddBlack 将用户加入全局黑名单
	AddBlack(ctx context.Context, blacks []*model.UserGlobalBlack) error
	// RemoveBlack 按 userID 将用户从全局黑名单移除
	RemoveBlack(ctx context.Context, userIDs []string) error
	// IsBlocked 检查用户是否在全局黑名单
	IsBlocked(ctx context.Context, userID string) (bool, error)
	// GetBlackList 分页获取黑名单列表
	GetBlackList(ctx context.Context, pagination pagination.Pagination) (count int64, blacks []*model.UserGlobalBlack, err error)
}

type userGlobalBlackDatabase struct {
	db database.UserGlobalBlack
}

func NewUserGlobalBlackDatabase(db database.UserGlobalBlack) UserGlobalBlackDatabase {
	return &userGlobalBlackDatabase{db: db}
}

func (u *userGlobalBlackDatabase) AddBlack(ctx context.Context, blacks []*model.UserGlobalBlack) error {
	return u.db.Add(ctx, blacks)
}

func (u *userGlobalBlackDatabase) RemoveBlack(ctx context.Context, userIDs []string) error {
	return u.db.Remove(ctx, userIDs)
}

func (u *userGlobalBlackDatabase) IsBlocked(ctx context.Context, userID string) (bool, error) {
	return u.db.IsBlocked(ctx, userID)
}

func (u *userGlobalBlackDatabase) GetBlackList(ctx context.Context, pagination pagination.Pagination) (int64, []*model.UserGlobalBlack, error) {
	return u.db.Page(ctx, pagination)
}
