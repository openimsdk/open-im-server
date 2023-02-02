package relation

import (
	"Open_IM/pkg/common/db/table"
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

func (u *UserGorm) Create(ctx context.Context, users []*table.UserModel, tx ...*gorm.DB) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "users", users)
	}()
	return utils.Wrap(getDBConn(u.DB, tx).Model(&table.UserModel{}).Create(&users).Error, "")
}

func (u *UserGorm) UpdateByMap(ctx context.Context, userID string, args map[string]interface{}) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userID", userID, "args", args)
	}()
	return utils.Wrap(u.DB.Model(&table.UserModel{}).Where("user_id = ?", userID).Updates(args).Error, "")
}

func (u *UserGorm) Update(ctx context.Context, users []*table.UserModel) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "users", users)
	}()
	return utils.Wrap(u.DB.Model(&table.UserModel{}).Updates(&users).Error, "")
}

func (u *UserGorm) Find(ctx context.Context, userIDs []string) (users []*table.UserModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userIDs", userIDs, "users", users)
	}()
	err = utils.Wrap(u.DB.Model(&table.UserModel{}).Where("user_id in (?)", userIDs).Find(&users).Error, "")
	return users, err
}

func (u *UserGorm) Take(ctx context.Context, userID string) (user *table.UserModel, err error) {
	user = &table.UserModel{}
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userID", userID, "user", *user)
	}()
	err = utils.Wrap(u.DB.Model(&table.UserModel{}).Where("user_id = ?", userID).Take(&user).Error, "")
	return user, err
}

func (u *UserGorm) GetByName(ctx context.Context, userName string, showNumber, pageNumber int32) (users []*table.UserModel, count int64, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userName", userName, "showNumber", showNumber, "pageNumber", pageNumber, "users", users, "count", count)
	}()
	err = u.DB.Model(&table.UserModel{}).Where(" name like ?", fmt.Sprintf("%%%s%%", userName)).Limit(int(showNumber)).Offset(int(showNumber * pageNumber)).Find(&users).Error
	if err != nil {
		return nil, 0, utils.Wrap(err, "")
	}
	return users, count, utils.Wrap(u.DB.Model(&table.UserModel{}).Where(" name like ? ", fmt.Sprintf("%%%s%%", userName)).Count(&count).Error, "")
}

func (u *UserGorm) GetByNameAndID(ctx context.Context, content string, showNumber, pageNumber int32) (users []*table.UserModel, count int64, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "content", content, "showNumber", showNumber, "pageNumber", pageNumber, "users", users)
	}()
	db := u.DB.Model(&table.UserModel{}).Where(" name like ? or user_id = ? ", fmt.Sprintf("%%%s%%", content), content)
	if err := db.Model(&table.UserModel{}).Count(&count).Error; err != nil {
		return nil, 0, utils.Wrap(err, "")
	}
	err = utils.Wrap(db.Limit(int(showNumber)).Offset(int(showNumber*pageNumber)).Find(&users).Error, "")
	return
}

func (u *UserGorm) Get(ctx context.Context, showNumber, pageNumber int32) (users []*table.UserModel, count int64, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "showNumber", showNumber, "pageNumber", pageNumber, "users", users, "count", count)
	}()
	err = u.DB.Model(&table.UserModel{}).Model(u).Count(&count).Error
	if err != nil {
		return nil, 0, utils.Wrap(err, "")
	}
	err = utils.Wrap(u.DB.Model(&table.UserModel{}).Limit(int(showNumber)).Offset(int(pageNumber*showNumber)).Find(&users).Error, "")
	return
}
