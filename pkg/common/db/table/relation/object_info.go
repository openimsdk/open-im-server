package relation

import (
	"context"
	"time"
)

const (
	ObjectInfoModelTableName = "object_info"
)

type ObjectInfoModel struct {
	Name           string     `gorm:"column:name;primary_key"`
	Hash           string     `gorm:"column:hash"`
	ExpirationTime *time.Time `gorm:"column:expiration_time"`
	CreateTime     time.Time  `gorm:"column:create_time"`
}

func (ObjectInfoModel) TableName() string {
	return ObjectInfoModelTableName
}

type ObjectInfoModelInterface interface {
	NewTx(tx any) ObjectInfoModelInterface
	SetObject(ctx context.Context, obj *ObjectInfoModel) error
	Take(ctx context.Context, name string) (*ObjectInfoModel, error)
}
