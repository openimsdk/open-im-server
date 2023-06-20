package api

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/a2r"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/apiresp"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/apistruct"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func NewMsg(discov discoveryregistry.SvcDiscoveryRegistry) *Message {
	conn, err := discov.GetConn(context.Background(), config.Config.RpcRegisterName.OpenImMsgName)
	if err != nil {
		panic(err)
	}
	client := msg.NewMsgClient(conn)
	return &Message{validate: validator.New(), discov: discov, conn: conn, client: client}
}

type Message struct {
	conn     *grpc.ClientConn
	client   msg.MsgClient
	validate *validator.Validate
	discov   discoveryregistry.SvcDiscoveryRegistry
}

func (Message) SetOptions(options map[string]bool, value bool) {
	utils.SetSwitchFromOptions(options, constant.IsHistory, value)
	utils.SetSwitchFromOptions(options, constant.IsPersistent, value)
	utils.SetSwitchFromOptions(options, constant.IsSenderSync, value)
	utils.SetSwitchFromOptions(options, constant.IsConversationUpdate, value)
}

func (m Message) newUserSendMsgReq(c *gin.Context, params *apistruct.ManagementSendMsgReq) *msg.SendMsgReq {
	var newContent string
	var err error
	switch params.ContentType {
	case constant.Text:
		newContent = params.Content["text"].(string)
	case constant.Picture:
		fallthrough
	case constant.Custom:
		fallthrough
	case constant.Voice:
		fallthrough
	case constant.Video:
		fallthrough
	case constant.File:
		fallthrough
	case constant.CustomNotTriggerConversation:
		fallthrough
	case constant.CustomOnlineOnly:
		fallthrough
	case constant.Revoke:
		newContent = params.Content["revokeMsgClientID"].(string)
	default:
	}
	options := make(map[string]bool, 5)
	if params.IsOnlineOnly {
		m.SetOptions(options, false)
	}
	if params.NotOfflinePush {
		utils.SetSwitchFromOptions(options, constant.IsOfflinePush, false)
	}
	if params.ContentType == constant.CustomOnlineOnly {
		m.SetOptions(options, false)
	} else if params.ContentType == constant.CustomNotTriggerConversation {
		utils.SetSwitchFromOptions(options, constant.IsConversationUpdate, false)
	}
	pbData := msg.SendMsgReq{
		MsgData: &sdkws.MsgData{
			SendID:           params.SendID,
			GroupID:          params.GroupID,
			ClientMsgID:      utils.GetMsgID(params.SendID),
			SenderPlatformID: params.SenderPlatformID,
			SenderNickname:   params.SenderNickname,
			SenderFaceURL:    params.SenderFaceURL,
			SessionType:      params.SessionType,
			MsgFrom:          constant.SysMsgType,
			ContentType:      params.ContentType,
			Content:          []byte(newContent),
			RecvID:           params.RecvID,
			CreateTime:       utils.GetCurrentTimestampByMill(),
			Options:          options,
			OfflinePushInfo:  params.OfflinePushInfo,
		},
	}
	if params.ContentType == constant.OANotification {
		var tips sdkws.TipsComm
		tips.JsonDetail = utils.StructToJsonString(params.Content)
		pbData.MsgData.Content, err = proto.Marshal(&tips)
		if err != nil {
			log.ZError(c, "Marshal failed ", err, "tips", tips.String())
		}
	}
	return &pbData
}

func (m *Message) Client() msg.MsgClient {
	return m.client
}

func (m *Message) GetSeq(c *gin.Context) {
	a2r.Call(msg.MsgClient.GetMaxSeq, m.Client, c)
}

func (m *Message) PullMsgBySeqs(c *gin.Context) {
	a2r.Call(msg.MsgClient.PullMessageBySeqs, m.Client, c)
}

func (m *Message) RevokeMsg(c *gin.Context) {
	a2r.Call(msg.MsgClient.RevokeMsg, m.Client, c)
}

func (m *Message) MarkMsgsAsRead(c *gin.Context) {
	a2r.Call(msg.MsgClient.MarkMsgsAsRead, m.Client, c)
}

func (m *Message) MarkConversationAsRead(c *gin.Context) {
	a2r.Call(msg.MsgClient.MarkConversationAsRead, m.Client, c)
}

func (m *Message) GetConversationsHasReadAndMaxSeq(c *gin.Context) {
	a2r.Call(msg.MsgClient.GetConversationsHasReadAndMaxSeq, m.Client, c)
}

func (m *Message) SetConversationHasReadSeq(c *gin.Context) {
	a2r.Call(msg.MsgClient.SetConversationHasReadSeq, m.Client, c)
}

func (m *Message) ClearConversationsMsg(c *gin.Context) {
	a2r.Call(msg.MsgClient.ClearConversationsMsg, m.Client, c)
}

func (m *Message) UserClearAllMsg(c *gin.Context) {
	a2r.Call(msg.MsgClient.UserClearAllMsg, m.Client, c)
}

func (m *Message) DeleteMsgs(c *gin.Context) {
	a2r.Call(msg.MsgClient.DeleteMsgs, m.Client, c)
}

func (m *Message) DeleteMsgPhysicalBySeq(c *gin.Context) {
	a2r.Call(msg.MsgClient.DeleteMsgPhysicalBySeq, m.Client, c)
}

func (m *Message) DeleteMsgPhysical(c *gin.Context) {
	a2r.Call(msg.MsgClient.DeleteMsgPhysical, m.Client, c)
}

func (m *Message) SetMessageReactionExtensions(c *gin.Context) {
	a2r.Call(msg.MsgClient.SetMessageReactionExtensions, m.Client, c)
}

func (m *Message) GetMessageListReactionExtensions(c *gin.Context) {
	a2r.Call(msg.MsgClient.GetMessagesReactionExtensions, m.Client, c)
}

func (m *Message) AddMessageReactionExtensions(c *gin.Context) {
	a2r.Call(msg.MsgClient.AddMessageReactionExtensions, m.Client, c)
}

func (m *Message) DeleteMessageReactionExtensions(c *gin.Context) {
	a2r.Call(msg.MsgClient.DeleteMessageReactionExtensions, m.Client, c)
}

func (m *Message) SendMessage(c *gin.Context) {
	params := apistruct.ManagementSendMsgReq{}
	if err := c.BindJSON(&params); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WithDetail(err.Error()).Wrap())
		return
	}
	// todo
	//if !tokenverify.IsAppManagerUid(c) {
	//	apiresp.GinError(c, errs.ErrNoPermission.Wrap("only app manager can send message"))
	//	return
	//}

	var data interface{}
	switch params.ContentType {
	case constant.Text:
		data = apistruct.TextElem{}
	case constant.Picture:
		data = apistruct.PictureElem{}
	case constant.Voice:
		data = apistruct.SoundElem{}
	case constant.Video:
		data = apistruct.VideoElem{}
	case constant.File:
		data = apistruct.FileElem{}
	case constant.Custom:
		data = apistruct.CustomElem{}
	case constant.Revoke:
		data = apistruct.RevokeElem{}
	case constant.OANotification:
		data = apistruct.OANotificationElem{}
		params.SessionType = constant.NotificationChatType
	case constant.CustomNotTriggerConversation:
		data = apistruct.CustomElem{}
	case constant.CustomOnlineOnly:
		data = apistruct.CustomElem{}
	default:
		apiresp.GinError(c, errs.ErrArgs.WithDetail("not support err contentType").Wrap())
		return
	}
	if err := mapstructure.WeakDecode(params.Content, &data); err != nil {
		apiresp.GinError(c, errs.ErrArgs.Wrap(err.Error()))
		return
	} else if err := m.validate.Struct(data); err != nil {
		apiresp.GinError(c, errs.ErrArgs.Wrap(err.Error()))
		return
	}
	pbReq := m.newUserSendMsgReq(c, &params)
	client := msg.NewMsgClient(m.conn)
	var status int
	respPb, err := client.SendMsg(c, pbReq)
	if err != nil {
		status = constant.MsgSendFailed
		apiresp.GinError(c, err)
		return
	}
	status = constant.MsgSendSuccessed
	_, err = client.SetSendMsgStatus(c, &msg.SetSendMsgStatusReq{
		Status: int32(status),
	})
	if err != nil {
		log.ZError(c, "SetSendMsgStatus failed", err)
	}
	apiresp.GinSuccess(c, respPb)
}

func (m *Message) ManagementBatchSendMsg(c *gin.Context) {
	a2r.Call(msg.MsgClient.SendMsg, m.Client, c)
}

func (m *Message) CheckMsgIsSendSuccess(c *gin.Context) {
	a2r.Call(msg.MsgClient.GetSendMsgStatus, m.Client, c)
}

func (m *Message) GetUsersOnlineStatus(c *gin.Context) {
	a2r.Call(msg.MsgClient.GetSendMsgStatus, m.Client, c)
}
