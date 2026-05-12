package database

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

// UserMute 用户静音持久化接口（支持好友与非好友）
type UserMute interface {
	// Upsert 新增或更新静音记录
	Upsert(ctx context.Context, mute *model.UserMute) error
	// Delete 取消静音（删除记录）
	Delete(ctx context.Context, ownerUserID, mutedUserID string) error
	// IsMuted 检查 ownerUserID 是否对 mutedUserID 设置了有效的静音
	IsMuted(ctx context.Context, ownerUserID, mutedUserID string) (bool, error)
	// Get 按 owner + muted 查询一条记录；不存在则 (nil, nil)
	Get(ctx context.Context, ownerUserID, mutedUserID string) (*model.UserMute, error)
}
