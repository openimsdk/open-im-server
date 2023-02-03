package unrelation

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
}
