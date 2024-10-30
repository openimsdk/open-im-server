package database

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/db/pagination"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Application interface {
	LatestVersion(ctx context.Context, platform string, hot bool) (*model.Application, error)
	AddVersion(ctx context.Context, val *model.Application) error
	UpdateVersion(ctx context.Context, id primitive.ObjectID, update map[string]any) error
	DeleteVersion(ctx context.Context, id []primitive.ObjectID) error
	PageVersion(ctx context.Context, platforms []string, page pagination.Pagination) (int64, []*model.Application, error)
	FindPlatform(ctx context.Context, id []primitive.ObjectID) ([]string, error)
}
