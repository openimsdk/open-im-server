package database

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

type CryptoDevice interface {
	Create(ctx context.Context, device *model.CryptoDevice) error
	FindByUserID(ctx context.Context, userID string) ([]*model.CryptoDevice, error)
	FindByUserIDAndDeviceID(ctx context.Context, userID, deviceID string) (*model.CryptoDevice, error)
	UpdateStatus(ctx context.Context, userID, deviceID, status string) error
	UpdateLastSeen(ctx context.Context, userID, deviceID string) error
}

type GroupKeyVersion interface {
	Find(ctx context.Context, groupID string) (*model.GroupKeyVersion, error)
	IncrVersion(ctx context.Context, groupID string) (int64, error)
}

type GroupKeyEvent interface {
	Create(ctx context.Context, event *model.GroupKeyEvent) error
	FindSinceVersion(ctx context.Context, groupID string, sinceVersion int64) ([]*model.GroupKeyEvent, error)
}
