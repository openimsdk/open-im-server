package relation

import "context"

const (
	AesKeyModelTableName = "aes_keys"
)

type AesKeyModel struct {
	KeyConversationsID string `gorm:"column:key_conversations_id;primary_key;size:64" json:"keyConversationsID"`
	Key                string `gorm:"column:key" json:"key"`
	ConversationType   int32  `gorm:"column:conversation_type" json:"conversationType"`
	OwnerUserID        string `gorm:"column:owner_user_id;size:64" json:"ownerUserID"`
	FriendUserID       string `gorm:"column:friend_user_id;size:64" json:"friendUserID"`
	GroupID            string `gorm:"column:group_id;size:64" json:"groupID"`
}

func (AesKeyModel) TableName() string {
	return AesKeyModelTableName
}

type AesKeyModelInterface interface {
	Installs(ctx context.Context, friends []*AesKeyModel) (err error)
	GetAesKey(tx context.Context, KeyConversationsID string) (key *AesKeyModel, err error)
	GetAllAesKey(tx context.Context, UserID string) (key []*AesKeyModel, err error)
}
