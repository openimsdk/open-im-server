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

package data_conversion

import (
	"context"
	"fmt"
	"time"

	"github.com/OpenIMSDK/tools/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
)

var (
	MysqlDb_v2 *gorm.DB
	MysqlDb_v3 *gorm.DB
)

const (
	username_v2 = "root"
	password_v2 = "123456"
	ip_v2       = "127.0.0.1:3306"
	database_v2 = "openim_v2"
)

const (
	username_v3 = "root"
	password_v3 = "123456"
	ip_v3       = "127.0.0.1:3306"
	database_v3 = "openim_v3"
)

func init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username_v2,
		password_v2,
		ip_v2,
		database_v2,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	MysqlDb_v2 = db
	if err != nil {
		log.ZDebug(context.Background(), "err", err)
	}

	dsn_v3 := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username_v3,
		password_v3,
		ip_v3,
		database_v3,
	)
	db_v3, err := gorm.Open(mysql.Open(dsn_v3), &gorm.Config{})
	MysqlDb_v3 = db_v3
	if err != nil {
		log.ZDebug(context.Background(), "err", err)
	}
}

func UserConversion() {
	var count int64
	var user relation.UserModel
	MysqlDb_v2.Model(&user).Count(&count)
	batchSize := 100
	offset := 0

	for int64(offset) < count {
		var results []relation.UserModel
		MysqlDb_v2.Limit(batchSize).Offset(offset).Find(&results)
		// Process query results
		fmt.Println("============================batch data===================", offset, batchSize)
		fmt.Println(results)
		MysqlDb_v3.Create(results)
		fmt.Println("======================================================")
		offset += batchSize
	}
}

func FriendConversion() {
	var count int64
	var friend relation.FriendModel
	MysqlDb_v2.Model(&friend).Count(&count)
	batchSize := 100
	offset := 0

	for int64(offset) < count {
		var results []relation.FriendModel
		MysqlDb_v2.Limit(batchSize).Offset(offset).Find(&results)
		// Process query results
		fmt.Println("============================batch data===================", offset, batchSize)
		fmt.Println(results)
		MysqlDb_v3.Create(results)
		fmt.Println("======================================================")
		offset += batchSize
	}
}

func RequestConversion() {
	var count int64
	var friendRequest relation.FriendRequestModel
	MysqlDb_v2.Model(&friendRequest).Count(&count)
	batchSize := 100
	offset := 0

	for int64(offset) < count {
		var results []relation.FriendRequestModel
		MysqlDb_v2.Limit(batchSize).Offset(offset).Find(&results)
		// Process query results
		fmt.Println("============================batch data===================", offset, batchSize)
		fmt.Println(results)
		MysqlDb_v3.Create(results)
		fmt.Println("======================================================")
		offset += batchSize
	}

	var groupRequest relation.GroupRequestModel
	MysqlDb_v2.Model(&groupRequest).Count(&count)
	batchSize = 100
	offset = 0

	for int64(offset) < count {
		var results []relation.GroupRequestModel
		MysqlDb_v2.Limit(batchSize).Offset(offset).Find(&results)
		// Process query results
		fmt.Println("============================batch data===================", offset, batchSize)
		fmt.Println(results)
		MysqlDb_v3.Create(results)
		fmt.Println("======================================================")
		offset += batchSize
	}
}

func GroupConversion() {
	var count int64
	var group relation.GroupModel
	MysqlDb_v2.Model(&group).Count(&count)
	batchSize := 100
	offset := 0

	for int64(offset) < count {
		var results []relation.GroupModel
		MysqlDb_v2.Limit(batchSize).Offset(offset).Find(&results)
		for i, val := range results {
			temp := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
			if val.NotificationUpdateTime.Equal(temp) {
				results[i].NotificationUpdateTime = time.Now()
				// fmt.Println(val.NotificationUpdateTime)
			}
		}
		// Process query results
		fmt.Println("============================batch data===================", offset, batchSize)
		fmt.Println(results)
		MysqlDb_v3.Create(results)
		fmt.Println("======================================================")
		offset += batchSize
	}
}

func GroupMemberConversion() {
	var count int64
	var groupMember relation.GroupMemberModel
	MysqlDb_v2.Model(&groupMember).Count(&count)
	batchSize := 100
	offset := 0

	for int64(offset) < count {
		var results []relation.GroupMemberModel
		MysqlDb_v2.Limit(batchSize).Offset(offset).Find(&results)
		// Process query results
		fmt.Println("============================batch data===================", offset, batchSize)
		fmt.Println(results)
		MysqlDb_v3.Create(results)
		fmt.Println("======================================================")
		offset += batchSize
	}
}

func BlacksConversion() {
	var count int64
	var black relation.BlackModel
	MysqlDb_v2.Model(&black).Count(&count)
	batchSize := 100
	offset := 0

	for int64(offset) < count {
		var results []relation.BlackModel
		MysqlDb_v2.Limit(batchSize).Offset(offset).Find(&results)
		// Process query results
		fmt.Println("============================batch data===================", offset, batchSize)
		fmt.Println(results)
		MysqlDb_v3.Create(results)
		fmt.Println("======================================================")
		offset += batchSize
	}
}

func ChatLogsConversion() {
	var count int64
	var chat relation.ChatLogModel
	MysqlDb_v2.Model(&chat).Count(&count)
	batchSize := 100
	offset := 0

	for int64(offset) < count {
		var results []relation.ChatLogModel
		MysqlDb_v2.Limit(batchSize).Offset(offset).Find(&results)
		// Process query results
		fmt.Println("============================batch data===================", offset, batchSize)
		// fmt.Println(results)
		MysqlDb_v3.Create(results)
		fmt.Println("======================================================")
		offset += batchSize
	}
}
