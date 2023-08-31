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

package mysql

import (
	"context"
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	config "github.com/OpenIMSDK/Open-IM-Server/tools/conversion/common"
	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/tools/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	MysqldbV2 *gorm.DB
	MysqldbV3 *gorm.DB
)

func init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.UsernameV2,
		config.PasswordV2,
		config.IpV2,
		config.DatabaseV2,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	MysqldbV2 = db
	if err != nil {
		log.ZDebug(context.Background(), "err", err)
	}

	dsnV3 := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.UsernameV3,
		config.PasswordV3,
		config.IpV3,
		config.DatabaseV3,
	)
	dbV3, err := gorm.Open(mysql.Open(dsnV3), &gorm.Config{})
	MysqldbV3 = dbV3
	if err != nil {
		log.ZDebug(context.Background(), "err", err)
	}
}

func UserConversion() {
	var count int64
	var user relation.UserModel
	MysqldbV2.Model(&user).Count(&count)
	batchSize := 100
	offset := 0

	for int64(offset) < count {
		var results []relation.UserModel
		MysqldbV2.Limit(batchSize).Offset(offset).Find(&results)
		// Process query results
		fmt.Println("============================batch data===================", offset, batchSize)
		fmt.Println(results)
		MysqldbV3.Create(results)
		fmt.Println("======================================================")
		offset += batchSize
	}
}

func FriendConversion() {
	var count int64
	var friend relation.FriendModel
	MysqldbV2.Model(&friend).Count(&count)
	batchSize := 100
	offset := 0

	for int64(offset) < count {
		var results []relation.FriendModel
		MysqldbV2.Limit(batchSize).Offset(offset).Find(&results)
		// Process query results
		fmt.Println("============================batch data===================", offset, batchSize)
		fmt.Println(results)
		MysqldbV3.Create(results)
		fmt.Println("======================================================")
		offset += batchSize
	}
}

func RequestConversion() {
	var count int64
	var friendRequest relation.FriendRequestModel
	MysqldbV2.Model(&friendRequest).Count(&count)
	batchSize := 100
	offset := 0

	for int64(offset) < count {
		var results []relation.FriendRequestModel
		MysqldbV2.Limit(batchSize).Offset(offset).Find(&results)
		// Process query results
		fmt.Println("============================batch data===================", offset, batchSize)
		fmt.Println(results)
		MysqldbV3.Create(results)
		fmt.Println("======================================================")
		offset += batchSize
	}

	var groupRequest relation.GroupRequestModel
	MysqldbV2.Model(&groupRequest).Count(&count)
	batchSize = 100
	offset = 0

	for int64(offset) < count {
		var results []relation.GroupRequestModel
		MysqldbV2.Limit(batchSize).Offset(offset).Find(&results)
		// Process query results
		fmt.Println("============================batch data===================", offset, batchSize)
		fmt.Println(results)
		MysqldbV3.Create(results)
		fmt.Println("======================================================")
		offset += batchSize
	}
}

func GroupConversion() {
	var count int64
	var group relation.GroupModel
	MysqldbV2.Model(&group).Count(&count)
	batchSize := 100
	offset := 0

	for int64(offset) < count {
		var results []relation.GroupModel
		MysqldbV2.Limit(batchSize).Offset(offset).Find(&results)
		for i, val := range results {
			results[i].GroupType = constant.WorkingGroup // After version 3.0, there is only one group type, which is the work group
			temp := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
			if val.NotificationUpdateTime.Equal(temp) {
				results[i].NotificationUpdateTime = time.Now()
				//fmt.Println(val.NotificationUpdateTime)
			}
		}
		// Process query results
		fmt.Println("============================batch data===================", offset, batchSize)
		fmt.Println(results)
		MysqldbV3.Create(results)
		fmt.Println("======================================================")
		offset += batchSize
	}
}

func GroupMemberConversion() {
	var count int64
	var groupMember relation.GroupMemberModel
	MysqldbV2.Model(&groupMember).Count(&count)
	batchSize := 100
	offset := 0

	for int64(offset) < count {
		var results []relation.GroupMemberModel
		MysqldbV2.Limit(batchSize).Offset(offset).Find(&results)
		// Process query results
		fmt.Println("============================batch data===================", offset, batchSize)
		fmt.Println(results)
		MysqldbV3.Create(results)
		fmt.Println("======================================================")
		offset += batchSize
	}
}

func BlacksConversion() {
	var count int64
	var black relation.BlackModel
	MysqldbV2.Model(&black).Count(&count)
	batchSize := 100
	offset := 0

	for int64(offset) < count {
		var results []relation.BlackModel
		MysqldbV2.Limit(batchSize).Offset(offset).Find(&results)
		// Process query results
		fmt.Println("============================batch data===================", offset, batchSize)
		fmt.Println(results)
		MysqldbV3.Create(results)
		fmt.Println("======================================================")
		offset += batchSize
	}
}

func ChatLogsConversion() {
	var count int64
	var chat relation.ChatLogModel
	MysqldbV2.Model(&chat).Count(&count)
	batchSize := 100
	offset := 0

	for int64(offset) < count {
		var results []relation.ChatLogModel
		MysqldbV2.Limit(batchSize).Offset(offset).Find(&results)
		// Process query results
		fmt.Println("============================batch data===================", offset, batchSize)
		//fmt.Println(results)
		MysqldbV3.Create(results)
		fmt.Println("======================================================")
		offset += batchSize
	}
}
