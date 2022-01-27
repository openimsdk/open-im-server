/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@tuoyun.net").
** time(2021/9/15 15:23).
 */
package manage

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbChat "Open_IM/pkg/proto/chat"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"strings"
)

var validate *validator.Validate

func newUserSendMsgReq(params *ManagementSendMsgReq) *pbChat.SendMsgReq {
	var newContent string
	switch params.ContentType {
	case constant.Text:
		newContent = params.Content["text"].(string)
	case constant.Picture:
		fallthrough
	case constant.Custom:
		fallthrough
	case constant.Voice:
		fallthrough
	case constant.File:
		newContent = utils.StructToJsonString(params.Content)
	default:
	}
	options := make(map[string]bool, 2)
	if params.IsOnlineOnly {
		utils.SetSwitchFromOptions(options, constant.IsOfflinePush, false)
		utils.SetSwitchFromOptions(options, constant.IsHistory, false)
		utils.SetSwitchFromOptions(options, constant.IsPersistent, false)
	}
	pbData := pbChat.SendMsgReq{
		OperationID: params.OperationID,
		MsgData: &open_im_sdk.MsgData{
			SendID:           params.SendID,
			RecvID:           params.RecvID,
			GroupID:          params.GroupID,
			ClientMsgID:      utils.GetMsgID(params.SendID),
			SenderPlatformID: params.SenderPlatformID,
			SenderNickname:   params.SenderNickname,
			SenderFaceURL:    params.SenderFaceURL,
			SessionType:      params.SessionType,
			MsgFrom:          constant.SysMsgType,
			ContentType:      params.ContentType,
			Content:          []byte(newContent),
			//	ForceList:        params.ForceList,
			CreateTime:      utils.GetCurrentTimestampByMill(),
			Options:         options,
			OfflinePushInfo: params.OfflinePushInfo,
		},
	}
	return &pbData
}
func init() {
	validate = validator.New()
}

func ManagementSendMsg(c *gin.Context) {
	var data interface{}
	params := ManagementSendMsgReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		log.ErrorByKv("json unmarshal err", c.PostForm("operationID"), "err", err.Error(), "content", c.PostForm("content"))
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
	//case constant.AtText:
	//	data = AtElem{}
	//case constant.Merger:
	//	data =
	//case constant.Card:
	//case constant.Location:
	case constant.Custom:
		data = CustomElem{}
	//case constant.Revoke:
	//case constant.HasReadReceipt:
	//case constant.Typing:
	//case constant.Quote:
	default:
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 404, "errMsg": "contentType err"})
		log.ErrorByKv("contentType err", c.PostForm("operationID"), "content", c.PostForm("content"))
		return
	}
	if err := mapstructure.WeakDecode(params.Content, &data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 401, "errMsg": err.Error()})
		log.ErrorByKv("content to Data struct  err", "", "err", err.Error())
		return
	} else if err := validate.Struct(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 403, "errMsg": err.Error()})
		log.ErrorByKv("data args validate  err", "", "err", err.Error())
		return
	}

	token := c.Request.Header.Get("token")
	claims, err := token_verify.ParseToken(token)
	if err != nil {
		log.NewError(params.OperationID, "parse token failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "parse token failed", "sendTime": 0, "MsgID": ""})
	}
	if !utils.IsContain(claims.UID, config.Config.Manager.AppManagerUid) {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "not authorized", "sendTime": 0, "MsgID": ""})
		return

	}
	switch params.SessionType {
	case constant.SingleChatType:
		if len(params.RecvID) == 0 {
			log.NewError(params.OperationID, "recvID is a null string")
			c.JSON(http.StatusBadRequest, gin.H{"errCode": 405, "errMsg": "recvID is a null string", "sendTime": 0, "MsgID": ""})
		}
	case constant.GroupChatType:
		if len(params.GroupID) == 0 {
			log.NewError(params.OperationID, "groupID is a null string")
			c.JSON(http.StatusBadRequest, gin.H{"errCode": 405, "errMsg": "groupID is a null string", "sendTime": 0, "MsgID": ""})
		}

	}
	log.InfoByKv("Ws call success to ManagementSendMsgReq", params.OperationID, "Parameters", params)

	pbData := newUserSendMsgReq(&params)
	log.Info("", "", "api ManagementSendMsg call start..., [data: %s]", pbData.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfflineMessageName)
	client := pbChat.NewChatClient(etcdConn)

	log.Info("", "", "api ManagementSendMsg call, api call rpc...")

	reply, err := client.SendMsg(context.Background(), pbData)
	if err != nil {
		log.NewError(params.OperationID, "call delete UserSendMsg rpc server failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call UserSendMsg  rpc server failed"})
		return
	}
	log.Info("", "", "api ManagementSendMsg call end..., [data: %s] [reply: %s]", pbData.String(), reply.String())

	c.JSON(http.StatusOK, gin.H{
		"errCode":  reply.ErrCode,
		"errMsg":   reply.ErrMsg,
		"sendTime": reply.SendTime,
		"msgID":    reply.ClientMsgID,
	})

}

//
//type MergeElem struct {
//	Title        string       `json:"title"`
//	AbstractList []string     `json:"abstractList"`
//	MultiMessage []*MsgStruct `json:"multiMessage"`
//}
//
//type QuoteElem struct {
//	Text         string     `json:"text"`
//	QuoteMessage *MsgStruct `json:"quoteMessage"`
//}
type ManagementSendMsgReq struct {
	OperationID      string                       `json:"operationID" binding:"required"`
	SendID           string                       `json:"sendID" binding:"required"`
	RecvID           string                       `json:"recvID" `
	GroupID          string                       `json:"groupID" `
	SenderNickname   string                       `json:"senderNickname" `
	SenderFaceURL    string                       `json:"senderFaceURL" `
	SenderPlatformID int32                        `json:"senderPlatformID"`
	ForceList        []string                     `json:"forceList" `
	Content          map[string]interface{}       `json:"content" binding:"required"`
	ContentType      int32                        `json:"contentType" binding:"required"`
	SessionType      int32                        `json:"sessionType" binding:"required"`
	IsOnlineOnly     bool                         `json:"isOnlineOnly"`
	OfflinePushInfo  *open_im_sdk.OfflinePushInfo `json:"offlinePushInfo"`
}

type PictureBaseInfo struct {
	UUID   string `mapstructure:"uuid"`
	Type   string `mapstructure:"type" validate:"required"`
	Size   int64  `mapstructure:"size" validate:"required"`
	Width  int32  `mapstructure:"width" validate:"required"`
	Height int32  `mapstructure:"height" validate:"required"`
	Url    string `mapstructure:"url" validate:"required"`
}

type PictureElem struct {
	SourcePath      string          `mapstructure:"sourcePath"`
	SourcePicture   PictureBaseInfo `mapstructure:"sourcePicture" validate:"required"`
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
