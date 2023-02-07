package relation

import (
	"Open_IM/pkg/common/db/table/relation"
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/utils"
	"context"
	"fmt"
	"gorm.io/gorm"
)

type UserGorm struct {
	DB *gorm.DB
}

func NewUserGorm(db *gorm.DB) *UserGorm {
	var user UserGorm
	user.DB = db
	return &user
}

// 插入多条
func (u *UserGorm) Create(ctx context.Context, users []*relation.UserModel, tx ...*gorm.DB) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "users", users)
	}()
	return utils.Wrap(getDBConn(u.DB, tx).Create(&users).Error, "")
}

// 更新用户信息 零值
func (u *UserGorm) UpdateByMap(ctx context.Context, userID string, args map[string]interface{}, tx ...*gorm.DB) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userID", userID, "args", args)
	}()
	return utils.Wrap(getDBConn(u.DB, tx).Model(&relation.UserModel{}).Where("user_id = ?", userID).Updates(args).Error, "")
}

// 更新多个用户信息 非零值
func (u *UserGorm) Update(ctx context.Context, users []*relation.UserModel, tx ...*gorm.DB) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "users", users)
	}()
	return utils.Wrap(getDBConn(u.DB, tx).Updates(&users).Error, "")
}

// 获取指定用户信息  不存在，也不返回错误
func (u *UserGorm) Find(ctx context.Context, userIDs []string, tx ...*gorm.DB) (users []*relation.UserModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userIDs", userIDs, "users", users)
	}()
	err = utils.Wrap(getDBConn(u.DB, tx).Where("user_id in (?)", userIDs).Find(&users).Error, "")
	return users, err
}

// 获取某个用户信息  不存在，则返回错误
func (u *UserGorm) Take(ctx context.Context, userID string, tx ...*gorm.DB) (user *relation.UserModel, err error) {
	user = &relation.UserModel{}
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userID", userID, "user", *user)
	}()
	err = utils.Wrap(getDBConn(u.DB, tx).Where("user_id = ?", userID).Take(&user).Error, "")
	return user, err
}

// 通过名字查找用户 不存在，不返回错误
func (u *UserGorm) GetByName(ctx context.Context, userName string, pageNumber, showNumber int32, tx ...*gorm.DB) (users []*relation.UserModel, count int64, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userName", userName, "pageNumber", pageNumber, "showNumber", showNumber, "users", users, "count", count)
	}()
	err = utils.Wrap(getDBConn(u.DB, tx).Model(&relation.UserModel{}).Where(" name like ?", fmt.Sprintf("%%%s%%", userName)).Count(&count).Error, "")
	if err != nil {
		return
	}
	err = utils.Wrap(getDBConn(u.DB, tx).Model(&relation.UserModel{}).Where(" name like ?", fmt.Sprintf("%%%s%%", userName)).Limit(int(showNumber)).Offset(int(showNumber*pageNumber)).Find(&users).Error, "")
	return
}

// 通过名字或userID查找用户 不存在，不返回错误
func (u *UserGorm) GetByNameAndID(ctx context.Context, content string, pageNumber, showNumber int32, tx ...*gorm.DB) (users []*relation.UserModel, count int64, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "content", content, "pageNumber", pageNumber, "showNumber", showNumber, "users", users, "count", count)
	}()
	db := getDBConn(u.DB, tx).Model(&relation.UserModel{}).Where(" name like ? or user_id = ? ", fmt.Sprintf("%%%s%%", content), content)
	if err = db.Count(&count).Error; err != nil {
		return
	}
	err = utils.Wrap(db.Limit(int(showNumber)).Offset(int(showNumber*pageNumber)).Find(&users).Error, "")
	return
}

// 获取用户信息 不存在，不返回错误
func (u *UserGorm) Get(ctx context.Context, pageNumber, showNumber int32, tx ...*gorm.DB) (users []*relation.UserModel, count int64, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "pageNumber", pageNumber, "showNumber", showNumber, "users", users, "count", count)
	}()
	err = utils.Wrap(getDBConn(u.DB, tx).Model(&relation.UserModel{}).Count(&count).Error, "")
	if err != nil {
		return
	}
	err = utils.Wrap(getDBConn(u.DB, tx).Limit(int(showNumber)).Offset(int(pageNumber*showNumber)).Find(&users).Error, "")
	return
}
