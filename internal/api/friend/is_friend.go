package friend

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbFriend "Open_IM/pkg/proto/friend"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type paramsIsFriend struct {
	OperationID string `json:"operationID" binding:"required"`
	ReceiveUid  string `json:"receive_uid"`
}

func IsFriend(c *gin.Context) {
	log.Info("", "", "api is friend init....")

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := pbFriend.NewFriendClient(etcdConn)
	//defer etcdConn.Close()

	params := paramsIsFriend{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pbFriend.IsFriendReq{
		OperationID: params.OperationID,
		ReceiveUid:  params.OperationID,
		Token:       c.Request.Header.Get("token"),
	}
	log.Info(req.Token, req.OperationID, "api is friend is server")
	RpcResp, err := client.IsFriend(context.Background(), req)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,call add friend rpc server failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call add friend rpc server failed"})
		return
	}
	log.InfoByArgs("call is friend rpc server success,args=%s", RpcResp.String())
	resp := gin.H{"errCode": RpcResp.ErrorCode, "errMsg": RpcResp.ErrorMsg, "isFriend": RpcResp.ShipType}
	c.JSON(http.StatusOK, resp)
	log.InfoByArgs("api is friend success return,get args=%s,return args=%s", req.String(), RpcResp.String())
}
