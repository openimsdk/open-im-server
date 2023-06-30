package relation

import (
	"context"
	"time"
)

const (
	ObjectHashModelTableName = "object_hash"
)

type ObjectHashModel struct {
	Hash       string    `gorm:"column:hash;primary_key;size:32"`
	Engine     string    `gorm:"column:engine;primary_key;size:16"`
	Size       int64     `gorm:"column:size"`
	Bucket     string    `gorm:"column:bucket"`
	Name       string    `gorm:"column:name"`
	CreateTime time.Time `gorm:"column:create_time"`
}

func (ObjectHashModel) TableName() string {
	return ObjectHashModelTableName
}

type ObjectHashModelInterface interface {
	NewTx(tx any) ObjectHashModelInterface
	Take(ctx context.Context, hash string, engine string) (*ObjectHashModel, error)
	Create(ctx context.Context, h []*ObjectHashModel) error
	DeleteNoCitation(ctx context.Context, engine string, num int) (list []*ObjectHashModel, err error)
}
