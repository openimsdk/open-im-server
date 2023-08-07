package relation

import (
	"context"
)

const (
	AesKeyModelTableName = "aes-keys"
)

type AesKeyModel struct {
	UserID           string `gorm:"column:user_id;size:64"`
	ConversationID   string `gorm:"column:conversation_id;uniqueIndex:idx_key"   json:"conversationID"`
	AesKey           string `gorm:"column:conversation_id"   json:"aesKey"`
	ConversationType int32  `gorm:"column:conversation_type;uniqueIndex:idx_key" json:"conversationType"`
}

func (AesKeyModel) TableName() string {
	return AesKeyModelTableName
}

type AesKeyModelInterface interface {
	Install(ctx context.Context, aesKey AesKeyModel) (err error)
	GetAesKey(ctx context.Context, userId, cid string, cType int32) (aesKey *AesKeyModel, err error)
	GetAllAesKey(ctx context.Context, userId string) (aesKey []*AesKeyModel, err error)
}
