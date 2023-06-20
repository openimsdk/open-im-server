package relation

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"gorm.io/gorm"
)

type UserGorm struct {
	*MetaDB
}

func NewUserGorm(db *gorm.DB) relation.UserModelInterface {
	return &UserGorm{NewMetaDB(db, &relation.UserModel{})}
}

// 插入多条
func (u *UserGorm) Create(ctx context.Context, users []*relation.UserModel) (err error) {
	return utils.Wrap(u.db(ctx).Create(&users).Error, "")
}

// 更新用户信息 零值
func (u *UserGorm) UpdateByMap(ctx context.Context, userID string, args map[string]interface{}) (err error) {
	return utils.Wrap(u.db(ctx).Where("user_id = ?", userID).Updates(args).Error, "")
}

// 更新多个用户信息 非零值
func (u *UserGorm) Update(ctx context.Context, user *relation.UserModel) (err error) {
	return utils.Wrap(u.db(ctx).Model(user).Updates(user).Error, "")
}

// 获取指定用户信息  不存在，也不返回错误
func (u *UserGorm) Find(ctx context.Context, userIDs []string) (users []*relation.UserModel, err error) {
	err = utils.Wrap(u.db(ctx).Where("user_id in (?)", userIDs).Find(&users).Error, "")
	return users, err
}

// 获取某个用户信息  不存在，则返回错误
func (u *UserGorm) Take(ctx context.Context, userID string) (user *relation.UserModel, err error) {
	user = &relation.UserModel{}
	err = utils.Wrap(u.db(ctx).Where("user_id = ?", userID).Take(&user).Error, "")
	return user, err
}

// 获取用户信息 不存在，不返回错误
func (u *UserGorm) Page(ctx context.Context, pageNumber, showNumber int32) (users []*relation.UserModel, count int64, err error) {
	err = utils.Wrap(u.db(ctx).Count(&count).Error, "")
	if err != nil {
		return
	}
	err = utils.Wrap(u.db(ctx).Limit(int(showNumber)).Offset(int((pageNumber-1)*showNumber)).Find(&users).Order("create_time DESC").Error, "")
	return
}

// 获取所有用户ID
func (u *UserGorm) GetAllUserID(ctx context.Context) (userIDs []string, err error) {
	err = u.db(ctx).Pluck("user_id", &userIDs).Error
	return userIDs, err
}

func (u *UserGorm) GetUserGlobalRecvMsgOpt(ctx context.Context, userID string) (opt int, err error) {
	err = u.db(ctx).Model(&relation.UserModel{}).Where("user_id = ?", userID).Pluck("global_recv_msg_opt", &opt).Error
	return opt, err
}

func (u *UserGorm) CountTotal(ctx context.Context) (count int64, err error) {
	err = u.db(ctx).Model(&relation.UserModel{}).Count(&count).Error
	return count, errs.Wrap(err)
}

func (u *UserGorm) CountRangeEverydayTotal(ctx context.Context, start time.Time, end time.Time) (map[string]int64, error) {
	var res []struct {
		Date  string `gorm:"column:date"`
		Count int64  `gorm:"column:count"`
	}
	err := u.db(ctx).Model(&relation.UserModel{}).Select("DATE(create_time) AS date, count(1) AS count").Where("create_time >= ? and create_time < ?", start, end).Group("date").Find(&res).Error
	if err != nil {
		return nil, errs.Wrap(err)
	}
	v := make(map[string]int64)
	for _, r := range res {
		v[r.Date] = r.Count
	}
	return v, nil
}
