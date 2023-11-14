package relation

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/pagination"
	"time"
)

//type Log struct {
//	LogID      string    `gorm:"column:log_id;primary_key;type:char(64)"`
//	Platform   string    `gorm:"column:platform;type:varchar(32)"`
//	UserID     string    `gorm:"column:user_id;type:char(64)"`
//	CreateTime time.Time `gorm:"index:,sort:desc"`
//	Url        string    `gorm:"column:url;type varchar(255)"`
//	FileName   string    `gorm:"column:filename;type varchar(255)"`
//	SystemType string    `gorm:"column:system_type;type varchar(255)"`
//	Version    string    `gorm:"column:version;type varchar(255)"`
//	Ex         string    `gorm:"column:ex;type varchar(255)"`
//}

type Log struct {
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
	Create(ctx context.Context, log []*Log) error
	Search(ctx context.Context, keyword string, start time.Time, end time.Time, pagination pagination.Pagination) (int64, []*Log, error)
	Delete(ctx context.Context, logID []string, userID string) error
	Get(ctx context.Context, logIDs []string, userID string) ([]*Log, error)
}
