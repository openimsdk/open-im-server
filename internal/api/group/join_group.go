package group

import (
	pb "Open_IM/pkg/proto/group"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type paramsJoinGroup struct {
	GroupID     string `json:"groupID" binding:"required"`
	Message     string `json:"message"`
	OperationID string `json:"operationID" binding:"required"`
}

func JoinGroup(c *gin.Context) {
	log.Info("", "", "api join group init....")

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pb.NewGroupClient(etcdConn)
	//defer etcdConn.Close()

	params := paramsJoinGroup{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pb.JoinGroupReq{
		GroupID:     params.GroupID,
		Message:     params.Message,
		Token:       c.Request.Header.Get("token"),
		OperationID: params.OperationID,
	}
	log.Info(req.Token, req.OperationID, "api join group is server,params=%s", req.String())
	RpcResp, err := client.JoinGroup(context.Background(), req)
	if err != nil {
		log.Error(req.Token, req.OperationID, "call join group  rpc server failed,err=%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	log.InfoByArgs("call join group rpc server success,args=%s", RpcResp.String())
	c.JSON(http.StatusOK, gin.H{"errCode": RpcResp.ErrorCode, "errMsg": RpcResp.ErrorMsg})
}
