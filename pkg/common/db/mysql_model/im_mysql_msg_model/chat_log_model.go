/*
** description("").
** copyright('tuoyun,www.tuoyun.net').
** author("fg,Gordon@tuoyun.net").
** time(2021/3/4 11:18).
 */
package im_mysql_msg_model

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	pbMsg "Open_IM/pkg/proto/chat"
	"Open_IM/pkg/utils"
	"github.com/jinzhu/copier"
)

func InsertMessageToChatLog(msg pbMsg.MsgDataToMQ) error {
	dbConn, err := db.DB.MysqlDB.DefaultGormDB()
	if err != nil {
		return err
	}
	chatLog := new(db.ChatLog)
	copier.Copy(chatLog, msg.MsgData)
	switch msg.MsgData.SessionType {
	case constant.GroupChatType:
		chatLog.RecvID = msg.MsgData.GroupID
	case constant.SingleChatType:
		chatLog.RecvID = msg.MsgData.RecvID
	}
	chatLog.Content = string(msg.MsgData.Content)
	chatLog.CreateTime = utils.UnixNanoSecondToTime(msg.MsgData.CreateTime)
	chatLog.SendTime = utils.UnixNanoSecondToTime(msg.MsgData.SendTime)
	return dbConn.Table("chat_logs").Create(chatLog).Error
}
