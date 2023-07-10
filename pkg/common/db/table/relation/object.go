package relation

import (
	"context"
	"time"
)

const (
	ObjectInfoModelTableName = "object"
)

type ObjectModel struct {
	Name        string    `gorm:"column:name;primary_key"`
	UserID      string    `gorm:"column:user_id"`
	Hash        string    `gorm:"column:hash"`
	Key         string    `gorm:"column:key"`
	Size        int64     `gorm:"column:size"`
	ContentType string    `gorm:"column:content_type"`
	Cause       string    `gorm:"column:cause"`
	CreateTime  time.Time `gorm:"column:create_time"`
}

func (ObjectModel) TableName() string {
	return ObjectInfoModelTableName
}

type ObjectInfoModelInterface interface {
	NewTx(tx any) ObjectInfoModelInterface
	SetObject(ctx context.Context, obj *ObjectModel) error
	Take(ctx context.Context, name string) (*ObjectModel, error)
}
