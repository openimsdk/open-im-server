package group

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	"Open_IM/pkg/proto/group"
	"Open_IM/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type paramsTransferGroupOwner struct {
	OperationID string `json:"operationID" binding:"required"`
	GroupID     string `json:"groupID" binding:"required"`
	UID         string `json:"uid" binding:"required"`
}

func newTransferGroupOwnerReq(params *paramsTransferGroupOwner) *group.TransferGroupOwnerReq {
	pbData := group.TransferGroupOwnerReq{
		OperationID: params.OperationID,
		GroupID:     params.GroupID,
		NewOwner:    params.UID,
	}
	return &pbData
}

func TransferGroupOwner(c *gin.Context) {
	log.Info("", "", "api TransferGroupOwner init ....")
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := group.NewGroupClient(etcdConn)
	//defer etcdConn.Close()

	params := paramsTransferGroupOwner{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	pbData := newTransferGroupOwnerReq(&params)

	token := c.Request.Header.Get("token")
	if claims, err := utils.ParseToken(token); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "token validate err"})
		return
	} else {
		pbData.OldOwner = claims.UID
	}

	log.Info("", "", "api TransferGroupOwner is server, [data: %s]", pbData.String())
	reply, err := client.TransferGroupOwner(context.Background(), pbData)
	if err != nil {
		log.Error("", "", "api TransferGroupOwner call rpc fail, [data: %s] [err: %s]", pbData.String(), err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	log.Info("", "", "api TransferGroupOwner call rpc success, [data: %s] [reply: %s]", pbData.String(), reply.String())

	c.JSON(http.StatusOK, gin.H{
		"errCode": reply.ErrCode,
		"errMsg":  reply.ErrMsg,
	})

}
