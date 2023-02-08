package unrelation

import "context"

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

type WorkMomentModelInterface interface {
	CreateOneWorkMoment(ctx context.Context, workMoment *WorkMoment) error
	DeleteOneWorkMoment(ctx context.Context, workMomentID string) error
	DeleteComment(ctx context.Context, workMomentID, contentID, opUserID string) error
	GetWorkMomentByID(ctx context.Context, workMomentID string) (*WorkMoment, error)
	LikeOneWorkMoment(ctx context.Context, likeUserID, userName, workMomentID string) (*WorkMoment, bool, error)
	CommentOneWorkMoment(ctx context.Context, comment *Comment, workMomentID string) (*WorkMoment, error)
	GetUserSelfWorkMoments(ctx context.Context, userID string, showNumber, pageNumber int32) ([]*WorkMoment, error)
	GetUserWorkMoments(ctx context.Context, opUserID, userID string, showNumber, pageNumber int32, friendIDList []string) ([]*WorkMoment, error)
	GetUserFriendWorkMoments(ctx context.Context, showNumber, pageNumber int32, userID string, friendIDList []string) ([]*WorkMoment, error)
}
