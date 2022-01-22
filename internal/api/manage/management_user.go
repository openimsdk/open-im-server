/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@tuoyun.net").
** time(2021/9/15 10:28).
 */
package manage

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbRelay "Open_IM/pkg/proto/relay"
	pbUser "Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type paramsDeleteUsers struct {
	OperationID   string   `json:"operationID" binding:"required"`
	DeleteUidList []string `json:"deleteUidList" binding:"required"`
}
type paramsGetAllUsersUid struct {
	OperationID string `json:"operationID" binding:"required"`
}
type paramsGetUsersOnlineStatus struct {
	OperationID string   `json:"operationID" binding:"required"`
	UserIDList  []string `json:"userIDList" binding:"required,lte=200"`
}
type paramsAccountCheck struct {
	OperationID string   `json:"operationID" binding:"required"`
	UserIDList  []string `json:"userIDList" binding:"required,lte=100"`
}

func DeleteUser(c *gin.Context) {
	params := paramsDeleteUsers{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.InfoByKv("DeleteUser req come here", params.OperationID, "DeleteUidList", params.DeleteUidList)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := pbUser.NewUserClient(etcdConn)
	//defer etcdConn.Close()

	req := &pbUser.DeleteUsersReq{
		OperationID:   params.OperationID,
		DeleteUidList: params.DeleteUidList,
		Token:         c.Request.Header.Get("token"),
	}
	RpcResp, err := client.DeleteUsers(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "call delete users rpc server failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call delete users rpc server failed"})
		return
	}
	failedUidList := make([]string, 0)
	for _, v := range RpcResp.FailedUidList {
		failedUidList = append(failedUidList, v)
	}
	log.InfoByKv("call delete user rpc server is success", params.OperationID, "resp args", RpcResp.String())
	resp := gin.H{"errCode": RpcResp.CommonResp.ErrorCode, "errMsg": RpcResp.CommonResp.ErrorMsg, "failedUidList": RpcResp.FailedUidList}
	c.JSON(http.StatusOK, resp)
}
func GetAllUsersUid(c *gin.Context) {
	params := paramsGetAllUsersUid{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.InfoByKv("GetAllUsersUid req come here", params.OperationID)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := pbUser.NewUserClient(etcdConn)
	//defer etcdConn.Close()

	req := &pbUser.GetAllUsersUidReq{
		OperationID: params.OperationID,
		Token:       c.Request.Header.Get("token"),
	}
	RpcResp, err := client.GetAllUsersUid(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error(), "uidList": []string{}})
		return
	}
	log.InfoByKv("call GetAllUsersUid rpc server is success", params.OperationID, "resp args", RpcResp.String())
	resp := gin.H{"errCode": RpcResp.CommonResp.ErrorCode, "errMsg": RpcResp.CommonResp.ErrorMsg, "uidList": RpcResp.UidList}
	c.JSON(http.StatusOK, resp)

}
func AccountCheck(c *gin.Context) {
	params := paramsAccountCheck{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.InfoByKv("AccountCheck req come here", params.OperationID, params.UserIDList)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := pbUser.NewUserClient(etcdConn)
	//defer etcdConn.Close()

	req := &pbUser.AccountCheckReq{
		OperationID: params.OperationID,
		Token:       c.Request.Header.Get("token"),
		UidList:     params.UserIDList,
	}
	RpcResp, err := client.AccountCheck(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	log.InfoByKv("call AccountCheck rpc server is success", params.OperationID, "resp args", RpcResp.String())
	resp := gin.H{"errCode": RpcResp.CommonResp.ErrorCode, "errMsg": RpcResp.CommonResp.ErrorMsg, "result": RpcResp.Result}
	c.JSON(http.StatusOK, resp)

}
func GetUsersOnlineStatus(c *gin.Context) {
	params := paramsGetUsersOnlineStatus{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	claims, err := token_verify.ParseToken(c.Request.Header.Get("token"))
	if err != nil {
		log.ErrorByKv("parse token failed", params.OperationID, "err", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 401, "errMsg": err.Error()})
		return
	}
	if !utils.IsContain(claims.UID, config.Config.Manager.AppManagerUid) {
		log.ErrorByKv(" Authentication failed", params.OperationID, "args", c)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 402, "errMsg": "not authorized"})
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
