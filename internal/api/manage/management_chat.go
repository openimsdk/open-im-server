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
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbChat "Open_IM/pkg/proto/chat"
	"Open_IM/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"strings"
)

var validate *validator.Validate

type paramsManagementSendMsg struct {
	OperationID    string                 `json:"operationID" binding:"required"`
	SendID         string                 `json:"sendID" binding:"required"`
	RecvID         string                 `json:"recvID" binding:"required"`
	SenderNickName string                 `json:"senderNickName" `
	SenderFaceURL  string                 `json:"senderFaceURL" `
	ForceList      []string               `json:"forceList" `
	Content        map[string]interface{} `json:"content" binding:"required"`
	ContentType    int32                  `json:"contentType" binding:"required"`
	SessionType    int32                  `json:"sessionType" binding:"required"`
}

func newUserSendMsgReq(params *paramsManagementSendMsg) *pbChat.UserSendMsgReq {
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
	pbData := pbChat.UserSendMsgReq{
		ReqIdentifier:  constant.WSSendMsg,
		SendID:         params.SendID,
		SenderNickName: params.SenderNickName,
		SenderFaceURL:  params.SenderFaceURL,
		OperationID:    params.OperationID,
		PlatformID:     0,
		SessionType:    params.SessionType,
		MsgFrom:        constant.UserMsgType,
		ContentType:    params.ContentType,
		RecvID:         params.RecvID,
		ForceList:      params.ForceList,
		Content:        newContent,
		ClientMsgID:    utils.GetMsgID(params.SendID),
	}
	return &pbData
}
func init() {
	validate = validator.New()
}
func ManagementSendMsg(c *gin.Context) {
	var data interface{}
	params := paramsManagementSendMsg{}
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
	case constant.Custom:
		data = CustomElem{}
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
	claims, err := utils.ParseToken(token)
	if err != nil {
		log.NewError(params.OperationID, "parse token failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "parse token failed", "sendTime": 0, "MsgID": ""})
	}
	if !utils.IsContain(claims.UID, config.Config.Manager.AppManagerUid) {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "not authorized", "sendTime": 0, "MsgID": ""})
		return

	}
	log.InfoByKv("Ws call success to ManagementSendMsgReq", params.OperationID, "Parameters", params)

	pbData := newUserSendMsgReq(&params)
	log.Info("", "", "api ManagementSendMsg call start..., [data: %s]", pbData.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfflineMessageName)
	client := pbChat.NewChatClient(etcdConn)

	log.Info("", "", "api ManagementSendMsg call, api call rpc...")

	reply, _ := client.UserSendMsg(context.Background(), pbData)
	log.Info("", "", "api ManagementSendMsg call end..., [data: %s] [reply: %s]", pbData.String(), reply.String())

	c.JSON(http.StatusOK, gin.H{
		"errCode":  reply.ErrCode,
		"errMsg":   reply.ErrMsg,
		"sendTime": reply.SendTime,
		"msgID":    reply.ClientMsgID,
	})

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

//type MergeElem struct {
//	Title        string       `json:"title"`
//	AbstractList []string     `json:"abstractList"`
//	MultiMessage []*MsgStruct `json:"multiMessage"`
//}
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

//type QuoteElem struct {
//	Text         string     `json:"text"`
//	QuoteMessage *MsgStruct `json:"quoteMessage"`
//}
