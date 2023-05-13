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

package multi_terminal_login

//
//import (
//	"Open_IM/internal/push/content_struct"
//	"Open_IM/internal/push/logic"
//	"Open_IM/pkg/common/config"
//	"Open_IM/pkg/common/constant"
//	"Open_IM/pkg/common/db"
//	pbChat "Open_IM/pkg/proto/msg"
//	"Open_IM/pkg/utils"
//)
//
//const  DayOfSecond = 24*60*60
//func MultiTerminalLoginChecker(uid, token string, platformID int32) error {
//	//	1.check userid and platform class   0 not exists and  1 exists
//	exists, err := db.DB.ExistsUserIDAndPlatform(uid, constant.PlatformNameToClass(constant.PlatformIDToName(platformID)))
//	if err != nil {
//		return err
//	}
//	//get config multi login policy
//	if config.Config.MultiLoginPolicy  {
//		//OnlyOneTerminalAccess policy need to check all terminal
//		if constant.PlatformNameToClass(constant.PlatformIDToName(platformID)) == "PC" {
//			exists, err = db.DB.ExistsUserIDAndPlatform(uid, "Mobile")
//			if err != nil {
//				return err
//			}
//		} else {
//			exists, err = db.DB.ExistsUserIDAndPlatform(uid, "PC")
//			if err != nil {
//				return err
//			}
//		}
//		if exists == 1 {
//			err := db.DB.SetUserIDAndPlatform(uid, constant.PlatformNameToClass(constant.PlatformIDToName(platformID)), token, config.Config.TokenPolicy.AccessExpire*DayOfSecond)
//			if err != nil {
//				return err
//			}
//			PushMessageToTheTerminal(uid, platformID)
//			return nil
//		}
//	} else if config.Config.MultiLoginPolicy.MobileAndPCTerminalAccessButOtherTerminalKickEachOther {
//		//	common terminal need to kick eich other
//		if exists == 1 {
//			err := db.DB.SetUserIDAndPlatform(uid, constant.PlatformNameToClass(constant.PlatformIDToName(platformID)), token, config.Config.TokenPolicy.AccessExpire*DayOfSecond)
//			if err != nil {
//				return err
//			}
//			PushMessageToTheTerminal(uid, platformID)
//			return nil
//		}
//	}
//	err = db.DB.SetUserIDAndPlatform(uid, constant.PlatformNameToClass(constant.PlatformIDToName(platformID)), token, config.Config.TokenPolicy.AccessExpire*DayOfSecond)
//	if err != nil {
//		return err
//	}
//	PushMessageToTheTerminal(uid, platformID)
//	return nil
//}
////
////func PushMessageToTheTerminal(uid string, platform int32) {
////
////	logic.SendMsgByWS(&pbChat.WSToMsgSvrChatMsg{
////		SendID:      uid,
////		RecvID:      uid,
////		Content:     content_struct.NewContentStructString(1, "", "Your account is already logged on other terminal,please confirm"),
////		SendTime:    utils.GetCurrentTimestampBySecond(),
////		MsgFrom:     constant.SysMsgType,
////		ContentType: constant.KickOnlineTip,
////		PlatformID:  platform,
////	})
////}
