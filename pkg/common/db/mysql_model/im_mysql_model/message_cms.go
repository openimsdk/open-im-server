package im_mysql_model

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"fmt"
)

func GetChatLog(chatLog db.ChatLog, pageNumber, showNumber int32) ([]db.ChatLog, error) {
	var chatLogs []db.ChatLog
	db := db.DB.MysqlDB.DefaultGormDB().Table("chat_logs").
		Limit(int(showNumber)).Offset(int(showNumber * (pageNumber - 1)))
	if chatLog.SendTime.Unix() > 0 {
		db = db.Where("send_time > ? and send_time < ?", chatLog.SendTime, chatLog.SendTime.AddDate(0, 0, 1))
	}
	if chatLog.Content != "" {
		db = db.Where(" content like ? ", fmt.Sprintf("%%%s%%", chatLog.Content))
	}
	if chatLog.SessionType == 1 {
		db = db.Where("session_type = ?", chatLog.SessionType)
	} else if chatLog.SessionType == 2 {
		db = db.Where("content_type in (?)", []int{constant.GroupChatType, constant.SuperGroupChatType})
	}
	if chatLog.ContentType != 0 {
		db = db.Where("content_type = ?", chatLog.ContentType)
	}
	if chatLog.SendID != "" {
		db = db.Where("send_id = ?", chatLog.SendID)
	}
	if chatLog.RecvID != "" {
		db = db.Where("recv_id = ?", chatLog.RecvID)
	}

	err := db.Find(&chatLogs).Error
	return chatLogs, err
}

func GetChatLogCount(chatLog db.ChatLog) (int64, error) {
	var count int64
	db := db.DB.MysqlDB.DefaultGormDB().Table("chat_logs")
	if chatLog.SendTime.Unix() > 0 {
		log.NewDebug("", utils.GetSelfFuncName(), chatLog.SendTime, chatLog.SendTime.AddDate(0, 0, 1))
		db = db.Where("send_time > ? and send_time < ?", chatLog.SendTime, chatLog.SendTime.AddDate(0, 0, 1))
	}
	if chatLog.Content != "" {
		db = db.Where(" content like ? ", fmt.Sprintf("%%%s%%", chatLog.Content))
	}
	if chatLog.SessionType != 0 {
		db = db.Where("session_type = ?", chatLog.SessionType)
	}
	if chatLog.ContentType == 1 {
		db = db.Where("content_type = ?", chatLog.ContentType)
	} else if chatLog.ContentType == 2 {
		db = db.Where("content_type in (?)", []int{constant.GroupChatType, constant.SuperGroupChatType})
	}
	if chatLog.SendID != "" {
		db = db.Where("send_id = ?", chatLog.SendID)
	}
	if chatLog.RecvID != "" {
		db = db.Where("recv_id = ?", chatLog.RecvID)
	}

	err := db.Count(&count).Error
	return count, err
}
