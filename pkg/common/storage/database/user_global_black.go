package database

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/db/pagination"
)

// UserGlobalBlack 全局黑名单持久化接口
type UserGlobalBlack interface {
	// Add 批量添加用户到全局黑名单
	Add(ctx context.Context, blacks []*model.UserGlobalBlack) error
	// Remove 按 userID 从全局黑名单移除用户
	Remove(ctx context.Context, userIDs []string) error
	// Find 查询指定用户是否在黑名单（返回在黑名单中的记录）
	Find(ctx context.Context, userIDs []string) ([]*model.UserGlobalBlack, error)
	// IsBlocked 检查单个用户是否在黑名单
	IsBlocked(ctx context.Context, userID string) (bool, error)
	// Page 分页查询黑名单列表
	Page(ctx context.Context, pagination pagination.Pagination) (count int64, blacks []*model.UserGlobalBlack, err error)
}
