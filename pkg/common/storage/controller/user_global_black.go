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
	// IsBlocked 检查用户是否在全局黑名单（含冻结）
	IsBlocked(ctx context.Context, userID string) (bool, error)
	// GetStatus 返回用户限制状态：0=正常，1=冻结，2=黑名单
	GetStatus(ctx context.Context, userID string) (int32, error)
	// FindBlocked 批量查询哪些 userID 在全局黑名单中，返回被封禁的记录
	FindBlocked(ctx context.Context, userIDs []string) ([]*model.UserGlobalBlack, error)
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

func (u *userGlobalBlackDatabase) GetStatus(ctx context.Context, userID string) (int32, error) {
	return u.db.GetStatus(ctx, userID)
}

func (u *userGlobalBlackDatabase) GetBlackList(ctx context.Context, pagination pagination.Pagination) (int64, []*model.UserGlobalBlack, error) {
	return u.db.Page(ctx, pagination)
}

func (u *userGlobalBlackDatabase) FindBlocked(ctx context.Context, userIDs []string) ([]*model.UserGlobalBlack, error) {
	return u.db.Find(ctx, userIDs)
}
