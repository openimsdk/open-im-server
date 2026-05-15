package database

import (
	"context"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

// UserOfflineRecord 管理 user_offline_record 集合。
// 集合中的每条记录代表一个当前处于全平台离线状态的用户及其首次全离线时刻。
type UserOfflineRecord interface {
	// Upsert 写入用户的离线记录；若记录已存在则不覆盖（保留最早的离线时刻）。
	// deadline = offlineTime + delete_account_interval，供 FindExpiredUsers 快速过滤。
	Upsert(ctx context.Context, userID string, offlineTime, deadline time.Time) error

	// RefreshOfflineTime 将已存在的离线记录的 offline_time 与 delete_user_deadline
	// 同时刷新，使删除倒计时从 newOfflineTime 重新起算。
	// 若记录不存在（用户在线）则无操作。
	RefreshOfflineTime(ctx context.Context, userID string, newOfflineTime, newDeadline time.Time) error

	// Delete 删除用户的离线记录（用户重新上线时调用）。
	Delete(ctx context.Context, userID string) error

	// FindExpiredUsers 返回 delete_user_deadline <= now 的用户（$lookup user 集合获取完整信息）。
	// limit 限制单次返回条数，防止单批处理量过大。
	FindExpiredUsers(ctx context.Context, now time.Time, limit int) ([]*model.User, error)
}
