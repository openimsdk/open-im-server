package friend

import (
	"Open_IM/src/common/config"
	"Open_IM/src/common/log"
	pbFriend "Open_IM/src/proto/friend"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/skiffer-git/grpc-etcdv3/getcdv3"
	"net/http"
	"strings"
)

type paramsAddFriend struct {
	OperationID string `json:"operationID" binding:"required"`
	UID         string `json:"uid" binding:"required"`
	ReqMessage  string `json:"reqMessage"`
}

func AddFriend(c *gin.Context) {
	log.Info("", "", "api add friend init ....")

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := pbFriend.NewFriendClient(etcdConn)
	defer etcdConn.Close()

	params := paramsAddFriend{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pbFriend.AddFriendReq{
		Uid:         params.UID,
		OperationID: params.OperationID,
		ReqMessage:  params.ReqMessage,
		Token:       c.Request.Header.Get("token"),
	}
	log.Info(req.Token, req.OperationID, "api add friend is server")
	RpcResp, err := client.AddFriend(context.Background(), req)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,call add friend rpc server failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call add friend rpc server failed"})
		return
	}
	log.InfoByArgs("call add friend rpc server success,args=%s", RpcResp.String())
	resp := gin.H{"errCode": RpcResp.ErrorCode, "errMsg": RpcResp.ErrorMsg}
	c.JSON(http.StatusOK, resp)
	log.InfoByArgs("api add friend success return,get args=%s,return args=%s", req.String(), RpcResp.String())
}
