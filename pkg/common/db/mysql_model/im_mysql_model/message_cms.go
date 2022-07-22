package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"fmt"
)

func GetChatLog(chatLog db.ChatLog, pageNumber, showNumber int32) ([]db.ChatLog, error) {
	var chatLogs []db.ChatLog
	db := db.DB.MysqlDB.DefaultGormDB().Table("chat_logs").
		Where(fmt.Sprintf(" content like '%%%s%%'", chatLog.Content)).
		Limit(int(showNumber)).Offset(int(showNumber * (pageNumber - 1)))
	if chatLog.SessionType != 0 {
		db = db.Where("session_type = ?", chatLog.SessionType)
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
	if chatLog.SendTime.Unix() > 0 {
		db = db.Where("send_time > ? and send_time < ?", chatLog.SendTime, chatLog.SendTime.AddDate(0, 0, 1))
	}
	err := db.Find(&chatLogs).Error
	return chatLogs, err
}

func GetChatLogCount(chatLog db.ChatLog) (int64, error) {
	var chatLogs []db.ChatLog
	var count int64
	db := db.DB.MysqlDB.DefaultGormDB().Table("chat_logs").
		Where(fmt.Sprintf(" content like '%%%s%%'", chatLog.Content))
	if chatLog.SessionType != 0 {
		db = db.Where("session_type = ?", chatLog.SessionType)
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
	if chatLog.SendTime.Unix() > 0 {
		log.NewDebug("", utils.GetSelfFuncName(), chatLog.SendTime, chatLog.SendTime.AddDate(0, 0, 1))
		db = db.Where("send_time > ? and send_time < ?", chatLog.SendTime, chatLog.SendTime.AddDate(0, 0, 1))
	}

	err := db.Find(&chatLogs).Count(&count).Error
	return count, err
}
