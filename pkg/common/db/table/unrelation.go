package table

type SuperGroup struct {
	GroupID      string   `bson:"group_id" json:"groupID"`
	MemberIDList []string `bson:"member_id_list" json:"memberIDList"`
}

type UserToSuperGroup struct {
	UserID      string   `bson:"user_id" json:"userID"`
	GroupIDList []string `bson:"group_id_list" json:"groupIDList"`
}

type Tag struct {
	UserID   string   `bson:"user_id"`
	TagID    string   `bson:"tag_id"`
	TagName  string   `bson:"tag_name"`
	UserList []string `bson:"user_list"`
}

type CommonUser struct {
	UserID   string `bson:"user_id"`
	UserName string `bson:"user_name"`
}

type TagSendLog struct {
	UserList         []CommonUser `bson:"tag_list"`
	SendID           string       `bson:"send_id"`
	SenderPlatformID int32        `bson:"sender_platform_id"`
	Content          string       `bson:"content"`
	SendTime         int64        `bson:"send_time"`
}

type WorkMoment struct {
	WorkMomentID         string        `bson:"work_moment_id"`
	UserID               string        `bson:"user_id"`
	UserName             string        `bson:"user_name"`
	FaceURL              string        `bson:"face_url"`
	Content              string        `bson:"content"`
	LikeUserList         []*CommonUser `bson:"like_user_list"`
	AtUserList           []*CommonUser `bson:"at_user_list"`
	PermissionUserList   []*CommonUser `bson:"permission_user_list"`
	Comments             []*CommonUser `bson:"comments"`
	PermissionUserIDList []string      `bson:"permission_user_id_list"`
	Permission           int32         `bson:"permission"`
	CreateTime           int32         `bson:"create_time"`
}

type Comment struct {
	UserID        string `bson:"user_id" json:"user_id"`
	UserName      string `bson:"user_name" json:"user_name"`
	ReplyUserID   string `bson:"reply_user_id" json:"reply_user_id"`
	ReplyUserName string `bson:"reply_user_name" json:"reply_user_name"`
	ContentID     string `bson:"content_id" json:"content_id"`
	Content       string `bson:"content" json:"content"`
	CreateTime    int32  `bson:"create_time" json:"create_time"`
}
