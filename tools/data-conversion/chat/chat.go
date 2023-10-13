package main

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/openimsdk/open-im-server/v3/tools/data-conversion/chat/conversion"
	"github.com/openimsdk/open-im-server/v3/tools/data-conversion/utils"
)

func main() {
	var (
		usernameV2 = "root"            // v2版本mysql用户名
		passwordV2 = "openIM"          // v2版本mysql密码
		addrV2     = "127.0.0.1:13306" // v2版本mysql地址
		databaseV2 = "admin_chat"      // v2版本mysql数据库名字
	)

	var (
		usernameV3 = "root"              // v3版本mysql用户名
		passwordV3 = "openIM123"         // v3版本mysql密码
		addrV3     = "127.0.0.1:13306"   // v3版本mysql地址
		databaseV3 = "openim_enterprise" // v3版本mysql数据库名字
	)

	var concurrency = 1 // 并发数量

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

	tasks.Append(func() (string, error) { return utils.FindAndInsert(dbV2, dbV3, conversion.Account) })
	tasks.Append(func() (string, error) { return utils.FindAndInsert(dbV2, dbV3, conversion.Attribute) })
	tasks.Append(func() (string, error) { return utils.FindAndInsert(dbV2, dbV3, conversion.Register) })
	tasks.Append(func() (string, error) { return utils.FindAndInsert(dbV2, dbV3, conversion.UserLoginRecord) })
	tasks.Append(func() (string, error) { return utils.FindAndInsert(dbV2, dbV3, conversion.Admin) })
	tasks.Append(func() (string, error) { return utils.FindAndInsert(dbV2, dbV3, conversion.Applet) })
	tasks.Append(func() (string, error) { return utils.FindAndInsert(dbV2, dbV3, conversion.ForbiddenAccount) })
	tasks.Append(func() (string, error) { return utils.FindAndInsert(dbV2, dbV3, conversion.InvitationRegister) })
	tasks.Append(func() (string, error) { return utils.FindAndInsert(dbV2, dbV3, conversion.IPForbidden) })
	tasks.Append(func() (string, error) { return utils.FindAndInsert(dbV2, dbV3, conversion.LimitUserLoginIP) })
	tasks.Append(func() (string, error) { return utils.FindAndInsert(dbV2, dbV3, conversion.RegisterAddFriend) })
	tasks.Append(func() (string, error) { return utils.FindAndInsert(dbV2, dbV3, conversion.RegisterAddGroup) })

	utils.RunTask(concurrency, tasks)

}
