package multi_terminal_login

import (
	"Open_IM/internal/push/content_struct"
	"Open_IM/internal/push/logic"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	pbChat "Open_IM/pkg/proto/chat"
	"Open_IM/pkg/utils"
)

func MultiTerminalLoginChecker(uid, token string, platformID int32) error {
	//	1.check userid and platform class   0 not exists and  1 exists
	existsInterface, err := db.DB.ExistsUserIDAndPlatform(uid, utils.PlatformNameToClass(utils.PlatformIDToName(platformID)))
	if err != nil {
		return err
	}
	exists := existsInterface.(int64)
	//get config multi login policy
	if config.Config.MultiLoginPolicy.OnlyOneTerminalAccess {
		//OnlyOneTerminalAccess policy need to check all terminal
		if utils.PlatformNameToClass(utils.PlatformIDToName(platformID)) == "PC" {
			existsInterface, err = db.DB.ExistsUserIDAndPlatform(uid, "Mobile")
			if err != nil {
				return err
			}
		} else {
			existsInterface, err = db.DB.ExistsUserIDAndPlatform(uid, "PC")
			if err != nil {
				return err
			}
		}
		exists = existsInterface.(int64)
		if exists == 1 {
			err := db.DB.SetUserIDAndPlatform(uid, utils.PlatformNameToClass(utils.PlatformIDToName(platformID)), token, config.Config.TokenPolicy.AccessExpire)
			if err != nil {
				return err
			}
			PushMessageToTheTerminal(uid, platformID)
			return nil
		}
	} else if config.Config.MultiLoginPolicy.MobileAndPCTerminalAccessButOtherTerminalKickEachOther {
		//	common terminal need to kick eich other
		if exists == 1 {
			err := db.DB.SetUserIDAndPlatform(uid, utils.PlatformNameToClass(utils.PlatformIDToName(platformID)), token, config.Config.TokenPolicy.AccessExpire)
			if err != nil {
				return err
			}
			PushMessageToTheTerminal(uid, platformID)
			return nil
		}
	}
	err = db.DB.SetUserIDAndPlatform(uid, utils.PlatformNameToClass(utils.PlatformIDToName(platformID)), token, config.Config.TokenPolicy.AccessExpire)
	if err != nil {
		return err
	}
	PushMessageToTheTerminal(uid, platformID)
	return nil
}

func PushMessageToTheTerminal(uid string, platform int32) {

	logic.SendMsgByWS(&pbChat.WSToMsgSvrChatMsg{
		SendID:      uid,
		RecvID:      uid,
		Content:     content_struct.NewContentStructString(1, "", "Your account is already logged on other terminal,please confirm"),
		SendTime:    utils.GetCurrentTimestampBySecond(),
		MsgFrom:     constant.SysMsgType,
		ContentType: constant.KickOnlineTip,
		PlatformID:  platform,
	})
}
