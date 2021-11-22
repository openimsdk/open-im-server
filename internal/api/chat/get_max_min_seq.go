package apiChat

import (
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbMsg "Open_IM/pkg/proto/chat"
	"Open_IM/pkg/utils"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// paramsUserNewestSeq struct
type paramsUserNewestSeq struct {
	ReqIdentifier int    `json:"reqIdentifier" binding:"required"`
	SendID        string `json:"sendID" binding:"required"`
	OperationID   string `json:"operationID" binding:"required"`
	MsgIncr       int    `json:"msgIncr" binding:"required"`
}

// resultUserNewestSeq struct
type resultUserNewestSeq struct {
	ErrCode       int32  `json:"errCode`
	ErrMsg        string `json:"errMsg"`
	MsgIncr       int    `json:"msgIncr"`
	ReqIdentifier int    `json:"reqIdentifier"`
	Data          struct {
		MaxSeq int64 `json:"maxSeq,omitempty"`
		MinSeq int64 `json:"minSeq,omitempty"`
	} `json:"data"`
}

// @Summary
// @Schemes
// @Description get latest message seq
// @Tags chat
// @Accept json
// @Produce json
// @Param body body apiChat.paramsUserNewestSeq true "user get latest seq params"
// @Param token header string true "token"
// @Success 200 {object} apiChat.resultUserNewestSeq
// @Failure 400 {object} user.result
// @Failure 500 {object} user.result
// @Router /chat/newest_seq [post]
func UserGetSeq(c *gin.Context) {
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
	pbData := pbMsg.GetMaxAndMinSeqReq{}
	pbData.UserID = params.SendID
	pbData.OperationID = params.OperationID
	grpcConn := getcdv3.GetOfflineMessageConn()

	if grpcConn == nil {
		log.ErrorByKv("get grpcConn err", pbData.OperationID, "args", params)
	}
	msgClient := pbMsg.NewChatClient(grpcConn)
	reply, err := msgClient.GetMaxAndMinSeq(context.Background(), &pbData)
	if err != nil {
		log.ErrorByKv("rpc call failed to getNewSeq", pbData.OperationID, "err", err, "pbData", pbData.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errCode":       reply.ErrCode,
		"errMsg":        reply.ErrMsg,
		"msgIncr":       params.MsgIncr,
		"reqIdentifier": params.ReqIdentifier,
		"data": gin.H{
			"maxSeq": reply.MaxSeq,
			"minSeq": reply.MinSeq,
		},
	})

}
