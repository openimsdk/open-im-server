package user

import (
	jsonData "Open_IM/internal/utils"
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	rpc "Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

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
	ok, req.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
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
			&open_im_sdk.PublicUserInfo{UserID: v.UserID, Nickname: v.Nickname, FaceURL: v.FaceURL, Gender: v.Gender})
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
	ok, req.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
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
	ok, req.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
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
