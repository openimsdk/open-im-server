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
	getUserInfoReq := &rpc.GetUserInfoReq{}
	getUserInfoReq.OperationID = params.OperationID
	var ok bool
	ok, getUserInfoReq.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), getUserInfoReq.OperationID)
	if !ok {
		log.NewError(getUserInfoReq.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(params.OperationID, "GetUserInfo args ", getUserInfoReq.String())
	reqCacheGetUserInfo := &cacheRpc.GetUserInfoReq{}
	utils.CopyStructFields(reqCacheGetUserInfo, &params)
	var userInfoList []*open_im_sdk.UserInfo
	var publicUserInfoList []*open_im_sdk.PublicUserInfo
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCacheName)
	cacheClient := cacheRpc.NewCacheClient(etcdConn)
	cacheResp, err := cacheClient.GetUserInfo(context.Background(), reqCacheGetUserInfo)
	if err != nil {
		log.NewError(getUserInfoReq.OperationID, utils.GetSelfFuncName(), "GetUserInfo failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed: " + err.Error()})
		return
	}
	if cacheResp.CommonResp.ErrCode != 0 {
		log.NewError(getUserInfoReq.OperationID, utils.GetSelfFuncName(), "GetUserInfo failed", cacheResp.CommonResp)
		resp := api.GetUsersInfoResp{CommResp: api.CommResp{ErrCode: cacheResp.CommonResp.ErrCode, ErrMsg: cacheResp.CommonResp.ErrMsg}}
		resp.Data = []map[string]interface{}{}
		log.NewInfo(getUserInfoReq.OperationID, "GetUserInfo api return ", resp)
		c.JSON(http.StatusOK, resp)
		return
	}
	log.NewInfo(getUserInfoReq.OperationID, utils.GetSelfFuncName(), "cacheResp:", cacheResp.String())
	userInfoList = cacheResp.UserInfoList
	var needCacheUserIDList []string
	for _, userID := range reqCacheGetUserInfo.UserIDList {
		isGetUserInfoFromCache := false
		for _, cacheUser := range userInfoList {
			if cacheUser.UserID == userID {
				isGetUserInfoFromCache = true
			}
		}
		if !isGetUserInfoFromCache {
			needCacheUserIDList = append(needCacheUserIDList, userID)
		}
	}
	if len(needCacheUserIDList) == 0 {
		log.NewInfo(getUserInfoReq.OperationID, utils.GetSelfFuncName(), "get all userInfo from cache success")
		for _, v := range userInfoList {
			publicUserInfoList = append(publicUserInfoList,
				&open_im_sdk.PublicUserInfo{UserID: v.UserID, Nickname: v.Nickname, FaceURL: v.FaceURL, Gender: v.Gender, Ex: v.Ex})
		}
		resp := api.GetUsersInfoResp{CommResp: api.CommResp{ErrCode: cacheResp.CommonResp.ErrCode, ErrMsg: cacheResp.CommonResp.ErrMsg}, UserInfoList: publicUserInfoList}
		resp.Data = jsonData.JsonDataList(resp.UserInfoList)
		log.NewInfo(getUserInfoReq.OperationID, "GetUserInfo api return ", resp)
		c.JSON(http.StatusOK, resp)
		return
	}

	log.NewDebug(getUserInfoReq.OperationID, utils.GetSelfFuncName(), "need cache user list", needCacheUserIDList)
	getUserInfoReq.UserIDList = needCacheUserIDList
	etcdConn = getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := rpc.NewUserClient(etcdConn)
	rpcResp, err := client.GetUserInfo(context.Background(), getUserInfoReq)
	if err != nil {
		log.NewError(getUserInfoReq.OperationID, "GetUserInfo failed ", err.Error(), getUserInfoReq.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed" + err.Error()})
		return
	}
	if rpcResp.CommonResp.ErrCode != 0 {
		log.NewError(getUserInfoReq.OperationID, utils.GetSelfFuncName(), "GetUserInfo failed", cacheResp.CommonResp)
		resp := api.GetUsersInfoResp{CommResp: api.CommResp{ErrCode: cacheResp.CommonResp.ErrCode, ErrMsg: cacheResp.CommonResp.ErrMsg}}
		resp.Data = []map[string]interface{}{}
		log.NewInfo(getUserInfoReq.OperationID, "GetUserInfo api return ", resp)
		c.JSON(http.StatusOK, resp)
		return
	}
	userInfoList = append(userInfoList, rpcResp.UserInfoList...)
	cacheUpdateUserInfoReq := &cacheRpc.UpdateUserInfoReq{
		UserInfoList: rpcResp.UserInfoList,
		OperationID:  getUserInfoReq.OperationID,
	}
	_, err = cacheClient.UpdateUserInfo(context.Background(), cacheUpdateUserInfoReq)
	if err != nil {
		log.NewError(getUserInfoReq.OperationID, "GetUserInfo failed ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed:" + err.Error()})
		return
	}
	userInfoList = rpcResp.UserInfoList
	for _, v := range userInfoList {
		publicUserInfoList = append(publicUserInfoList,
			&open_im_sdk.PublicUserInfo{UserID: v.UserID, Nickname: v.Nickname, FaceURL: v.FaceURL, Gender: v.Gender, Ex: v.Ex})
	}
	resp := api.GetUsersInfoResp{CommResp: api.CommResp{ErrCode: rpcResp.CommonResp.ErrCode, ErrMsg: rpcResp.CommonResp.ErrMsg}, UserInfoList: publicUserInfoList}
	resp.Data = jsonData.JsonDataList(resp.UserInfoList)
	log.NewInfo(getUserInfoReq.OperationID, "GetUserInfo api return ", resp)
	c.JSON(http.StatusOK, resp)
}

func GetUserFriendFromCache(c *gin.Context) {

}

func GetBlackListFromCache(c *gin.Context) {

}

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
	ok, req.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(params.OperationID, "GetUserInfo args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
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
	ok, req.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(params.OperationID, "UpdateUserInfo args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
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

func GetSelfUserInfo(c *gin.Context) {
	params := api.GetSelfUserInfoReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": http.StatusBadRequest, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetUserInfoReq{}

	utils.CopyStructFields(req, &params)
	var ok bool
	ok, req.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	req.UserIDList = append(req.UserIDList, req.OpUserID)
	log.NewInfo(params.OperationID, "GetUserInfo args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
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

func GetUsersOnlineStatus(c *gin.Context) {
	params := api.GetUsersOnlineStatusReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &pbRelay.GetUsersOnlineStatusReq{}
	utils.CopyStructFields(req, &params)
	var ok bool
	ok, req.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
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
