package apiChat

import (
	"Open_IM/src/common/config"
	"Open_IM/src/common/log"
	pbMsg "Open_IM/src/proto/chat"
	"Open_IM/src/utils"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/skiffer-git/grpc-etcdv3/getcdv3"
	"net/http"
	"strings"
)

type paramsUserNewestSeq struct {
	ReqIdentifier int    `json:"reqIdentifier" binding:"required"`
	SendID        string `json:"sendID" binding:"required"`
	OperationID   string `json:"operationID" binding:"required"`
	MsgIncr       int    `json:"msgIncr" binding:"required"`
}

func UserNewestSeq(c *gin.Context) {
	params := paramsUserNewestSeq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	token := c.Request.Header.Get("token")
	if !utils.VerifyToken(token, params.SendID) {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "token validate err"})
		return
	}
	pbData := pbMsg.GetNewSeqReq{}
	pbData.UserID = params.SendID
	pbData.OperationID = params.OperationID
	grpcConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfflineMessageName)
	if grpcConn == nil {
		log.ErrorByKv("get grpcConn err", pbData.OperationID, "args", params)
	}
	msgClient := pbMsg.NewChatClient(grpcConn)
	reply, err := msgClient.GetNewSeq(context.Background(), &pbData)
	if err != nil {
		log.ErrorByKv("rpc call failed to getNewSeq", pbData.OperationID, "err", err, "pbData", pbData.String())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errCode":       reply.ErrCode,
		"errMsg":        reply.ErrMsg,
		"msgIncr":       params.MsgIncr,
		"reqIdentifier": params.ReqIdentifier,
		"data": gin.H{
			"seq": reply.Seq,
		},
	})

}
