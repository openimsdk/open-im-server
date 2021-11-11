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

type paramsGetGroupInfo struct {
	GroupIDList []string `json:"groupIDList" binding:"required"`
	OperationID string   `json:"operationID" binding:"required"`
}

func GetGroupsInfo(c *gin.Context) {
	log.Info("", "", "api get groups info init ....")

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pb.NewGroupClient(etcdConn)
	//defer etcdConn.Close()

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
	log.InfoByArgs("call get groups info rpc server success", RpcResp.String())
	if RpcResp.ErrorCode == 0 {
		groupsInfo := make([]pb.GroupInfo, 0)
		for _, v := range RpcResp.Data {
			var groupInfo pb.GroupInfo
			groupInfo.GroupId = v.GroupId
			groupInfo.GroupName = v.GroupName
			groupInfo.Notification = v.Notification
			groupInfo.Introduction = v.Introduction
			groupInfo.FaceUrl = v.FaceUrl
			groupInfo.CreateTime = v.CreateTime
			groupInfo.OwnerId = v.OwnerId
			groupInfo.MemberCount = v.MemberCount

			groupsInfo = append(groupsInfo, groupInfo)
		}
		c.JSON(http.StatusOK, gin.H{
			"errCode": RpcResp.ErrorCode,
			"errMsg":  RpcResp.ErrorMsg,
			"data":    groupsInfo,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{"errCode": RpcResp.ErrorCode, "errMsg": RpcResp.ErrorMsg})
	}
}
