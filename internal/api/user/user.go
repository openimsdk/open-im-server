package user

import (
	jsonData "Open_IM/internal/utils"
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	cacheRpc "Open_IM/pkg/proto/cache"
	pbRelay "Open_IM/pkg/proto/relay"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	rpc "Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func GetUsersInfoFromCache(c *gin.Context) {
	params := api.GetUsersInfoReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": http.StatusBadRequest, "errMsg": err.Error()})
		return
	}
	log.NewInfo(params.OperationID, "GetUsersInfoFromCache req: ", params)
	req := &rpc.GetUserInfoReq{}
	utils.CopyStructFields(req, &params)
	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)
	RpcResp, err := client.GetUserInfo(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "GetUserInfo failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	var publicUserInfoList []*open_im_sdk.PublicUserInfo
	for _, v := range RpcResp.UserInfoList {
		publicUserInfoList = append(publicUserInfoList,
			&open_im_sdk.PublicUserInfo{UserID: v.UserID, Nickname: v.Nickname, FaceURL: v.FaceURL, Gender: v.Gender, Ex: v.Ex})
	}
	resp := api.GetUsersInfoResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}, UserInfoList: publicUserInfoList}
	resp.Data = jsonData.JsonDataList(resp.UserInfoList)
	log.NewInfo(req.OperationID, "GetUserInfo api return ", resp)
	c.JSON(http.StatusOK, resp)
}

func GetFriendIDListFromCache(c *gin.Context) {
	var (
		req    api.GetFriendIDListFromCacheReq
		resp   api.GetFriendIDListFromCacheResp
		reqPb  cacheRpc.GetFriendIDListFromCacheReq
		respPb *cacheRpc.GetFriendIDListFromCacheResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": http.StatusBadRequest, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), req)
	reqPb.OperationID = req.OperationID
	var ok bool
	var errInfo string
	ok, reqPb.UserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCacheName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := cacheRpc.NewCacheClient(etcdConn)
	respPb, err := client.GetFriendIDListFromCache(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetFriendIDListFromCache", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed:" + err.Error()})
		return
	}
	resp.UserIDList = respPb.UserIDList
	resp.CommResp = api.CommResp{ErrMsg: respPb.CommonResp.ErrMsg, ErrCode: respPb.CommonResp.ErrCode}
	c.JSON(http.StatusOK, resp)
}

func GetBlackIDListFromCache(c *gin.Context) {
	var (
		req    api.GetBlackIDListFromCacheReq
		resp   api.GetBlackIDListFromCacheResp
		reqPb  cacheRpc.GetBlackIDListFromCacheReq
		respPb *cacheRpc.GetBlackIDListFromCacheResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": http.StatusBadRequest, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)
	reqPb.OperationID = req.OperationID
	var ok bool
	var errInfo string
	ok, reqPb.UserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCacheName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := cacheRpc.NewCacheClient(etcdConn)
	respPb, err := client.GetBlackIDListFromCache(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetFriendIDListFromCache", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed:" + err.Error()})
		return
	}
	resp.UserIDList = respPb.UserIDList
	resp.CommResp = api.CommResp{ErrMsg: respPb.CommonResp.ErrMsg, ErrCode: respPb.CommonResp.ErrCode}
	c.JSON(http.StatusOK, resp)
}

// @Summary 获取用户信息
// @Description 根据用户列表批量获取用户信息
// @Tags 用户相关
// @ID GetUsersInfo
// @Accept json
// @Param token header string true "im token"
// @Param req body api.GetUsersInfoReq true "请求体"
// @Produce json
// @Success 0 {object} api.GetUsersInfoResp{Data=[]open_im_sdk.PublicUserInfo}
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /user/get_users_info [post]
func GetUsersInfo(c *gin.Context) {
	params := api.GetUsersInfoReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": http.StatusBadRequest, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetUserInfoReq{}
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

	log.NewInfo(params.OperationID, "GetUserInfo args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)
	RpcResp, err := client.GetUserInfo(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "GetUserInfo failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	var publicUserInfoList []*open_im_sdk.PublicUserInfo
	for _, v := range RpcResp.UserInfoList {
		publicUserInfoList = append(publicUserInfoList,
			&open_im_sdk.PublicUserInfo{UserID: v.UserID, Nickname: v.Nickname, FaceURL: v.FaceURL, Gender: v.Gender, Ex: v.Ex})
	}

	resp := api.GetUsersInfoResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}, UserInfoList: publicUserInfoList}
	resp.Data = jsonData.JsonDataList(resp.UserInfoList)
	log.NewInfo(req.OperationID, "GetUserInfo api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary 修改用户信息
// @Description 修改用户信息 userID faceURL等
// @Tags 用户相关
// @ID UpdateUserInfo
// @Accept json
// @Param token header string true "im token"
// @Param req body api.UpdateSelfUserInfoReq true "请求体"
// @Produce json
// @Success 0 {object} api.UpdateUserInfoResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /user/update_user_info [post]
func UpdateUserInfo(c *gin.Context) {
	params := api.UpdateSelfUserInfoReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.UpdateUserInfoReq{UserInfo: &open_im_sdk.UserInfo{}}
	utils.CopyStructFields(req.UserInfo, &params)
	req.OperationID = params.OperationID
	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.NewInfo(params.OperationID, "UpdateUserInfo args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)
	RpcResp, err := client.UpdateUserInfo(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "UpdateUserInfo failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	resp := api.UpdateUserInfoResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}}
	log.NewInfo(req.OperationID, "UpdateUserInfo api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary 设置全局免打扰
// @Description 设置全局免打扰
// @Tags 用户相关
// @ID SetGlobalRecvMessageOpt
// @Accept json
// @Param token header string true "im token"
// @Param req body api.SetGlobalRecvMessageOptReq true "globalRecvMsgOpt为全局免打扰设置0为关闭 1为开启"
// @Produce json
// @Success 0 {object} api.SetGlobalRecvMessageOptResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /user/set_global_msg_recv_opt [post]
func SetGlobalRecvMessageOpt(c *gin.Context) {
	params := api.SetGlobalRecvMessageOptReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.SetGlobalRecvMessageOptReq{}
	utils.CopyStructFields(req, &params)
	req.OperationID = params.OperationID
	var ok bool
	var errInfo string
	ok, req.UserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	log.NewInfo(params.OperationID, "SetGlobalRecvMessageOpt args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)
	RpcResp, err := client.SetGlobalRecvMessageOpt(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "SetGlobalRecvMessageOpt failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	resp := api.UpdateUserInfoResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}}
	log.NewInfo(req.OperationID, "SetGlobalRecvMessageOpt api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary 获取自己的信息
// @Description 传入ID获取自己的信息
// @Tags 用户相关
// @ID GetSelfUserInfo
// @Accept json
// @Param token header string true "im token"
// @Param req body api.GetSelfUserInfoReq true "请求体"
// @Produce json
// @Success 0 {object} api.GetSelfUserInfoResp{data=open_im_sdk.UserInfo}
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /user/get_self_user_info [post]
func GetSelfUserInfo(c *gin.Context) {
	params := api.GetSelfUserInfoReq{}
	if err := c.BindJSON(&params); err != nil {
		errMsg := "BindJSON failed " + err.Error()
		log.NewError("0", "BindJSON failed ", errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	req := &rpc.GetUserInfoReq{}

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

	req.UserIDList = append(req.UserIDList, req.OpUserID)
	log.NewInfo(params.OperationID, "GetUserInfo args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewUserClient(etcdConn)
	RpcResp, err := client.GetUserInfo(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "GetUserInfo failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	if len(RpcResp.UserInfoList) == 1 {
		resp := api.GetSelfUserInfoResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}, UserInfo: RpcResp.UserInfoList[0]}
		resp.Data = jsonData.JsonDataOne(resp.UserInfo)
		log.NewInfo(req.OperationID, "GetUserInfo api return ", resp)
		c.JSON(http.StatusOK, resp)
	} else {
		resp := api.GetSelfUserInfoResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}}
		log.NewInfo(req.OperationID, "GetUserInfo api return ", resp)
		c.JSON(http.StatusOK, resp)
	}

}

// @Summary 获取用户在线状态
// @Description 获取用户在线状态
// @Tags 用户相关
// @ID GetUsersOnlineStatus
// @Accept json
// @Param token header string true "im token"
// @Param req body api.GetUsersOnlineStatusReq true "请求体"
// @Produce json
// @Success 0 {object} api.GetUsersOnlineStatusResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /user/get_users_online_status [post]
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
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	if len(config.Config.Manager.AppManagerUid) == 0 {
		log.NewError(req.OperationID, "Manager == 0")
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "Manager == 0"})
		return
	}
	req.OpUserID = config.Config.Manager.AppManagerUid[0]

	log.NewInfo(params.OperationID, "GetUsersOnlineStatus args ", req.String())
	var wsResult []*pbRelay.GetUsersOnlineStatusResp_SuccessResult
	var respResult []*pbRelay.GetUsersOnlineStatusResp_SuccessResult
	flag := false
	grpcCons := getcdv3.GetConn4Unique(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImRelayName)
	for _, v := range grpcCons {
		log.Debug(params.OperationID, "get node ", *v, v.Target())
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
