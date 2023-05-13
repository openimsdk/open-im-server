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

package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"errors"

	_ "gorm.io/gorm"
)

func GetRegister(account, areaCode, userID string) (*db.Register, error) {
	var r db.Register
	return &r, db.DB.MysqlDB.DefaultGormDB().Table("registers").Where("user_id = ? and user_id != ? or account = ? or account =? and area_code=?",
		userID, "", account, account, areaCode).Take(&r).Error
}

func GetRegisterInfo(userID string) (*db.Register, error) {
	var r db.Register
	return &r, db.DB.MysqlDB.DefaultGormDB().Table("registers").Where("user_id = ?", userID).Take(&r).Error
}

func SetPassword(account, password, ex, userID, areaCode, ip string) error {
	r := db.Register{
		Account:    account,
		Password:   password,
		Ex:         ex,
		UserID:     userID,
		RegisterIP: ip,
		AreaCode:   areaCode,
	}
	return db.DB.MysqlDB.DefaultGormDB().Table("registers").Create(&r).Error
}

func ResetPassword(account, password string) error {
	r := db.Register{
		Password: password,
	}
	return db.DB.MysqlDB.DefaultGormDB().Table("registers").Where("account = ?", account).Updates(&r).Error
}

func GetRegisterAddFriendList(showNumber, pageNumber int32) ([]string, error) {
	var IDList []string
	var err error
	model := db.DB.MysqlDB.DefaultGormDB().Model(&db.RegisterAddFriend{})
	if showNumber == 0 {
		err = model.Pluck("user_id", &IDList).Error
	} else {
		err = model.Limit(int(showNumber)).Offset(int(showNumber*(pageNumber-1))).Pluck("user_id", &IDList).Error
	}
	return IDList, err
}

func AddUserRegisterAddFriendIDList(userIDList ...string) error {
	var list []db.RegisterAddFriend
	for _, v := range userIDList {
		list = append(list, db.RegisterAddFriend{UserID: v})
	}
	result := db.DB.MysqlDB.DefaultGormDB().Create(list)
	if int(result.RowsAffected) < len(userIDList) {
		return errors.New("some line insert failed")
	}
	err := result.Error
	return err
}

func ReduceUserRegisterAddFriendIDList(userIDList ...string) error {
	var list []db.RegisterAddFriend
	for _, v := range userIDList {
		list = append(list, db.RegisterAddFriend{UserID: v})
	}
	err := db.DB.MysqlDB.DefaultGormDB().Delete(list).Error
	return err
}

func DeleteAllRegisterAddFriendIDList() error {
	err := db.DB.MysqlDB.DefaultGormDB().Where("1 = 1").Delete(&db.RegisterAddFriend{}).Error
	return err
}

func GetUserIPLimit(userID string) (db.UserIpLimit, error) {
	var limit db.UserIpLimit
	limit.UserID = userID
	err := db.DB.MysqlDB.DefaultGormDB().Model(&db.UserIpLimit{}).Take(&limit).Error
	return limit, err
}
