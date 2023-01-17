package mysql

import (
	"Open_IM/pkg/common/trace_log"
	"Open_IM/pkg/utils"
	"context"
	"gorm.io/gorm"
	"time"
)

type User struct {
	UserID           string    `gorm:"column:user_id;primary_key;size:64"`
	Nickname         string    `gorm:"column:name;size:255"`
	FaceURL          string    `gorm:"column:face_url;size:255"`
	Gender           int32     `gorm:"column:gender"`
	PhoneNumber      string    `gorm:"column:phone_number;size:32"`
	Birth            time.Time `gorm:"column:birth"`
	Email            string    `gorm:"column:email;size:64"`
	Ex               string    `gorm:"column:ex;size:1024"`
	CreateTime       time.Time `gorm:"column:create_time;index:create_time"`
	AppMangerLevel   int32     `gorm:"column:app_manger_level"`
	GlobalRecvMsgOpt int32     `gorm:"column:global_recv_msg_opt"`

	status int32    `gorm:"column:status"`
	DB     *gorm.DB `gorm:"-" json:"-"`
}

func NewUserDB() *User {
	var user User
	user.DB = initMysqlDB(&user)
	return &user
}

func (u *User) Create(ctx context.Context, users []*User) (err error) {
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetFuncName(1), err, "users", users)
	}()
	err = utils.Wrap(u.DB.Create(&users).Error, "")
	return err
}

func (u *User) UpdateByMap(ctx context.Context, userID string, args map[string]interface{}) (err error) {
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userID", userID, "args", args)
	}()
	return utils.Wrap(u.DB.Where("user_id = ?", userID).Updates(args).Error, "")
}

func (u *User) Update(ctx context.Context, users []*User) (err error) {
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetFuncName(1), err, "users", users)
	}()
	return utils.Wrap(u.DB.Updates(&users).Error, "")
}

func (u *User) Find(ctx context.Context, userIDs []string) (users []*User, err error) {
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userIDs", userIDs, "users", users)
	}()
	err = utils.Wrap(u.DB.Where("user_id in (?)", userIDs).Find(&users).Error, "")
	return users, err
}

func (u *User) Take(ctx context.Context, userID string) (user *User, err error) {
	user = &User{}
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userID", userID, "user", *user)
	}()
	err = utils.Wrap(u.DB.Where("user_id = ?", userID).Take(&user).Error, "")
	return user, err
}
