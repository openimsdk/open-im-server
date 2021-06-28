package group

import (
	"Open_IM/src/common/config"
	"Open_IM/src/common/log"
	pb "Open_IM/src/proto/group"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/skiffer-git/grpc-etcdv3/getcdv3"
	"net/http"
	"strings"
)

type paramsGetGroupInfo struct {
	GroupIDList []string `json:"groupIDList" binding:"required"`
	OperationID string   `json:"operationID" binding:"required"`
}

func GetGroupsInfo(c *gin.Context) {
	log.Info("", "", "api get groups info init ....")

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pb.NewGroupClient(etcdConn)
	defer etcdConn.Close()

	params := paramsGetGroupInfo{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pb.GetGroupsInfoReq{
		GroupIDList: params.GroupIDList,
		Token:       c.Request.Header.Get("token"),
		OperationID: params.OperationID,
	}
	log.Info(req.Token, req.OperationID, "get groups info is server,params=%s", req.String())
	RpcResp, err := client.GetGroupsInfo(context.Background(), req)
	if err != nil {
		log.Error(req.Token, req.OperationID, "call get groups info rpc server failed,err=%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	log.InfoByArgs("call get groups info rpc server success,info=%s", RpcResp.String())
	if RpcResp.ErrorCode == 0 {
		c.JSON(http.StatusOK, gin.H{
			"errCode": RpcResp.ErrorCode,
			"errMsg":  RpcResp.ErrorMsg,
			"data":    RpcResp.Data,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{"errCode": RpcResp.ErrorCode, "errMsg": RpcResp.ErrorMsg})
	}
}
