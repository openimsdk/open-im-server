package relation

import (
	"fmt"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	pbMsg "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	sdkws "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/golang/protobuf/jsonpb"
	"github.com/jinzhu/copier"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
)

type ChatLogGorm struct {
	*MetaDB
}

func NewChatLogGorm(db *gorm.DB) relation.ChatLogModelInterface {
	return &ChatLogGorm{NewMetaDB(db, &relation.ChatLogModel{})}
}

func (c *ChatLogGorm) Create(msg *pbMsg.MsgDataToMQ) error {
	chatLog := new(relation.ChatLogModel)
	copier.Copy(chatLog, msg.MsgData)
	switch msg.MsgData.SessionType {
	case constant.GroupChatType, constant.SuperGroupChatType:
		chatLog.RecvID = msg.MsgData.GroupID
	case constant.SingleChatType:
		chatLog.RecvID = msg.MsgData.RecvID
	}
	if msg.MsgData.ContentType >= constant.NotificationBegin && msg.MsgData.ContentType <= constant.NotificationEnd {
		var tips sdkws.TipsComm
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
