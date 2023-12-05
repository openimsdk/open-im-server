package relation

import (
	"context"
	"time"

	"github.com/OpenIMSDK/tools/pagination"
)

type LogModel struct {
	LogID      string    `bson:"log_id"`
	Platform   string    `bson:"platform"`
	UserID     string    `bson:"user_id"`
	CreateTime time.Time `bson:"create_time"`
	Url        string    `bson:"url"`
	FileName   string    `bson:"file_name"`
	SystemType string    `bson:"system_type"`
	Version    string    `bson:"version"`
	Ex         string    `bson:"ex"`
}

type LogInterface interface {
	Create(ctx context.Context, log []*LogModel) error
	Search(ctx context.Context, keyword string, start time.Time, end time.Time, pagination pagination.Pagination) (int64, []*LogModel, error)
	Delete(ctx context.Context, logID []string, userID string) error
	Get(ctx context.Context, logIDs []string, userID string) ([]*LogModel, error)
}
