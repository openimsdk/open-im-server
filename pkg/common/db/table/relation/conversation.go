package relation

import "context"

const (
	conversationModelTableName = "conversations"
)

type ConversationModel struct {
	OwnerUserID           string `gorm:"column:owner_user_id;primary_key;type:char(128)" json:"OwnerUserID"`
	ConversationID        string `gorm:"column:conversation_id;primary_key;type:char(128)" json:"conversationID"`
	ConversationType      int32  `gorm:"column:conversation_type" json:"conversationType"`
	UserID                string `gorm:"column:user_id;type:char(64)" json:"userID"`
	GroupID               string `gorm:"column:group_id;type:char(128)" json:"groupID"`
	RecvMsgOpt            int32  `gorm:"column:recv_msg_opt" json:"recvMsgOpt"`
	UnreadCount           int32  `gorm:"column:unread_count" json:"unreadCount"`
	DraftTextTime         int64  `gorm:"column:draft_text_time" json:"draftTextTime"`
	IsPinned              bool   `gorm:"column:is_pinned" json:"isPinned"`
	IsPrivateChat         bool   `gorm:"column:is_private_chat" json:"isPrivateChat"`
	BurnDuration          int32  `gorm:"column:burn_duration;default:30" json:"burnDuration"`
	GroupAtType           int32  `gorm:"column:group_at_type" json:"groupAtType"`
	IsNotInGroup          bool   `gorm:"column:is_not_in_group" json:"isNotInGroup"`
	UpdateUnreadCountTime int64  `gorm:"column:update_unread_count_time" json:"updateUnreadCountTime"`
	AttachedInfo          string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
	Ex                    string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
	MaxSeq                int64  `gorm:"column:max_seq" json:"maxSeq"`
	MinSeq                int64  `gorm:"column:min_seq" json:"minSeq"`
	IsReadSeq             int64  `gorm:"column:is_read_seq" json:"isReadSeq"`
}

func (ConversationModel) TableName() string {
	return conversationModelTableName
}

type ConversationModelInterface interface {
	Create(ctx context.Context, conversations []*ConversationModel) (err error)
	Delete(ctx context.Context, groupIDs []string) (err error)
	UpdateByMap(ctx context.Context, userIDs []string, conversationID string, args map[string]interface{}) (rows int64, err error)
	Update(ctx context.Context, conversation *ConversationModel) (err error)
	Find(ctx context.Context, ownerUserID string, conversationIDs []string) (conversations []*ConversationModel, err error)
	FindUserID(ctx context.Context, userIDs []string, conversationIDs []string) ([]string, error)
	FindUserIDAllConversationID(ctx context.Context, userID string) ([]string, error)
	Take(ctx context.Context, userID, conversationID string) (conversation *ConversationModel, err error)
	FindConversationID(ctx context.Context, userID string, conversationIDs []string) (existConversationID []string, err error)
	FindUserIDAllConversations(ctx context.Context, userID string) (conversations []*ConversationModel, err error)
	FindRecvMsgNotNotifyUserIDs(ctx context.Context, groupID string) ([]string, error)
	GetUserRecvMsgOpt(ctx context.Context, ownerUserID, conversationID string) (opt int, err error)
	FindSuperGroupRecvMsgNotNotifyUserIDs(ctx context.Context, groupID string) ([]string, error)
	NewTx(tx any) ConversationModelInterface
}
