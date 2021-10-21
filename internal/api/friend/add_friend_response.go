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

type paramsAddFriendResponse struct {
	OperationID string `json:"operationID" binding:"required"`
	UID         string `json:"uid" binding:"required"`
	Flag        int32  `json:"flag" binding:"required"`
}

func AddFriendResponse(c *gin.Context) {
	log.Info("", "", fmt.Sprintf("api add friend response init ...."))

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := pbFriend.NewFriendClient(etcdConn)
	//defer etcdConn.Close()

	params := paramsAddFriendResponse{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pbFriend.AddFriendResponseReq{
		Uid:         params.UID,
		Flag:        params.Flag,
		OperationID: params.OperationID,
		Token:       c.Request.Header.Get("token"),
	}
	log.Info(req.Token, req.OperationID, "api add friend response is server:userID=%s", req.Uid)
	RpcResp, err := client.AddFriendResponse(context.Background(), req)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,call add_friend_response rpc server failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call add_friend_response rpc server failed"})
		return
	}
	log.InfoByArgs("call add friend response rpc server success,args=%s", RpcResp.String())
	c.JSON(http.StatusOK, gin.H{"errCode": RpcResp.ErrorCode, "errMsg": RpcResp.ErrorMsg})
	log.InfoByArgs("api add friend response success return,get args=%s,return args=%s", req.String(), RpcResp.String())
}
