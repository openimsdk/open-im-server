/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@tuoyun.net").
** time(2021/9/15 10:28).
 */
package manage

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbRelay "Open_IM/pkg/proto/relay"
	rpc "Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func DeleteUser(c *gin.Context) {
	params := api.DeleteUsersReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.DeleteUsersReq{}
	utils.CopyStructFields(req, &params)

	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(params.OperationID, "DeleteUser args ", req.String())
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)

	RpcResp, err := client.DeleteUsers(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "call delete users rpc server failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call delete users rpc server failed"})
		return
	}
	resp := api.DeleteUsersResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}, FailedUserIDList: RpcResp.FailedUserIDList}
	if len(RpcResp.FailedUserIDList) == 0 {
		resp.FailedUserIDList = []string{}
	}
	log.NewInfo(req.OperationID, "DeleteUser api return", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary 获取所有用户uid列表
// @Description 获取所有用户uid列表
// @Tags 用户相关
// @ID GetAllUsersUid
// @Accept json
// @Param token header string true "im token"
// @Param req body api.GetAllUsersUidReq true "请求体"
// @Produce json
// @Success 0 {object} api.GetAllUsersUidResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /user/get_all_users_uid [post]
func GetAllUsersUid(c *gin.Context) {
	params := api.GetAllUsersUidReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetAllUserIDReq{}
	utils.CopyStructFields(req, &params)

	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(params.OperationID, "GetAllUsersUid args ", req.String())
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)
	RpcResp, err := client.GetAllUserID(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "call GetAllUsersUid users rpc server failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call GetAllUsersUid users rpc server failed"})
		return
	}
	resp := api.GetAllUsersUidResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}, UserIDList: RpcResp.UserIDList}
	if len(RpcResp.UserIDList) == 0 {
		resp.UserIDList = []string{}
	}
	log.NewInfo(req.OperationID, "GetAllUsersUid api return", resp)
	c.JSON(http.StatusOK, resp)

}

// @Summary 检查列表账户注册状态，并且返回结果
// @Description 传入UserIDList检查列表账户注册状态，并且返回结果
// @Tags 用户相关
// @ID AccountCheck
// @Accept json
// @Param token header string true "im token"
// @Param req body api.AccountCheckReq true "请求体"
// @Produce json
// @Success 0 {object} api.AccountCheckResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /user/account_check [post]
func AccountCheck(c *gin.Context) {
	params := api.AccountCheckReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.AccountCheckReq{}
	utils.CopyStructFields(req, &params)

	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(params.OperationID, "AccountCheck args ", req.String())
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)

	RpcResp, err := client.AccountCheck(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "call AccountCheck users rpc server failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call AccountCheck users rpc server failed"})
		return
	}
	resp := api.AccountCheckResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}, ResultList: RpcResp.ResultList}
	if len(RpcResp.ResultList) == 0 {
		resp.ResultList = []*rpc.AccountCheckResp_SingleUserStatus{}
	}
	log.NewInfo(req.OperationID, "AccountCheck api return", resp)
	c.JSON(http.StatusOK, resp)
}

func GetUsersOnlineStatus(c *gin.Context) {
	params := api.GetUsersOnlineStatusReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pbRelay.GetUsersOnlineStatusReq{}
	utils.CopyStructFields(req, &params)

	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(params.OperationID, "GetUsersOnlineStatus args ", req.String())
	var wsResult []*pbRelay.GetUsersOnlineStatusResp_SuccessResult
	var respResult []*pbRelay.GetUsersOnlineStatusResp_SuccessResult
	flag := false
	grpcCons := getcdv3.GetDefaultGatewayConn4Unique(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), params.OperationID)
	for _, v := range grpcCons {
		client := pbRelay.NewRelayClient(v)
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
	log.NewInfo(params.OperationID, "call GetUsersOnlineStatus rpc server is success", wsResult)
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
	resp := api.GetUsersOnlineStatusResp{CommResp: api.CommResp{ErrCode: 0, ErrMsg: ""}, SuccessResult: respResult}
	if len(respResult) == 0 {
		resp.SuccessResult = []*pbRelay.GetUsersOnlineStatusResp_SuccessResult{}
	}
	log.NewInfo(req.OperationID, "GetUsersOnlineStatus api return", resp)
	c.JSON(http.StatusOK, resp)
}
