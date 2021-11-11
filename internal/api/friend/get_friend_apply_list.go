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

type paramsGetApplyList struct {
	OperationID string `json:"operationID" binding:"required"`
}
type UserInfo struct {
	UID        string `json:"uid"`
	Name       string `json:"name"`
	Icon       string `json:"icon"`
	Gender     int32  `json:"gender"`
	Mobile     string `json:"mobile"`
	Birth      string `json:"birth"`
	Email      string `json:"email"`
	Ex         string `json:"ex"`
	ReqMessage string `json:"reqMessage"`
	ApplyTime  string `json:"applyTime"`
	Flag       int32  `json:"flag"`
}

func GetFriendApplyList(c *gin.Context) {
	log.Info("", "", "api get_friend_apply_list init ....")

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := pbFriend.NewFriendClient(etcdConn)
	//defer etcdConn.Close()

	params := paramsGetApplyList{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pbFriend.GetFriendApplyReq{
		OperationID: params.OperationID,
		Token:       c.Request.Header.Get("token"),
	}
	log.Info(req.Token, req.OperationID, "api get friend apply list  is server")
	RpcResp, err := client.GetFriendApplyList(context.Background(), req)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,call get friend apply list rpc server failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call get friend apply list rpc server failed"})
		return
	}
	log.InfoByArgs("call get friend apply list rpc server success,args=%s", RpcResp.String())
	if RpcResp.ErrorCode == 0 {
		userInfoList := make([]UserInfo, 0)
		for _, applyUserinfo := range RpcResp.Data {
			var un UserInfo
			un.UID = applyUserinfo.Uid
			un.Name = applyUserinfo.Name
			un.Icon = applyUserinfo.Icon
			un.Gender = applyUserinfo.Gender
			un.Mobile = applyUserinfo.Mobile
			un.Birth = applyUserinfo.Birth
			un.Email = applyUserinfo.Email
			un.Ex = applyUserinfo.Ex
			un.Flag = applyUserinfo.Flag
			un.ApplyTime = applyUserinfo.ApplyTime
			un.ReqMessage = applyUserinfo.ReqMessage
			userInfoList = append(userInfoList, un)
		}
		resp := gin.H{"errCode": RpcResp.ErrorCode, "errMsg": RpcResp.ErrorMsg, "data": userInfoList}
		c.JSON(http.StatusOK, resp)
	} else {
		resp := gin.H{"errCode": RpcResp.ErrorCode, "errMsg": RpcResp.ErrorMsg}
		c.JSON(http.StatusOK, resp)
	}
	log.InfoByArgs("api get friend apply list success return,get args=%s,return args=%s", req.String(), RpcResp.String())
}

func GetSelfApplyList(c *gin.Context) {
	log.Info("", "", "api get self friend apply list init ....")

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := pbFriend.NewFriendClient(etcdConn)
	//defer etcdConn.Close()

	params := paramsGetApplyList{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pbFriend.GetFriendApplyReq{
		OperationID: params.OperationID,
		Token:       c.Request.Header.Get("token"),
	}
	log.Info(req.Token, req.OperationID, "api get self apply list  is server")
	RpcResp, err := client.GetSelfApplyList(context.Background(), req)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,call get self apply list rpc server failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call get self apply list rpc server failed"})
		return
	}
	log.InfoByArgs("call get self apply list rpc server success,args=%s", RpcResp.String())
	if RpcResp.ErrorCode == 0 {
		userInfoList := make([]UserInfo, 0)
		for _, selfApplyOtherUserinfo := range RpcResp.Data {
			var un UserInfo
			un.UID = selfApplyOtherUserinfo.Uid
			un.Name = selfApplyOtherUserinfo.Name
			un.Icon = selfApplyOtherUserinfo.Icon
			un.Gender = selfApplyOtherUserinfo.Gender
			un.Mobile = selfApplyOtherUserinfo.Mobile
			un.Birth = selfApplyOtherUserinfo.Birth
			un.Email = selfApplyOtherUserinfo.Email
			un.Ex = selfApplyOtherUserinfo.Ex
			un.Flag = selfApplyOtherUserinfo.Flag
			un.ApplyTime = selfApplyOtherUserinfo.ApplyTime
			un.ReqMessage = selfApplyOtherUserinfo.ReqMessage
			userInfoList = append(userInfoList, un)
		}
		resp := gin.H{"errCode": RpcResp.ErrorCode, "errMsg": RpcResp.ErrorMsg, "data": userInfoList}
		c.JSON(http.StatusOK, resp)
	} else {
		resp := gin.H{"errCode": RpcResp.ErrorCode, "errMsg": RpcResp.ErrorMsg}
		c.JSON(http.StatusOK, resp)
	}
	log.InfoByArgs("api get self apply list success return,get args=%s,return args=%s", req.String(), RpcResp.String())
}
