package relation

import (
	"OpenIM/pkg/common/db/table/relation"
	"OpenIM/pkg/common/log"
	"OpenIM/pkg/utils"
	"context"
	"gorm.io/gorm"
)

type UserGorm struct {
	DB *gorm.DB
}

func NewUserGorm(db *gorm.DB) relation.UserModelInterface {
	return &UserGorm{DB: db.Model(&relation.UserModel{})}
}

func (u *UserGorm) db() *gorm.DB {
	newDB := *u.DB
	return &newDB
}

// 插入多条
func (u *UserGorm) Create(ctx context.Context, users []*relation.UserModel) (err error) {
	return utils.Wrap(u.db().Create(&users).Error, "")
}

// 更新用户信息 零值
func (u *UserGorm) UpdateByMap(ctx context.Context, userID string, args map[string]interface{}) (err error) {
	return utils.Wrap(u.db().Where("user_id = ?", userID).Updates(args).Error, "")
}

// 更新多个用户信息 非零值
func (u *UserGorm) Update(ctx context.Context, users []*relation.UserModel) (err error) {
	return utils.Wrap(u.db().Updates(&users).Error, "")
}

// 获取指定用户信息  不存在，也不返回错误
func (u *UserGorm) Find(ctx context.Context, userIDs []string) (users []*relation.UserModel, err error) {
	log.ZDebug(ctx, "Find args", "userIDs", userIDs)
	err = utils.Wrap(u.db().Where("user_id in (?)", userIDs).Find(&users).Error, "")
	return users, err
}

// 获取某个用户信息  不存在，则返回错误
func (u *UserGorm) Take(ctx context.Context, userID string) (user *relation.UserModel, err error) {
	user = &relation.UserModel{}
	err = utils.Wrap(u.db().Where("user_id = ?", userID).Take(&user).Error, "")
	return user, err
}

// 获取用户信息 不存在，不返回错误
func (u *UserGorm) Page(ctx context.Context, pageNumber, showNumber int32) (users []*relation.UserModel, count int64, err error) {
	err = utils.Wrap(u.db().Count(&count).Error, "")
	if err != nil {
		return
	}
	err = utils.Wrap(u.db().Limit(int(showNumber)).Offset(int(pageNumber*showNumber)).Find(&users).Error, "")
	return
}

// 获取所有用户ID
func (u *UserGorm) GetAllUserID(ctx context.Context) (userIDs []string, err error) {
	err = u.db().Pluck("user_id", &userIDs).Error
	return userIDs, err
}
