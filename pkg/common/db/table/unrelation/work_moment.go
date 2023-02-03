package unrelation

const (
	CWorkMoment = "work_moment"
)

type WorkMoment struct {
	WorkMomentID         string             `bson:"work_moment_id"`
	UserID               string             `bson:"user_id"`
	UserName             string             `bson:"user_name"`
	FaceURL              string             `bson:"face_url"`
	Content              string             `bson:"content"`
	LikeUserList         []*CommonUserModel `bson:"like_user_list"`
	AtUserList           []*CommonUserModel `bson:"at_user_list"`
	PermissionUserList   []*CommonUserModel `bson:"permission_user_list"`
	Comments             []*CommonUserModel `bson:"comments"`
	PermissionUserIDList []string           `bson:"permission_user_id_list"`
	Permission           int32              `bson:"permission"`
	CreateTime           int32              `bson:"create_time"`
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

func (WorkMoment) TableName() string {
	return CWorkMoment
}
