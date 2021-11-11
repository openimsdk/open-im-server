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

type paramsSetFriendComment struct {
	OperationID string `json:"operationID" binding:"required"`
	UID         string `json:"uid" binding:"required"`
	Comment     string `json:"comment"`
}

func SetFriendComment(c *gin.Context) {
	log.Info("", "", "api set friend comment init ....")

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := pbFriend.NewFriendClient(etcdConn)
	//defer etcdConn.Close()

	params := paramsSetFriendComment{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pbFriend.SetFriendCommentReq{
		Uid:         params.UID,
		OperationID: params.OperationID,
		Comment:     params.Comment,
		Token:       c.Request.Header.Get("token"),
	}
	log.Info(req.Token, req.OperationID, "api set friend comment is server")
	RpcResp, err := client.SetFriendComment(context.Background(), req)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,call set friend comment rpc server failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call set friend comment rpc server failed"})
		return
	}
	log.Info("", "", "call set friend comment rpc server success,args=%s", RpcResp.String())
	c.JSON(http.StatusOK, gin.H{"errCode": RpcResp.ErrorCode, "errMsg": RpcResp.ErrorMsg})
	log.Info("", "", "api set friend comment success return,get args=%s,return args=%s", req.String(), RpcResp.String())
}
