package im_mysql_model

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"fmt"
)

func GetChatLog(chatLog *ChatLog, pageNumber, showNumber int32, contentTypeList []int32) (int64, []ChatLog, error) {
	mdb := db.DB.MysqlDB.DefaultGormDB().Table("chat_logs")
	if chatLog.SendTime.Unix() > 0 {
		mdb = mdb.Where("send_time > ? and send_time < ?", chatLog.SendTime, chatLog.SendTime.AddDate(0, 0, 1))
	}
	if chatLog.Content != "" {
		mdb = mdb.Where(" content like ? ", fmt.Sprintf("%%%s%%", chatLog.Content))
	}
	if chatLog.SessionType == 1 {
		mdb = mdb.Where("session_type = ?", chatLog.SessionType)
	} else if chatLog.SessionType == 2 {
		mdb = mdb.Where("session_type in (?)", []int{constant.GroupChatType, constant.SuperGroupChatType})
	}
	if chatLog.ContentType != 0 {
		mdb = mdb.Where("content_type = ?", chatLog.ContentType)
	}
	if chatLog.SendID != "" {
		mdb = mdb.Where("send_id = ?", chatLog.SendID)
	}
	if chatLog.RecvID != "" {
		mdb = mdb.Where("recv_id = ?", chatLog.RecvID)
	}
	if len(contentTypeList) > 0 {
		mdb = mdb.Where("content_type in (?)", contentTypeList)
	}
	var count int64
	if err := mdb.Count(&count).Error; err != nil {
		return 0, nil, err
	}
	var chatLogs []ChatLog
	mdb = mdb.Limit(int(showNumber)).Offset(int(showNumber * (pageNumber - 1)))
	if err := mdb.Find(&chatLogs).Error; err != nil {
		return 0, nil, err
	}
	return count, chatLogs, nil
}
