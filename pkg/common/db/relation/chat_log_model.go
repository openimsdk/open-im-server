package relation

import (
	"Open_IM/pkg/common/constant"
	pbMsg "Open_IM/pkg/proto/msg"
	server_api_params "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"time"
)

type ChatLog struct {
	ServerMsgID      string    `gorm:"column:server_msg_id;primary_key;type:char(64)" json:"serverMsgID"`
	ClientMsgID      string    `gorm:"column:client_msg_id;type:char(64)" json:"clientMsgID"`
	SendID           string    `gorm:"column:send_id;type:char(64);index:send_id,priority:2" json:"sendID"`
	RecvID           string    `gorm:"column:recv_id;type:char(64);index:recv_id,priority:2" json:"recvID"`
	SenderPlatformID int32     `gorm:"column:sender_platform_id" json:"senderPlatformID"`
	SenderNickname   string    `gorm:"column:sender_nick_name;type:varchar(255)" json:"senderNickname"`
	SenderFaceURL    string    `gorm:"column:sender_face_url;type:varchar(255);" json:"senderFaceURL"`
	SessionType      int32     `gorm:"column:session_type;index:session_type,priority:2;index:session_type_alone" json:"sessionType"`
	MsgFrom          int32     `gorm:"column:msg_from" json:"msgFrom"`
	ContentType      int32     `gorm:"column:content_type;index:content_type,priority:2;index:content_type_alone" json:"contentType"`
	Content          string    `gorm:"column:content;type:varchar(3000)" json:"content"`
	Status           int32     `gorm:"column:status" json:"status"`
	SendTime         time.Time `gorm:"column:send_time;index:sendTime;index:content_type,priority:1;index:session_type,priority:1;index:recv_id,priority:1;index:send_id,priority:1" json:"sendTime"`
	CreateTime       time.Time `gorm:"column:create_time" json:"createTime"`
	Ex               string    `gorm:"column:ex;type:varchar(1024)" json:"ex"`
	DB               *gorm.DB
}

func (ChatLog) TableName() string {
	return "chat_logs"
}

func NewChatLog(db *gorm.DB) *ChatLog {
	return &ChatLog{DB: db}
}

func (c *ChatLog) Create(msg pbMsg.MsgDataToMQ) error {
	chatLog := new(ChatLog)
	copier.Copy(chatLog, msg.MsgData)
	switch msg.MsgData.SessionType {
	case constant.GroupChatType, constant.SuperGroupChatType:
		chatLog.RecvID = msg.MsgData.GroupID
	case constant.SingleChatType:
		chatLog.RecvID = msg.MsgData.RecvID
	}
	if msg.MsgData.ContentType >= constant.NotificationBegin && msg.MsgData.ContentType <= constant.NotificationEnd {
		var tips server_api_params.TipsComm
		_ = proto.Unmarshal(msg.MsgData.Content, &tips)
		marshaler := jsonpb.Marshaler{
			OrigName:     true,
			EnumsAsInts:  false,
			EmitDefaults: false,
		}
		chatLog.Content, _ = marshaler.MarshalToString(&tips)
	} else {
		chatLog.Content = string(msg.MsgData.Content)
	}
	chatLog.CreateTime = utils.UnixMillSecondToTime(msg.MsgData.CreateTime)
	chatLog.SendTime = utils.UnixMillSecondToTime(msg.MsgData.SendTime)
	return c.DB.Create(chatLog).Error
}

func (c *ChatLog) GetChatLog(chatLog *ChatLog, pageNumber, showNumber int32, contentTypeList []int32) (int64, []ChatLog, error) {
	mdb := c.DB.Model(chatLog)
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
