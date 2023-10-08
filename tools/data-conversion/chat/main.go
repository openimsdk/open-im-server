package main

import (
	"fmt"
	"github.com/openimsdk/open-im-server/v3/tools/data-conversion/chat/conversion"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func main() {
	var (
		usernameV2 = "root"
		passwordV2 = "openIM123"
		addrV2     = "127.0.0.1:13306"
		databaseV2 = "admin_chat"
	)

	var (
		usernameV3 = "root"
		passwordV3 = "openIM123"
		addrV3     = "127.0.0.1:13306"
		databaseV3 = "openim_enterprise"
	)

	dsnV2 := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", usernameV2, passwordV2, addrV2, databaseV2)
	dsnV3 := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", usernameV3, passwordV3, addrV3, databaseV3)
	dbV2, err := gorm.Open(mysql.Open(dsnV2), &gorm.Config{})
	if err != nil {
		log.Println("open v2 db failed", err)
		return
	}
	dbV3, err := gorm.Open(mysql.Open(dsnV3), &gorm.Config{})
	if err != nil {
		log.Println("open v3 db failed", err)
		return
	}

	var fns []func() (string, error)

	Append := func(fn func() (string, error)) {
		fns = append(fns, fn)
	}

	Append(func() (string, error) { return conversion.FindAndInsert(dbV2, dbV3, conversion.Account) })
	Append(func() (string, error) { return conversion.FindAndInsert(dbV2, dbV3, conversion.Attribute) })
	Append(func() (string, error) { return conversion.FindAndInsert(dbV2, dbV3, conversion.Register) })
	Append(func() (string, error) { return conversion.FindAndInsert(dbV2, dbV3, conversion.UserLoginRecord) })
	Append(func() (string, error) { return conversion.FindAndInsert(dbV2, dbV3, conversion.Admin) })
	Append(func() (string, error) { return conversion.FindAndInsert(dbV2, dbV3, conversion.Applet) })
	Append(func() (string, error) { return conversion.FindAndInsert(dbV2, dbV3, conversion.ForbiddenAccount) })
	Append(func() (string, error) { return conversion.FindAndInsert(dbV2, dbV3, conversion.InvitationRegister) })
	Append(func() (string, error) { return conversion.FindAndInsert(dbV2, dbV3, conversion.IPForbidden) })
	Append(func() (string, error) { return conversion.FindAndInsert(dbV2, dbV3, conversion.LimitUserLoginIP) })
	Append(func() (string, error) { return conversion.FindAndInsert(dbV2, dbV3, conversion.RegisterAddFriend) })
	Append(func() (string, error) { return conversion.FindAndInsert(dbV2, dbV3, conversion.RegisterAddGroup) })

	for i := range fns {
		name, err := fns[i]()
		if err == nil {
			log.Printf("[%d/%d] %s success\n", i+1, len(fns), name)
		} else {
			log.Printf("[%d/%d] %s failed %s\n", i+1, len(fns), name, err)
			return
		}
	}
}
