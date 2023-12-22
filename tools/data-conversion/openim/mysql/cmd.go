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
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/openimsdk/open-im-server/v3/tools/data-conversion/openim/mysql/conversion"
	"github.com/openimsdk/open-im-server/v3/tools/data-conversion/utils"
)

func Cmd() {
	var (
		usernameV2 = "root"
		passwordV2 = "openIM"
		addrV2     = "121.5.182.23:13306"
		databaseV2 = "openIM_v2"
	)

	var (
		usernameV3 = "root"
		passwordV3 = "openIM123"
		addrV3     = "203.56.175.233:13306"
		databaseV3 = "openIM_v3"
	)
	log.SetFlags(log.LstdFlags | log.Llongfile)
	dsnV2 := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", usernameV2, passwordV2, addrV2, databaseV2)
	dsnV3 := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", usernameV3, passwordV3, addrV3, databaseV3)
	dbV2, err := gorm.Open(mysql.Open(dsnV2), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		log.Println("open v2 db failed", err)
		return
	}
	dbV3, err := gorm.Open(mysql.Open(dsnV3), &gorm.Config{Logger: logger.Discard})
	if err != nil {
		log.Println("open v3 db failed", err)
		return
	}

	var tasks utils.TakeList

	tasks.Append(func() (string, error) { return utils.FindAndInsert(dbV2, dbV3, conversion.Friend) })
	tasks.Append(func() (string, error) { return utils.FindAndInsert(dbV2, dbV3, conversion.FriendRequest) })
	tasks.Append(func() (string, error) { return utils.FindAndInsert(dbV2, dbV3, conversion.Group) })
	tasks.Append(func() (string, error) { return utils.FindAndInsert(dbV2, dbV3, conversion.GroupMember) })
	tasks.Append(func() (string, error) { return utils.FindAndInsert(dbV2, dbV3, conversion.GroupRequest) })
	tasks.Append(func() (string, error) { return utils.FindAndInsert(dbV2, dbV3, conversion.User) })

	utils.RunTask(4, tasks)

}
