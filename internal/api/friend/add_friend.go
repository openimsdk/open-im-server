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

type paramsImportFriendReq struct {
	OperationID string   `json:"operationID" binding:"required"`
	UIDList     []string `json:"uidList" binding:"required"`
	OwnerUid    string   `json:"ownerUid" binding:"required"`
}

type paramsAddFriend struct {
	OperationID string `json:"operationID" binding:"required"`
	UID         string `json:"uid" binding:"required"`
	ReqMessage  string `json:"reqMessage"`
}

//
func ImportFriend(c *gin.Context) {
	log.Info("", "", "ImportFriend init ....")
	log.NewDebug("", "api importFriend start")
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := pbFriend.NewFriendClient(etcdConn)

	params := paramsImportFriendReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pbFriend.ImportFriendReq{
		UidList:     params.UIDList,
		OperationID: params.OperationID,
		OwnerUid:    params.OwnerUid,
		Token:       c.Request.Header.Get("token"),
	}
	log.NewDebug(req.OperationID, "args is ", req.String())
	RpcResp, err := client.ImportFriend(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "rpc importFriend failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "cImportFriend failed " + err.Error()})
		return
	}
	failedUidList := make([]string, 0)
	for _, v := range RpcResp.FailedUidList {
		failedUidList = append(failedUidList, v)
	}
	log.NewDebug(req.OperationID, "rpc importFriend success", RpcResp.CommonResp.ErrorMsg, RpcResp.CommonResp.ErrorCode, RpcResp.FailedUidList)
	c.JSON(http.StatusOK, gin.H{"errCode": RpcResp.CommonResp.ErrorCode, "errMsg": RpcResp.CommonResp.ErrorMsg, "failedUidList": failedUidList})
}

func AddFriend(c *gin.Context) {
	log.Info("", "", "api add friend init ....")

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := pbFriend.NewFriendClient(etcdConn)

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
	c.JSON(http.StatusOK, gin.H{
		"errCode": RpcResp.ErrorCode,
		"errMsg":  RpcResp.ErrorMsg,
	})
}
