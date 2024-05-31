package database

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"time"
)

const (
	FirstVersion         = 1
	DefaultDeleteVersion = 0
)

type VersionLog interface {
	IncrVersion(ctx context.Context, dId string, eIds []string, deleted bool) error
	FindChangeLog(ctx context.Context, dId string, version uint, limit int) (*model.VersionLog, error)
	DeleteAfterUnchangedLog(ctx context.Context, deadline time.Time) error
}
