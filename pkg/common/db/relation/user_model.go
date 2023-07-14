// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package relation

import (
	"context"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"

	"gorm.io/gorm"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
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
	return utils.Wrap(u.db(ctx).Model(&relation.UserModel{}).Where("user_id = ?", userID).Updates(args).Error, "")
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
func (u *UserGorm) Page(
	ctx context.Context,
	pageNumber, showNumber int32,
) (users []*relation.UserModel, count int64, err error) {
	err = utils.Wrap(u.db(ctx).Count(&count).Error, "")
	if err != nil {
		return
	}
	err = utils.Wrap(
		u.db(ctx).
			Limit(int(showNumber)).
			Offset(int((pageNumber-1)*showNumber)).
			Find(&users).
			Order("create_time DESC").
			Error,
		"",
	)
	return
}

// 获取所有用户ID
func (u *UserGorm) GetAllUserID(ctx context.Context, pageNumber, showNumber int32) (userIDs []string, err error) {
	return userIDs, errs.Wrap(u.db(ctx).Limit(int(showNumber)).Offset(int((pageNumber-1)*showNumber)).Pluck("user_id", &userIDs).Error)
}

func (u *UserGorm) GetUserGlobalRecvMsgOpt(ctx context.Context, userID string) (opt int, err error) {
	err = u.db(ctx).Model(&relation.UserModel{}).Where("user_id = ?", userID).Pluck("global_recv_msg_opt", &opt).Error
	return opt, err
}

func (u *UserGorm) CountTotal(ctx context.Context, before *time.Time) (count int64, err error) {
	db := u.db(ctx).Model(&relation.UserModel{})
	if before != nil {
		db = db.Where("create_time < ?", before)
	}
	if err := db.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (u *UserGorm) CountRangeEverydayTotal(
	ctx context.Context,
	start time.Time,
	end time.Time,
) (map[string]int64, error) {
	var res []struct {
		Date  time.Time `gorm:"column:date"`
		Count int64     `gorm:"column:count"`
	}
	err := u.db(ctx).
		Model(&relation.UserModel{}).
		Select("DATE(create_time) AS date, count(1) AS count").
		Where("create_time >= ? and create_time < ?", start, end).
		Group("date").
		Find(&res).
		Error
	if err != nil {
		return nil, errs.Wrap(err)
	}
	v := make(map[string]int64)
	for _, r := range res {
		v[r.Date.Format("2006-01-02")] = r.Count
	}
	return v, nil
}
