package relation

import (
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/utils"
	"context"
	"fmt"
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

func NewUserDB(db *gorm.DB) *User {
	var user User
	user.DB = db
	return &user
}

func (u *User) Create(ctx context.Context, users []*User) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "users", users)
	}()
	err = utils.Wrap(u.DB.Create(&users).Error, "")
	return err
}

func (u *User) UpdateByMap(ctx context.Context, userID string, args map[string]interface{}) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userID", userID, "args", args)
	}()
	return utils.Wrap(u.DB.Where("user_id = ?", userID).Updates(args).Error, "")
}

func (u *User) Update(ctx context.Context, users []*User) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "users", users)
	}()
	return utils.Wrap(u.DB.Updates(&users).Error, "")
}

func (u *User) Find(ctx context.Context, userIDs []string) (users []*User, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userIDs", userIDs, "users", users)
	}()
	err = utils.Wrap(u.DB.Where("user_id in (?)", userIDs).Find(&users).Error, "")
	return users, err
}

func (u *User) Take(ctx context.Context, userID string) (user *User, err error) {
	user = &User{}
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userID", userID, "user", *user)
	}()
	err = utils.Wrap(u.DB.Where("user_id = ?", userID).Take(&user).Error, "")
	return user, err
}

func (u *User) GetByName(ctx context.Context, userName string, showNumber, pageNumber int32) (users []*User, count int64, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userName", userName, "showNumber", showNumber, "pageNumber", pageNumber, "users", users, "count", count)
	}()
	err = u.DB.Where(" name like ?", fmt.Sprintf("%%%s%%", userName)).Limit(int(showNumber)).Offset(int(showNumber * pageNumber)).Find(&users).Error
	if err != nil {
		return nil, 0, utils.Wrap(err, "")
	}
	return users, count, utils.Wrap(u.DB.Where(" name like ? ", fmt.Sprintf("%%%s%%", userName)).Count(&count).Error, "")
}

func (u *User) GetByNameAndID(ctx context.Context, content string, showNumber, pageNumber int32) (users []*User, count int64, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "content", content, "showNumber", showNumber, "pageNumber", pageNumber, "users", users)
	}()
	db := u.DB.Where(" name like ? or user_id = ? ", fmt.Sprintf("%%%s%%", content), content)
	if err := db.Count(&count).Error; err != nil {
		return nil, 0, utils.Wrap(err, "")
	}
	err = utils.Wrap(db.Limit(int(showNumber)).Offset(int(showNumber*pageNumber)).Find(&users).Error, "")
	return
}

func (u *User) Get(ctx context.Context, showNumber, pageNumber int32) (users []*User, count int64, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "showNumber", showNumber, "pageNumber", pageNumber, "users", users, "count", count)
	}()
	err = u.DB.Model(u).Count(&count).Error
	if err != nil {
		return nil, 0, utils.Wrap(err, "")
	}
	err = utils.Wrap(u.DB.Limit(int(showNumber)).Offset(int(pageNumber*showNumber)).Find(&users).Error, "")
	return
}
