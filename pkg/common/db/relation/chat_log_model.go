package relation

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/table/relation"
	pbMsg "Open_IM/pkg/proto/msg"
	server_api_params "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type ChatLogGorm struct {
	DB *gorm.DB
}

func NewChatLog(db *gorm.DB) *ChatLogGorm {
	return &ChatLogGorm{DB: db}
}

func (c *ChatLogGorm) Create(msg pbMsg.MsgDataToMQ) error {
	chatLog := new(relation.ChatLogModel)
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

func (c *ChatLogGorm) GetChatLog(chatLog *relation.ChatLogModel, pageNumber, showNumber int32, contentTypeList []int32) (int64, []relation.ChatLogModel, error) {
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
	var chatLogs []relation.ChatLogModel
	mdb = mdb.Limit(int(showNumber)).Offset(int(showNumber * (pageNumber - 1)))
	if err := mdb.Find(&chatLogs).Error; err != nil {
		return 0, nil, err
	}
	return count, chatLogs, nil
}
