package apiChat

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	pbChat "Open_IM/pkg/proto/chat"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"context"

	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type paramsUserSendMsg struct {
	SenderPlatformID int32  `json:"senderPlatformID" binding:"required"`
	SendID           string `json:"sendID" binding:"required"`
	SenderNickName   string `json:"senderNickName"`
	SenderFaceURL    string `json:"senderFaceUrl"`
	OperationID      string `json:"operationID" binding:"required"`
	Data             struct {
		SessionType int32                        `json:"sessionType" binding:"required"`
		MsgFrom     int32                        `json:"msgFrom" binding:"required"`
		ContentType int32                        `json:"contentType" binding:"required"`
		RecvID      string                       `json:"recvID" `
		GroupID     string                       `json:"groupID" `
		ForceList   []string                     `json:"forceList"`
		Content     []byte                       `json:"content" binding:"required"`
		Options     map[string]bool              `json:"options" `
		ClientMsgID string                       `json:"clientMsgID" binding:"required"`
		CreateTime  int64                        `json:"createTime" binding:"required"`
		OffLineInfo *open_im_sdk.OfflinePushInfo `json:"offlineInfo" `
	}
}

func newUserSendMsgReq(token string, params *paramsUserSendMsg) *pbChat.SendMsgReq {
	pbData := pbChat.SendMsgReq{
		Token:       token,
		OperationID: params.OperationID,
		MsgData: &open_im_sdk.MsgData{
			SendID:           params.SendID,
			RecvID:           params.Data.RecvID,
			GroupID:          params.Data.GroupID,
			ClientMsgID:      params.Data.ClientMsgID,
			SenderPlatformID: params.SenderPlatformID,
			SenderNickname:   params.SenderNickName,
			SenderFaceURL:    params.SenderFaceURL,
			SessionType:      params.Data.SessionType,
			MsgFrom:          params.Data.MsgFrom,
			ContentType:      params.Data.ContentType,
			Content:          params.Data.Content,
			CreateTime:       params.Data.CreateTime,
			Options:          params.Data.Options,
			OfflinePushInfo:  params.Data.OffLineInfo,
		},
	}
	return &pbData
}

func SendMsg(c *gin.Context) {
	params := paramsUserSendMsg{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		log.ErrorByKv("json unmarshal err", "", "err", err.Error(), "data", c.PostForm("data"))
		return
	}

	token := c.Request.Header.Get("token")

	log.InfoByKv("api call success to sendMsgReq", params.OperationID, "Parameters", params)

	pbData := newUserSendMsgReq(token, &params)
	log.Info("", "", "api SendMsg call start..., [data: %s]", pbData.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfflineMessageName)
	client := pbChat.NewChatClient(etcdConn)

	log.Info("", "", "api SendMsg call, api call rpc...")

	reply, err := client.SendMsg(context.Background(), pbData)
	if err != nil {
		log.NewError(params.OperationID, "SendMsg rpc failed, ", params, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 401, "errMsg": "SendMsg rpc failed, " + err.Error()})
		return
	}
	log.Info("", "", "api SendMsg call end..., [data: %s] [reply: %s]", pbData.String(), reply.String())

	c.JSON(http.StatusOK, gin.H{
		"errCode": reply.ErrCode,
		"errMsg":  reply.ErrMsg,
		"data": gin.H{
			"clientMsgID": reply.ClientMsgID,
			"serverMsgID": reply.ServerMsgID,
			"sendTime":    reply.SendTime,
		},
	})

}
