package main

import (
	"fmt"
	"github.com/openimsdk/open-im-server/v3/tools/data-conversion/chat/conversion"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

func main() {
	const (
		usernameV2 = "root"
		passwordV2 = "openIM123"
		addrV2     = "127.0.0.1:3306"
		databaseV2 = "admin_chat"
	)

	const (
		usernameV3 = "root"
		passwordV3 = "openIM123"
		addrV3     = "127.0.0.1:3306"
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
	var fns []func() error

	Append := func(fn func() error) {
		fns = append(fns, fn)
	}

	Append(func() error { return conversion.FindAndInsert(dbV2, dbV3, conversion.Account) })
	Append(func() error { return conversion.FindAndInsert(dbV2, dbV3, conversion.Attribute) })
	Append(func() error { return conversion.FindAndInsert(dbV2, dbV3, conversion.Register) })
	Append(func() error { return conversion.FindAndInsert(dbV2, dbV3, conversion.UserLoginRecord) })
	Append(func() error { return conversion.FindAndInsert(dbV2, dbV3, conversion.Admin) })
	Append(func() error { return conversion.FindAndInsert(dbV2, dbV3, conversion.Applet) })
	Append(func() error { return conversion.FindAndInsert(dbV2, dbV3, conversion.ForbiddenAccount) })
	Append(func() error { return conversion.FindAndInsert(dbV2, dbV3, conversion.InvitationRegister) })
	Append(func() error { return conversion.FindAndInsert(dbV2, dbV3, conversion.IPForbidden) })
	Append(func() error { return conversion.FindAndInsert(dbV2, dbV3, conversion.LimitUserLoginIP) })
	Append(func() error { return conversion.FindAndInsert(dbV2, dbV3, conversion.RegisterAddFriend) })
	Append(func() error { return conversion.FindAndInsert(dbV2, dbV3, conversion.RegisterAddGroup) })

	for i := range fns {
		if err := fns[i](); err != nil {
			log.Printf("[%d] %s\n", i, err)
			return
		}
	}
}
