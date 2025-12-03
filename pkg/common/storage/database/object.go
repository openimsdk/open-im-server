package database

import (
	"context"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

type ObjectInfo interface {
	SetObject(ctx context.Context, obj *model.Object) error
	Take(ctx context.Context, engine string, name string) (*model.Object, error)
	Delete(ctx context.Context, engine string, name []string) error
	FindExpirationObject(ctx context.Context, engine string, expiration time.Time, needDelType []string, count int64) ([]*model.Object, error)
	GetKeyCount(ctx context.Context, engine string, key string) (int64, error)

	GetEngineCount(ctx context.Context, engine string) (int64, error)
	GetEngineInfo(ctx context.Context, engine string, limit int, skip int) ([]*model.Object, error)
	UpdateEngine(ctx context.Context, oldEngine, oldName string, newEngine string) error
}
