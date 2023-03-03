package api

import (
	"OpenIM/internal/api/a2r"
	"OpenIM/internal/apiresp"
	"OpenIM/pkg/apistruct"
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/log"
	"OpenIM/pkg/proto/msg"
	"OpenIM/pkg/proto/sdkws"
	"OpenIM/pkg/utils"
	"context"
	"errors"
	"github.com/OpenIMSDK/openKeeper"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang/protobuf/proto"
	"github.com/mitchellh/mapstructure"
	"net/http"
)

var _ context.Context // 解决goland编辑器bug

func NewMsg(zk *openKeeper.ZkClient) *Msg {
	return &Msg{zk: zk}
}

type Msg struct {
	zk *openKeeper.ZkClient
}

var validate *validator.Validate

func SetOptions(options map[string]bool, value bool) {
	utils.SetSwitchFromOptions(options, constant.IsHistory, value)
	utils.SetSwitchFromOptions(options, constant.IsPersistent, value)
	utils.SetSwitchFromOptions(options, constant.IsSenderSync, value)
	utils.SetSwitchFromOptions(options, constant.IsConversationUpdate, value)
}

func newUserSendMsgReq(params *apistruct.ManagementSendMsgReq) *msg.SendMsgReq {
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
	case constant.AdvancedRevoke:
		newContent = utils.StructToJsonString(params.Content)
	case constant.Revoke:
		newContent = params.Content["revokeMsgClientID"].(string)
	default:
	}
	options := make(map[string]bool, 5)
	if params.IsOnlineOnly {
		SetOptions(options, false)
	}
	if params.NotOfflinePush {
		utils.SetSwitchFromOptions(options, constant.IsOfflinePush, false)
	}
	if params.ContentType == constant.CustomOnlineOnly {
		SetOptions(options, false)
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
			log.Error(params.OperationID, "Marshal failed ", err.Error(), tips.String())
		}
	}
	return &pbData
}
func init() {
	validate = validator.New()
}

func (o *Msg) client() (msg.MsgClient, error) {
	conn, err := o.zk.GetConn(config.Config.RpcRegisterName.OpenImMsgName)
	if err != nil {
		return nil, err
	}
	return msg.NewMsgClient(conn), nil
}

func (o *Msg) GetSeq(c *gin.Context) {
	a2r.Call(msg.MsgClient.GetMaxAndMinSeq, o.client, c)
}

func (o *Msg) SendMsg(c *gin.Context) {
	a2r.Call(msg.MsgClient.SendMsg, o.client, c)
}

func (o *Msg) PullMsgBySeqs(c *gin.Context) {
	a2r.Call(msg.MsgClient.PullMessageBySeqs, o.client, c)
}

func (o *Msg) DelMsg(c *gin.Context) {
	a2r.Call(msg.MsgClient.DelMsgs, o.client, c)
}

func (o *Msg) DelSuperGroupMsg(c *gin.Context) {
	a2r.Call(msg.MsgClient.DelSuperGroupMsg, o.client, c)
}

func (o *Msg) ClearMsg(c *gin.Context) {
	a2r.Call(msg.MsgClient.ClearMsg, o.client, c)
}

func (o *Msg) SetMessageReactionExtensions(c *gin.Context) {
	a2r.Call(msg.MsgClient.SetMessageReactionExtensions, o.client, c)
}

func (o *Msg) GetMessageListReactionExtensions(c *gin.Context) {
	a2r.Call(msg.MsgClient.GetMessagesReactionExtensions, o.client, c)
}

func (o *Msg) AddMessageReactionExtensions(c *gin.Context) {
	a2r.Call(msg.MsgClient.AddMessageReactionExtensions, o.client, c)
}

func (o *Msg) DeleteMessageReactionExtensions(c *gin.Context) {
	a2r.Call(msg.MsgClient.DeleteMessageReactionExtensions, o.client, c)
}

func (o *Msg) ManagementSendMsg(c *gin.Context) {
	var data interface{}
	params := apistruct.ManagementSendMsgReq{}
	if err := c.BindJSON(&params); err != nil {
		apiresp.GinError(c, err)
		return
	}
	switch params.ContentType {
	case constant.Text:
		data = TextElem{}
	case constant.Picture:
		data = PictureElem{}
	case constant.Voice:
		data = SoundElem{}
	case constant.Video:
		data = VideoElem{}
	case constant.File:
		data = FileElem{}
	case constant.Custom:
		data = CustomElem{}
	case constant.Revoke:
		data = RevokeElem{}
	case constant.AdvancedRevoke:
		data = MessageRevoked{}
	case constant.OANotification:
		data = OANotificationElem{}
		params.SessionType = constant.NotificationChatType
	case constant.CustomNotTriggerConversation:
		data = CustomElem{}
	case constant.CustomOnlineOnly:
		data = CustomElem{}
	//case constant.HasReadReceipt:
	//case constant.Typing:
	//case constant.Quote:
	default:
		apiresp.GinError(c, errors.New("wrong contentType"))
		return
	}
	if err := mapstructure.WeakDecode(params.Content, &data); err != nil {
		apiresp.GinError(c, constant.ErrData)
		return
	} else if err := validate.Struct(data); err != nil {
		apiresp.GinError(c, constant.ErrData)
		return
	}
	log.NewInfo(params.OperationID, data, params)
	switch params.SessionType {
	case constant.SingleChatType:
		if len(params.RecvID) == 0 {
			apiresp.GinError(c, constant.ErrData)
			return
		}
	case constant.GroupChatType, constant.SuperGroupChatType:
		if len(params.GroupID) == 0 {
			apiresp.GinError(c, constant.ErrData)
			return
		}
	}
	log.NewInfo(params.OperationID, "Ws call success to ManagementSendMsgReq", params)
	pbData := newUserSendMsgReq(&params)
	conn, err := o.zk.GetConn(config.Config.RpcRegisterName.OpenImMsgName)
	if err != nil {
		apiresp.GinError(c, constant.ErrInternalServer)
		return
	}
	client := msg.NewMsgClient(conn)
	log.Info(params.OperationID, "", "api ManagementSendMsg call, api call rpc...")
	//var status int32
	RpcResp, err := client.SendMsg(context.Background(), pbData)
	if err != nil {
		//status = constant.MsgSendFailed
		apiresp.GinError(c, err)
		return
	}
	//status = constant.MsgSendSuccessed
	//_, err2 := client.SetSendMsgStatus(context.Background(), &msg.SetSendMsgStatusReq{OperationID: params.OperationID, Status: status})
	//if err2 != nil {
	//	log.NewError(params.OperationID, utils.GetSelfFuncName(), err2.Error())
	//}
	log.Info(params.OperationID, "", "api ManagementSendMsg call end..., [data: %s] [reply: %s]", pbData.String(), RpcResp.String())
	resp := apistruct.ManagementSendMsgResp{ResultList: sdkws.UserSendMsgResp{ServerMsgID: RpcResp.ServerMsgID, ClientMsgID: RpcResp.ClientMsgID, SendTime: RpcResp.SendTime}}
	log.Info(params.OperationID, "ManagementSendMsg return", resp)
	c.JSON(http.StatusOK, resp)
}

func (o *Msg) ManagementBatchSendMsg(c *gin.Context) {
	a2r.Call(msg.MsgClient.SendMsg, o.client, c)
}

func (o *Msg) CheckMsgIsSendSuccess(c *gin.Context) {
	a2r.Call(msg.MsgClient.GetSendMsgStatus, o.client, c)
}

func (o *Msg) GetUsersOnlineStatus(c *gin.Context) {
	a2r.Call(msg.MsgClient.GetSendMsgStatus, o.client, c)
}

func (o *Msg) AccountCheck(c *gin.Context) {
	a2r.Call(msg.MsgClient.GetSendMsgStatus, o.client, c)
}

type PictureBaseInfo struct {
	UUID   string `mapstructure:"uuid"`
	Type   string `mapstructure:"type" `
	Size   int64  `mapstructure:"size" `
	Width  int32  `mapstructure:"width" `
	Height int32  `mapstructure:"height"`
	Url    string `mapstructure:"url" `
}

type PictureElem struct {
	SourcePath      string          `mapstructure:"sourcePath"`
	SourcePicture   PictureBaseInfo `mapstructure:"sourcePicture"`
	BigPicture      PictureBaseInfo `mapstructure:"bigPicture" `
	SnapshotPicture PictureBaseInfo `mapstructure:"snapshotPicture"`
}
type SoundElem struct {
	UUID      string `mapstructure:"uuid"`
	SoundPath string `mapstructure:"soundPath"`
	SourceURL string `mapstructure:"sourceUrl"`
	DataSize  int64  `mapstructure:"dataSize"`
	Duration  int64  `mapstructure:"duration"`
}
type VideoElem struct {
	VideoPath      string `mapstructure:"videoPath"`
	VideoUUID      string `mapstructure:"videoUUID"`
	VideoURL       string `mapstructure:"videoUrl"`
	VideoType      string `mapstructure:"videoType"`
	VideoSize      int64  `mapstructure:"videoSize"`
	Duration       int64  `mapstructure:"duration"`
	SnapshotPath   string `mapstructure:"snapshotPath"`
	SnapshotUUID   string `mapstructure:"snapshotUUID"`
	SnapshotSize   int64  `mapstructure:"snapshotSize"`
	SnapshotURL    string `mapstructure:"snapshotUrl"`
	SnapshotWidth  int32  `mapstructure:"snapshotWidth"`
	SnapshotHeight int32  `mapstructure:"snapshotHeight"`
}
type FileElem struct {
	FilePath  string `mapstructure:"filePath"`
	UUID      string `mapstructure:"uuid"`
	SourceURL string `mapstructure:"sourceUrl"`
	FileName  string `mapstructure:"fileName"`
	FileSize  int64  `mapstructure:"fileSize"`
}
type AtElem struct {
	Text       string   `mapstructure:"text"`
	AtUserList []string `mapstructure:"atUserList"`
	IsAtSelf   bool     `mapstructure:"isAtSelf"`
}
type LocationElem struct {
	Description string  `mapstructure:"description"`
	Longitude   float64 `mapstructure:"longitude"`
	Latitude    float64 `mapstructure:"latitude"`
}
type CustomElem struct {
	Data        string `mapstructure:"data" validate:"required"`
	Description string `mapstructure:"description"`
	Extension   string `mapstructure:"extension"`
}
type TextElem struct {
	Text string `mapstructure:"text" validate:"required"`
}

type RevokeElem struct {
	RevokeMsgClientID string `mapstructure:"revokeMsgClientID" validate:"required"`
}
type OANotificationElem struct {
	NotificationName    string      `mapstructure:"notificationName" json:"notificationName" validate:"required"`
	NotificationFaceURL string      `mapstructure:"notificationFaceURL" json:"notificationFaceURL"`
	NotificationType    int32       `mapstructure:"notificationType" json:"notificationType" validate:"required"`
	Text                string      `mapstructure:"text" json:"text" validate:"required"`
	Url                 string      `mapstructure:"url" json:"url"`
	MixType             int32       `mapstructure:"mixType" json:"mixType"`
	PictureElem         PictureElem `mapstructure:"pictureElem" json:"pictureElem"`
	SoundElem           SoundElem   `mapstructure:"soundElem" json:"soundElem"`
	VideoElem           VideoElem   `mapstructure:"videoElem" json:"videoElem"`
	FileElem            FileElem    `mapstructure:"fileElem" json:"fileElem"`
	Ex                  string      `mapstructure:"ex" json:"ex"`
}
type MessageRevoked struct {
	RevokerID       string `mapstructure:"revokerID" json:"revokerID" validate:"required"`
	RevokerRole     int32  `mapstructure:"revokerRole" json:"revokerRole" validate:"required"`
	ClientMsgID     string `mapstructure:"clientMsgID" json:"clientMsgID" validate:"required"`
	RevokerNickname string `mapstructure:"revokerNickname" json:"revokerNickname"`
	SessionType     int32  `mapstructure:"sessionType" json:"sessionType" validate:"required"`
	Seq             uint32 `mapstructure:"seq" json:"seq" validate:"required"`
}
