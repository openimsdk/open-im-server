package apiAuth

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	rpc "Open_IM/pkg/proto/auth"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func UserRegister(c *gin.Context) {
	params := api.UserRegisterReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	if params.Secret != config.Config.Secret {
		log.NewError(params.OperationID, "params.Secret != config.Config.Secret", params.Secret, config.Config.Secret)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 401, "errMsg": "not authorized"})
		return
	}
	req := &rpc.UserRegisterReq{}
	log.NewInfo("UserRegister args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAuthName)
	client := rpc.NewAuthClient(etcdConn)
	reply, err := client.UserRegister(context.Background(), req)

	if err != nil || reply.CommonResp.ErrCode != 0 {
		log.NewError("0", "UserRegister failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	pbDataToken := &rpc.UserTokenReq{Platform: params.Platform, FromUserID: params.UserID, OperationID: params.OperationID}
	replyToken, err := client.UserToken(context.Background(), pbDataToken)
	if err != nil {
		log.NewError(req.OperationID, "UserToken failed ", err.Error(), pbDataToken)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	resp := api.UserRegisterResp{CommResp: api.CommResp{ErrCode: replyToken.CommonResp.ErrCode, ErrMsg: replyToken.CommonResp.ErrMsg},
		UserToken: api.UserTokenInfo{UserID: req.UserInfo.UserID, Token: replyToken.Token, ExpiredTime: replyToken.ExpiredTime}}
	c.JSON(http.StatusOK, resp)
	log.NewInfo(req.OperationID, "UserRegister return ", resp)
}

func UserToken(c *gin.Context) {
	params := api.UserTokenReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	if params.Secret != config.Config.Secret {
		log.NewError(params.OperationID, "params.Secret != config.Config.Secret", params.Secret, config.Config.Secret)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 401, "errMsg": "not authorized"})
		return
	}
	req := &rpc.UserTokenReq{Platform: params.Platform, FromUserID: params.UserID, OperationID: params.OperationID}
	log.NewInfo("UserToken args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAuthName)
	client := rpc.NewAuthClient(etcdConn)
	reply, err := client.UserToken(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "UserToken failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	resp := api.UserTokenResp{CommResp: api.CommResp{ErrCode: reply.CommonResp.ErrCode, ErrMsg: reply.CommonResp.ErrMsg},
		UserToken: api.UserTokenInfo{UserID: req.FromUserID, Token: reply.Token, ExpiredTime: reply.ExpiredTime}}
	c.JSON(http.StatusOK, resp)
	log.NewInfo(req.OperationID, "UserRegister return ", resp)
}
