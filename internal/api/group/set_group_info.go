package group

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pb "Open_IM/pkg/proto/group"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type paramsSetGroupInfo struct {
	GroupID      string `json:"groupId"  binding:"required"`
	GroupName    string `json:"groupName"`
	Notification string `json:"notification"`
	Introduction string `json:"introduction"`
	FaceUrl      string `json:"faceUrl"`
	OperationID  string `json:"operationID"  binding:"required"`
}

func SetGroupInfo(c *gin.Context) {
	log.Info("", "", "api set group info init...")

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pb.NewGroupClient(etcdConn)
	//defer etcdConn.Close()

	params := paramsSetGroupInfo{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pb.SetGroupInfoReq{
		GroupID:      params.GroupID,
		GroupName:    params.GroupName,
		Notification: params.Notification,
		Introduction: params.Introduction,
		FaceUrl:      params.FaceUrl,
		Token:        c.Request.Header.Get("token"),
		OperationID:  params.OperationID,
	}
	log.Info(req.Token, req.OperationID, "api set group info is server,params=%s", req.String())
	RpcResp, err := client.SetGroupInfo(context.Background(), req)
	if err != nil {
		log.Error(req.Token, req.OperationID, "call set group info rpc server failed,err=%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	log.InfoByArgs("call set group info rpc server success,args=%s", RpcResp.String())
	c.JSON(http.StatusOK, gin.H{"errCode": RpcResp.ErrorCode, "errMsg": RpcResp.ErrorMsg})
}
