package relation

import (
	"context"
	"time"

	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/ormutil"
	"gorm.io/gorm"

	relationtb "github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
)

type LogGorm struct {
	db *gorm.DB
}

func (l *LogGorm) Create(ctx context.Context, log []*relationtb.Log) error {
	return errs.Wrap(l.db.WithContext(ctx).Create(log).Error)
}

func (l *LogGorm) Search(ctx context.Context, keyword string, start time.Time, end time.Time, pageNumber int32, showNumber int32) (uint32, []*relationtb.Log, error) {
	db := l.db.WithContext(ctx).Where("create_time >= ?", start)
	if end.UnixMilli() != 0 {
		db = l.db.WithContext(ctx).Where("create_time <= ?", end)
	}
	return ormutil.GormSearch[relationtb.Log](db, []string{"user_id"}, keyword, pageNumber, showNumber)
}

func (l *LogGorm) Delete(ctx context.Context, logIDs []string, userID string) error {
	if userID == "" {
		return errs.Wrap(l.db.WithContext(ctx).Where("log_id in ?", logIDs).Delete(&relationtb.Log{}).Error)
	}
	return errs.Wrap(l.db.WithContext(ctx).Where("log_id in ? and user_id=?", logIDs, userID).Delete(&relationtb.Log{}).Error)
}

func (l *LogGorm) Get(ctx context.Context, logIDs []string, userID string) ([]*relationtb.Log, error) {
	var logs []*relationtb.Log
	if userID == "" {
		return logs, errs.Wrap(l.db.WithContext(ctx).Where("log_id in ?", logIDs).Find(&logs).Error)
	}
	return logs, errs.Wrap(l.db.WithContext(ctx).Where("log_id in ? and user_id=?", logIDs, userID).Find(&logs).Error)
}

func NewLogGorm(db *gorm.DB) relationtb.LogInterface {
	db.AutoMigrate(&relationtb.Log{})
	return &LogGorm{db: db}
}
