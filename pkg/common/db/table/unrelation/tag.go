package unrelation

import "context"

const (
	CTag     = "tag"
	CSendLog = "send_log"
)

type TagModel struct {
	UserID   string   `bson:"user_id"`
	TagID    string   `bson:"tag_id"`
	TagName  string   `bson:"tag_name"`
	UserList []string `bson:"user_list"`
}

func (TagModel) TableName() string {
	return CTag
}

type TagSendLogModel struct {
	UserList         []CommonUserModel `bson:"tag_list"`
	SendID           string            `bson:"send_id"`
	SenderPlatformID int32             `bson:"sender_platform_id"`
	Content          string            `bson:"content"`
	SendTime         int64             `bson:"send_time"`
}

func (TagSendLogModel) TableName() string {
	return CSendLog
}

type TagModelInterface interface {
	GetUserTags(ctx context.Context, userID string) ([]TagModel, error)
	CreateTag(ctx context.Context, userID, tagName string, userList []string) error
	GetTagByID(ctx context.Context, userID, tagID string) (TagModel, error)
	DeleteTag(ctx context.Context, userID, tagID string) error
	SetTag(ctx context.Context, userID, tagID, newName string, increaseUserIDList []string, reduceUserIDList []string) error
	GetUserIDListByTagID(ctx context.Context, userID, tagID string) ([]string, error)
	SaveTagSendLog(ctx context.Context, tagSendLog *TagSendLogModel) error
	GetTagSendLogs(ctx context.Context, userID string, showNumber, pageNumber int32) ([]TagSendLogModel, error)
}
