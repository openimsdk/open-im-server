package relation

import (
	"OpenIM/pkg/common/db/table/relation"
	"OpenIM/pkg/common/tracelog"
	"OpenIM/pkg/utils"
	"context"
	"gorm.io/gorm"
)

type UserGorm struct {
	DB *gorm.DB
}

func NewUserGorm(DB *gorm.DB) relation.UserModelInterface {
	return &UserGorm{DB: DB.Model(&relation.UserModel{})}
}

// 插入多条
func (u *UserGorm) Create(ctx context.Context, users []*relation.UserModel) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "users", users)
	}()
	return utils.Wrap(u.DB.Create(&users).Error, "")
}

// 更新用户信息 零值
func (u *UserGorm) UpdateByMap(ctx context.Context, userID string, args map[string]interface{}) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userID", userID, "args", args)
	}()
	return utils.Wrap(u.DB.Where("user_id = ?", userID).Updates(args).Error, "")
}

// 更新多个用户信息 非零值
func (u *UserGorm) Update(ctx context.Context, users []*relation.UserModel) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "users", users)
	}()
	return utils.Wrap(u.DB.Updates(&users).Error, "")
}

// 获取指定用户信息  不存在，也不返回错误
func (u *UserGorm) Find(ctx context.Context, userIDs []string) (users []*relation.UserModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userIDs", userIDs, "users", users)
	}()
	err = utils.Wrap(u.DB.Where("user_id in (?)", userIDs).Find(&users).Error, "")
	return users, err
}

// 获取某个用户信息  不存在，则返回错误
func (u *UserGorm) Take(ctx context.Context, userID string) (user *relation.UserModel, err error) {
	user = &relation.UserModel{}
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userID", userID, "user", *user)
	}()
	err = utils.Wrap(u.DB.Where("user_id = ?", userID).Take(&user).Error, "")
	return user, err
}

// 获取用户信息 不存在，不返回错误
func (u *UserGorm) Page(ctx context.Context, pageNumber, showNumber int32) (users []*relation.UserModel, count int64, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "pageNumber", pageNumber, "showNumber", showNumber, "users", users, "count", count)
	}()
	err = utils.Wrap(u.DB.Count(&count).Error, "")
	if err != nil {
		return
	}
	err = utils.Wrap(u.DB.Limit(int(showNumber)).Offset(int(pageNumber*showNumber)).Find(&users).Error, "")
	return
}

// 获取所有用户ID
func (u *UserGorm) GetAllUserID(ctx context.Context) (userIDs []string, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userIDs", userIDs)
	}()
	err = u.DB.Pluck("user_id", &userIDs).Error
	return userIDs, err
}
