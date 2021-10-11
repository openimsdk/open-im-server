package friend

import (
	pbFriend "Open_IM/pkg/proto/friend"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type paramsDeleteFriend struct {
	OperationID string `json:"operationID" binding:"required"`
	UID         string `json:"uid" binding:"required"`
}

func DeleteFriend(c *gin.Context) {
	log.Info("", "", fmt.Sprintf("api delete_friend init ...."))

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := pbFriend.NewFriendClient(etcdConn)
	//defer etcdConn.Close()

	params := paramsDeleteFriend{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pbFriend.DeleteFriendReq{
		Uid:         params.UID,
		OperationID: params.OperationID,
		Token:       c.Request.Header.Get("token"),
	}
	log.Info(req.Token, req.OperationID, "api delete_friend  is server:%s", req.Uid)
	RpcResp, err := client.DeleteFriend(context.Background(), req)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,call delete_friend rpc server failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call delete_friend rpc server failed"})
		return
	}
	log.InfoByArgs("call delete_friend rpc server,args=%s", RpcResp.String())
	resp := gin.H{"errCode": RpcResp.ErrorCode, "errMsg": RpcResp.ErrorMsg}
	c.JSON(http.StatusOK, resp)
	log.InfoByArgs("api delete_friend success return,get args=%s,return args=%s", req.String(), RpcResp.String())
}
