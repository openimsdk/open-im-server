package mysql

import (
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"strconv"
	"time"
)

func ImportUserToSuperGroup() {
	for i := 18000000700; i <= 18000000800; i++ {
		user := db.User{
			UserID:           strconv.Itoa(i),
			Nickname:         strconv.Itoa(i),
			FaceURL:          "",
			Gender:           0,
			PhoneNumber:      strconv.Itoa(i),
			Birth:            time.Time{},
			Email:            "",
			Ex:               "",
			CreateTime:       time.Time{},
			AppMangerLevel:   0,
			GlobalRecvMsgOpt: 0,
		}
		err := im_mysql_model.UserRegister(user)
		if err != nil {
			log.NewError("", err.Error(), user)
			continue
		}

		groupMember := db.GroupMember{
			GroupID:        "3907826375",
			UserID:         strconv.Itoa(i),
			Nickname:       strconv.Itoa(i),
			FaceURL:        "",
			RoleLevel:      0,
			JoinTime:       time.Time{},
			JoinSource:     0,
			InviterUserID:  "openIMAdmin",
			OperatorUserID: "openIMAdmin",
			MuteEndTime:    time.Time{},
			Ex:             "",
		}

		err = im_mysql_model.InsertIntoGroupMember(groupMember)
		if err != nil {
			log.NewError("", err.Error(), user)
			continue
		}

		log.NewInfo("success", i)

	}

}
