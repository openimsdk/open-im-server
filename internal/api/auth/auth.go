package apiAuth

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	rpc "Open_IM/pkg/proto/auth"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// @Summary 用户注册
// @Description 用户注册
// @Tags 鉴权认证
// @ID UserRegister
// @Accept json
// @Param req body api.UserRegisterReq true "secret为openIM密钥, 详细见服务端config.yaml secret字段 <br> platform为平台ID <br> ex为拓展字段 <br> gender为性别, 0为女, 1为男"
// @Produce json
// @Success 0 {object} api.UserRegisterResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /auth/user_register [post]
func UserRegister(c *gin.Context) {
	params := api.UserRegisterReq{}
	if err := c.BindJSON(&params); err != nil {
		errMsg := " BindJSON failed " + err.Error()
		log.NewError("0", errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}

	if params.Secret != config.Config.Secret {
		errMsg := " params.Secret != config.Config.Secret "
		log.NewError(params.OperationID, errMsg, params.Secret, config.Config.Secret)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 401, "errMsg": errMsg})
		return
	}
	req := &rpc.UserRegisterReq{UserInfo: &open_im_sdk.UserInfo{}}
	utils.CopyStructFields(req.UserInfo, &params)
	//copier.Copy(req.UserInfo, &params)
	req.OperationID = params.OperationID
	log.NewInfo(req.OperationID, "UserRegister args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAuthName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + " getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewAuthClient(etcdConn)
	reply, err := client.UserRegister(context.Background(), req)
	if err != nil {
		errMsg := req.OperationID + " " + "UserRegister failed " + err.Error() + req.String()
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	if reply.CommonResp.ErrCode != 0 {
		errMsg := req.OperationID + " " + " UserRegister failed " + reply.CommonResp.ErrMsg + req.String()
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	pbDataToken := &rpc.UserTokenReq{Platform: params.Platform, FromUserID: params.UserID, OperationID: params.OperationID}
	replyToken, err := client.UserToken(context.Background(), pbDataToken)
	if err != nil {
		errMsg := req.OperationID + " " + " client.UserToken failed " + err.Error() + pbDataToken.String()
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	resp := api.UserRegisterResp{CommResp: api.CommResp{ErrCode: replyToken.CommonResp.ErrCode, ErrMsg: replyToken.CommonResp.ErrMsg},
		UserToken: api.UserTokenInfo{UserID: req.UserInfo.UserID, Token: replyToken.Token, ExpiredTime: replyToken.ExpiredTime}}
	log.NewInfo(req.OperationID, "UserRegister return ", resp)
	c.JSON(http.StatusOK, resp)

}

// @Summary 用户登录
// @Description 获取用户的token
// @Tags 鉴权认证
// @ID UserToken
// @Accept json
// @Param req body api.UserTokenReq true "secret为openIM密钥, 详细见服务端config.yaml secret字段 <br> platform为平台ID"
// @Produce json
// @Success 0 {object} api.UserTokenResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /auth/user_token [post]
func UserToken(c *gin.Context) {
	params := api.UserTokenReq{}
	if err := c.BindJSON(&params); err != nil {
		errMsg := " BindJSON failed " + err.Error()
		log.NewError("0", errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}

	if params.Secret != config.Config.Secret {
		errMsg := params.OperationID + " params.Secret != config.Config.Secret "
		log.NewError(params.OperationID, "params.Secret != config.Config.Secret", params.Secret, config.Config.Secret)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 401, "errMsg": errMsg})
		return
	}
	req := &rpc.UserTokenReq{Platform: params.Platform, FromUserID: params.UserID, OperationID: params.OperationID}
	log.NewInfo(req.OperationID, "UserToken args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAuthName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + " getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewAuthClient(etcdConn)
	reply, err := client.UserToken(context.Background(), req)
	if err != nil {
		errMsg := req.OperationID + " UserToken failed " + err.Error() + " req: " + req.String()
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	resp := api.UserTokenResp{CommResp: api.CommResp{ErrCode: reply.CommonResp.ErrCode, ErrMsg: reply.CommonResp.ErrMsg},
		UserToken: api.UserTokenInfo{UserID: req.FromUserID, Token: reply.Token, ExpiredTime: reply.ExpiredTime}}
	log.NewInfo(req.OperationID, "UserToken return ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary 解析当前用户token
// @Description 解析当前用户token(token在请求头中传入)
// @Tags 鉴权认证
// @ID ParseToken
// @Accept json
// @Param token header string true "im token"
// @Param req body api.ParseTokenReq true "secret为openIM密钥, 详细见服务端config.yaml secret字段<br>platform为平台ID"
// @Produce json
// @Success 0 {object} api.ParseTokenResp{Data=api.ExpireTime}
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /auth/parse_token [post]
func ParseToken(c *gin.Context) {
	params := api.ParseTokenReq{}
	if err := c.BindJSON(&params); err != nil {
		errMsg := " BindJSON failed " + err.Error()
		log.NewError("0", errMsg)
		c.JSON(http.StatusOK, gin.H{"errCode": 1001, "errMsg": errMsg})
		return
	}

	var ok bool
	var errInfo string
	var expireTime int64
	ok, _, errInfo, expireTime = token_verify.GetUserIDFromTokenExpireTime(c.Request.Header.Get("token"), params.OperationID)
	if !ok {
		errMsg := params.OperationID + " " + "GetUserIDFromTokenExpireTime failed " + errInfo
		log.NewError(params.OperationID, errMsg)
		c.JSON(http.StatusOK, gin.H{"errCode": 1001, "errMsg": errMsg})
		return
	}

	resp := api.ParseTokenResp{CommResp: api.CommResp{ErrCode: 0, ErrMsg: ""}, ExpireTime: api.ExpireTime{ExpireTimeSeconds: uint32(expireTime)}}
	resp.Data = structs.Map(&resp.ExpireTime)
	log.NewInfo(params.OperationID, "ParseToken return ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary 强制登出
// @Description 对应的平台强制登出
// @Tags 鉴权认证
// @ID ForceLogout
// @Accept json
// @Param token header string true "im token"
// @Param req body api.ForceLogoutReq true "platform为平台ID <br> fromUserID为要执行强制登出的用户ID"
// @Produce json
// @Success 0 {object} api.ForceLogoutResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /auth/force_logout [post]
func ForceLogout(c *gin.Context) {
	params := api.ForceLogoutReq{}
	if err := c.BindJSON(&params); err != nil {
		errMsg := " BindJSON failed " + err.Error()
		log.NewError("0", errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}

	req := &rpc.ForceLogoutReq{}
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

	log.NewInfo(req.OperationID, "ForceLogout args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAuthName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + " getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewAuthClient(etcdConn)
	reply, err := client.ForceLogout(context.Background(), req)
	if err != nil {
		errMsg := req.OperationID + " UserToken failed " + err.Error() + req.String()
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	resp := api.ForceLogoutResp{CommResp: api.CommResp{ErrCode: reply.CommonResp.ErrCode, ErrMsg: reply.CommonResp.ErrMsg}}
	log.NewInfo(params.OperationID, utils.GetSelfFuncName(), " return ", resp)
	c.JSON(http.StatusOK, resp)
}
