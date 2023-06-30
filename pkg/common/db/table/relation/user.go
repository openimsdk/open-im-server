package relation

import (
	"context"
	"time"
)

const (
	UserModelTableName = "users"
)

type UserModel struct {
	UserID           string    `gorm:"column:user_id;primary_key;size:64"`
	Nickname         string    `gorm:"column:name;size:255"`
	FaceURL          string    `gorm:"column:face_url;size:255"`
	Ex               string    `gorm:"column:ex;size:1024"`
	CreateTime       time.Time `gorm:"column:create_time;index:create_time; autoCreateTime"`
	AppMangerLevel   int32     `gorm:"column:app_manger_level;default:18"`
	GlobalRecvMsgOpt int32     `gorm:"column:global_recv_msg_opt"`
}

func (u *UserModel) GetNickname() string {
	return u.Nickname
}

func (u *UserModel) GetFaceURL() string {
	return u.FaceURL
}

func (u *UserModel) GetUserID() string {
	return u.UserID
}

func (u *UserModel) GetEx() string {
	return u.Ex
}

func (UserModel) TableName() string {
	return UserModelTableName
}

type UserModelInterface interface {
	Create(ctx context.Context, users []*UserModel) (err error)
	UpdateByMap(ctx context.Context, userID string, args map[string]interface{}) (err error)
	Update(ctx context.Context, user *UserModel) (err error)
	// 获取指定用户信息  不存在，也不返回错误
	Find(ctx context.Context, userIDs []string) (users []*UserModel, err error)
	// 获取某个用户信息  不存在，则返回错误
	Take(ctx context.Context, userID string) (user *UserModel, err error)
	// 获取用户信息 不存在，不返回错误
	Page(ctx context.Context, pageNumber, showNumber int32) (users []*UserModel, count int64, err error)
	GetAllUserID(ctx context.Context) (userIDs []string, err error)
	GetUserGlobalRecvMsgOpt(ctx context.Context, userID string) (opt int, err error)
	// 获取用户总数
	CountTotal(ctx context.Context) (count int64, err error)
	// 获取范围内用户增量
	CountRangeEverydayTotal(ctx context.Context, start time.Time, end time.Time) (map[string]int64, error)
}
