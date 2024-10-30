package cache

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

type ApplicationCache interface {
	LatestVersion(ctx context.Context, platform string) (*model.Application, error)
	DeleteCache(ctx context.Context, platforms []string) error
}
