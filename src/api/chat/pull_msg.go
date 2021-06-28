package apiChat

import (
	"Open_IM/src/common/config"
	"Open_IM/src/common/log"
	"Open_IM/src/proto/chat"
	"Open_IM/src/utils"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/skiffer-git/grpc-etcdv3/getcdv3"
	"net/http"
	"strings"
)

type paramsUserPullMsg struct {
	ReqIdentifier *int   `json:"reqIdentifier" binding:"required"`
	SendID        string `json:"sendID" binding:"required"`
	OperationID   string `json:"operationID" binding:"required"`
	MsgIncr       *int   `json:"msgIncr" binding:"required"`
	Data          struct {
		SeqBegin *int64 `json:"seqBegin" binding:"required"`
		SeqEnd   *int64 `json:"seqEnd" binding:"required"`
	}
}

func UserPullMsg(c *gin.Context) {
	params := paramsUserPullMsg{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	token := c.Request.Header.Get("token")
	if !utils.VerifyToken(token, params.SendID) {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "token validate err"})
		return
	}

	pbData := pbChat.PullMessageReq{}
	pbData.UserID = params.SendID
	pbData.OperationID = params.OperationID
	pbData.SeqBegin = *params.Data.SeqBegin
	pbData.SeqEnd = *params.Data.SeqEnd
	grpcConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfflineMessageName)
	msgClient := pbChat.NewChatClient(grpcConn)
	reply, err := msgClient.PullMessage(context.Background(), &pbData)
	if err != nil {
		log.ErrorByKv("PullMessage error", pbData.OperationID, "err", err.Error())
		return
	}
	log.InfoByKv("rpc call success to pullMsgRep", pbData.OperationID, "ReplyArgs", reply.String(), "maxSeq", reply.GetMaxSeq(),
		"MinSeq", reply.GetMinSeq(), "singLen", len(reply.GetSingleUserMsg()), "groupLen", len(reply.GetGroupUserMsg()))

	msg := make(map[string]interface{})
	if v := reply.GetSingleUserMsg(); v != nil {
		msg["single"] = v
	} else {
		msg["single"] = []pbChat.GatherFormat{}
	}
	if v := reply.GetGroupUserMsg(); v != nil {
		msg["group"] = v
	} else {
		msg["group"] = []pbChat.GatherFormat{}
	}
	msg["maxSeq"] = reply.GetMaxSeq()
	msg["minSeq"] = reply.GetMinSeq()
	c.JSON(http.StatusOK, gin.H{
		"errCode":       reply.ErrCode,
		"errMsg":        reply.ErrMsg,
		"msgIncr":       *params.MsgIncr,
		"reqIdentifier": *params.ReqIdentifier,
		"data":          msg,
	})

}
