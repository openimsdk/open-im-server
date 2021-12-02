package user

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbRelay "Open_IM/pkg/proto/relay"
	pbUser "Open_IM/pkg/proto/user"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type userInfo struct {
	UID    string `json:"uid"`
	Name   string `json:"name"`
	Icon   string `json:"icon"`
	Gender int32  `json:"gender"`
	Mobile string `json:"mobile"`
	Birth  string `json:"birth"`
	Email  string `json:"email"`
	Ex     string `json:"ex"`
}

type paramsGetUsersOnlineStatus struct {
	OperationID string   `json:"operationID" binding:"required"`
	UserIDList  []string `json:"userIDList" binding:"required,lte=200"`
	Secret      string   `json:"secret" binding:"required,max=32"`
}

func GetUsersOnlineStatus(c *gin.Context) {
	params := paramsGetUsersOnlineStatus{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError(params.OperationID, "bind json failed ", err.Error(), c)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	if params.Secret != config.Config.Secret {
		log.NewError(params.OperationID, "parse token failed ", params.Secret, config.Config.Secret)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 401, "errMsg": "secret failed"})
		return
	}

	req := &pbRelay.GetUsersOnlineStatusReq{
		OperationID: params.OperationID,
		UserIDList:  params.UserIDList,
	}
	var wsResult []*pbRelay.GetUsersOnlineStatusResp_SuccessResult
	var respResult []*pbRelay.GetUsersOnlineStatusResp_SuccessResult
	flag := false
	log.NewDebug(params.OperationID, "GetUsersOnlineStatus req come here", params.UserIDList)

	grpcCons := getcdv3.GetConn4Unique(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOnlineMessageRelayName)
	for _, v := range grpcCons {
		client := pbRelay.NewOnlineMessageRelayServiceClient(v)
		reply, err := client.GetUsersOnlineStatus(context.Background(), req)
		if err != nil {
			log.NewError(params.OperationID, "GetUsersOnlineStatus rpc  err", req.String(), err.Error())
			continue
		} else {
			if reply.ErrCode == 0 {
				wsResult = append(wsResult, reply.SuccessResult...)
			}
		}
	}
	log.NewDebug(params.OperationID, "call GetUsersOnlineStatus rpc server is success", wsResult)
	//Online data merge of each node
	for _, v1 := range params.UserIDList {
		flag = false
		temp := new(pbRelay.GetUsersOnlineStatusResp_SuccessResult)
		for _, v2 := range wsResult {
			if v2.UserID == v1 {
				flag = true
				temp.UserID = v1
				temp.Status = constant.OnlineStatus
				temp.DetailPlatformStatus = append(temp.DetailPlatformStatus, v2.DetailPlatformStatus...)
			}
		}
		if !flag {
			temp.UserID = v1
			temp.Status = constant.OfflineStatus
		}
		respResult = append(respResult, temp)
	}
	log.NewDebug(params.OperationID, "Finished merged data", respResult)
	resp := gin.H{"errCode": 0, "errMsg": "", "successResult": respResult}
	c.JSON(http.StatusOK, resp)
}

func GetUserInfo(c *gin.Context) {
	log.InfoByKv("api get userinfo init...", "")

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := pbUser.NewUserClient(etcdConn)
	//defer etcdConn.Close()

	params := paramsStruct{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pbUser.GetUserInfoReq{
		UserIDList:  params.UIDList,
		OperationID: params.OperationID,
		Token:       c.Request.Header.Get("token"),
	}
	log.InfoByKv("api get user info is server", c.PostForm("OperationID"), c.Request.Header.Get("token"))
	RpcResp, err := client.GetUserInfo(context.Background(), req)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,call get user info rpc server failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"errorCode": 500,
			"errorMsg":  "call  rpc server failed",
		})
		return
	}
	log.InfoByKv("call get user info rpc server success", params.OperationID)
	if RpcResp.ErrorCode == 0 {
		userInfoList := make([]userInfo, 0)
		for _, user := range RpcResp.Data {
			var ui userInfo
			ui.UID = user.Uid
			ui.Name = user.Name
			ui.Icon = user.Icon
			ui.Gender = user.Gender
			ui.Mobile = user.Mobile
			ui.Birth = user.Birth
			ui.Email = user.Email
			ui.Ex = user.Ex
			userInfoList = append(userInfoList, ui)
		}
		resp := gin.H{"errCode": RpcResp.ErrorCode, "errMsg": RpcResp.ErrorMsg, "data": userInfoList}
		c.JSON(http.StatusOK, resp)
	} else {
		resp := gin.H{"errCode": RpcResp.ErrorCode, "errMsg": RpcResp.ErrorMsg}
		c.JSON(http.StatusOK, resp)
	}
	log.InfoByKv("api get user info return success", params.OperationID, "args=%s", RpcResp.String())
}
