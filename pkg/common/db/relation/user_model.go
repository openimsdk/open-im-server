package relation

import (
	"Open_IM/pkg/common/db/table/relation"
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/utils"
	"context"
	"gorm.io/gorm"
)

type User interface {
	Create(ctx context.Context, users []*relation.UserModel, tx ...any) (err error)
	UpdateByMap(ctx context.Context, userID string, args map[string]interface{}, tx ...any) (err error)
	Update(ctx context.Context, users []*relation.UserModel, tx ...any) (err error)
	// 获取指定用户信息  不存在，也不返回错误
	Find(ctx context.Context, userIDs []string, tx ...any) (users []*relation.UserModel, err error)
	// 获取某个用户信息  不存在，则返回错误
	Take(ctx context.Context, userID string, tx ...any) (user *relation.UserModel, err error)
	// 获取用户信息 不存在，不返回错误
	Page(ctx context.Context, pageNumber, showNumber int32, tx ...any) (users []*relation.UserModel, count int64, err error)
	GetAllUserID(ctx context.Context) (userIDs []string, err error)
}

type UserGorm struct {
	DB *gorm.DB
}

func NewUserGorm(DB *gorm.DB) *UserGorm {
	return &UserGorm{DB: DB}
}

// 插入多条
func (u *UserGorm) Create(ctx context.Context, users []*relation.UserModel, tx ...any) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "users", users)
	}()
	return utils.Wrap(getDBConn(u.DB, tx).Create(&users).Error, "")
}

// 更新用户信息 零值
func (u *UserGorm) UpdateByMap(ctx context.Context, userID string, args map[string]interface{}, tx ...any) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userID", userID, "args", args)
	}()
	return utils.Wrap(getDBConn(u.DB, tx).Model(&relation.UserModel{}).Where("user_id = ?", userID).Updates(args).Error, "")
}

// 更新多个用户信息 非零值
func (u *UserGorm) Update(ctx context.Context, users []*relation.UserModel, tx ...any) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "users", users)
	}()
	return utils.Wrap(getDBConn(u.DB, tx).Updates(&users).Error, "")
}

// 获取指定用户信息  不存在，也不返回错误
func (u *UserGorm) Find(ctx context.Context, userIDs []string, tx ...any) (users []*relation.UserModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userIDs", userIDs, "users", users)
	}()
	err = utils.Wrap(getDBConn(u.DB, tx).Where("user_id in (?)", userIDs).Find(&users).Error, "")
	return users, err
}

// 获取某个用户信息  不存在，则返回错误
func (u *UserGorm) Take(ctx context.Context, userID string, tx ...any) (user *relation.UserModel, err error) {
	user = &relation.UserModel{}
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userID", userID, "user", *user)
	}()
	err = utils.Wrap(getDBConn(u.DB, tx).Where("user_id = ?", userID).Take(&user).Error, "")
	return user, err
}

// 获取用户信息 不存在，不返回错误
func (u *UserGorm) Page(ctx context.Context, pageNumber, showNumber int32, tx ...any) (users []*relation.UserModel, count int64, err error) {
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

// 获取所有用户ID
func (u *UserGorm) GetAllUserID(ctx context.Context) (userIDs []string, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userIDs", userIDs)
	}()

	err = u.DB.Pluck("user_id", &userIDs).Error
	return userIDs, err
}
