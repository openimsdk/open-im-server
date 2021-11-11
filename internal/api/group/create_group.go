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

type paramsCreateGroupStruct struct {
	MemberList   []*pb.GroupAddMemberInfo `json:"memberList"`
	GroupName    string                   `json:"groupName"`
	Introduction string                   `json:"introduction"`
	Notification string                   `json:"notification"`
	FaceUrl      string                   `json:"faceUrl"`
	OperationID  string                   `json:"operationID" binding:"required"`
	Ex           string                   `json:"ex"`
}

func CreateGroup(c *gin.Context) {
	log.Info("", "", "api create group init ....")

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pb.NewGroupClient(etcdConn)
	//defer etcdConn.Close()

	params := paramsCreateGroupStruct{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pb.CreateGroupReq{
		MemberList:   params.MemberList,
		GroupName:    params.GroupName,
		Introduction: params.Introduction,
		Notification: params.Notification,
		FaceUrl:      params.FaceUrl,
		OperationID:  params.OperationID,
		Ex:           params.Ex,
		Token:        c.Request.Header.Get("token"),
	}
	log.Info(req.Token, req.OperationID, "api create group is server,params=%s", req.String())
	RpcResp, err := client.CreateGroup(context.Background(), req)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,call create group  rpc server failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	log.InfoByArgs("call create group  rpc server success,args=%s", RpcResp.String())
	if RpcResp.ErrorCode == 0 {
		resp := gin.H{"errCode": RpcResp.ErrorCode, "errMsg": RpcResp.ErrorMsg, "data": gin.H{"groupID": RpcResp.GroupID}}
		c.JSON(http.StatusOK, resp)
	} else {
		c.JSON(http.StatusOK, gin.H{"errCode": RpcResp.ErrorCode, "errMsg": RpcResp.ErrorMsg})
	}
	log.InfoByArgs("api create group success return,get args=%s,return args=%s", req.String(), RpcResp.String())
}
