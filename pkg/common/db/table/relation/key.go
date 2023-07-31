package relation

import "context"

const (
	KeyModelTableName = "keys"
)

type KeyModel struct {
	CId   string `gorm:"column:cid;primary_key;size:64" json:"cid"  binding:"required"`
	Key   string `gorm:"column:key;size:64" json:"key"  binding:"required"`
	CType int32  `gorm:"column:cType" json:"cType"  binding:"required"`
}

func (KeyModel) TableName() string {
	return KeyModelTableName
}

type KeyModelInterface interface {
	InstallKey(ctx context.Context, key KeyModel) (err error)
	GetKey(ctx context.Context, cid string) (key KeyModel, err error)
}
