package controller

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

// UserMuteDatabase 用户静音业务接口
type UserMuteDatabase interface {
	// Upsert 新增或更新静音记录
	Upsert(ctx context.Context, mute *model.UserMute) error
	// Delete 取消静音
	Delete(ctx context.Context, ownerUserID, mutedUserID string) error
	// IsMuted 检查 ownerUserID 是否对 mutedUserID 设置了有效静音
	IsMuted(ctx context.Context, ownerUserID, mutedUserID string) (bool, error)
	// Get 查询静音记录；不存在则 (nil, nil)
	Get(ctx context.Context, ownerUserID, mutedUserID string) (*model.UserMute, error)
}

type userMuteDatabase struct {
	db database.UserMute
}

func NewUserMuteDatabase(db database.UserMute) UserMuteDatabase {
	return &userMuteDatabase{db: db}
}

func (u *userMuteDatabase) Upsert(ctx context.Context, mute *model.UserMute) error {
	return u.db.Upsert(ctx, mute)
}

func (u *userMuteDatabase) Delete(ctx context.Context, ownerUserID, mutedUserID string) error {
	return u.db.Delete(ctx, ownerUserID, mutedUserID)
}

func (u *userMuteDatabase) IsMuted(ctx context.Context, ownerUserID, mutedUserID string) (bool, error) {
	return u.db.IsMuted(ctx, ownerUserID, mutedUserID)
}

func (u *userMuteDatabase) Get(ctx context.Context, ownerUserID, mutedUserID string) (*model.UserMute, error) {
	return u.db.Get(ctx, ownerUserID, mutedUserID)
}
