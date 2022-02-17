package conversation

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	"Open_IM/pkg/proto/user"
	rpc "Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func GetAllConversationMessageOpt(c *gin.Context) {
	params := api.GetAllConversationMessageOptReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}
	req := &rpc.GetAllConversationMsgOptReq{}
	utils.CopyStructFields(req, &params)
	var ok bool
	ok, req.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed or not set token in header"})
		return
	}
	log.NewInfo(params.OperationID, "GetAllConversationMessageOpt args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := user.NewUserClient(etcdConn)
	RpcResp, err := client.GetAllConversationMsgOpt(context.Background(), req)
	if err != nil {
		log.NewError(params.OperationID, "GetAllConversationMsgOpt rpc failed, ", req, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": "GetAllConversationMsgOpt rpc failed, " + err.Error()})
		return
	}
	optResult := make([]*api.OptResult, 0)
	for _, v := range RpcResp.ConversationOptResultList {
		temp := new(api.OptResult)
		temp.ConversationID = v.ConversationID
		temp.Result = &v.Result
		optResult = append(optResult, temp)
	}
	resp := api.GetAllConversationMessageOptResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}, ConversationOptResultList: optResult}
	log.NewInfo(req.OperationID, "GetAllConversationMsgOpt api return: ", resp)
	c.JSON(http.StatusOK, resp)
}

func GetReceiveMessageOpt(c *gin.Context) {
	params := api.GetReceiveMessageOptReq{}
	if err := c.BindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}
	req := &rpc.GetReceiveMessageOptReq{}
	utils.CopyStructFields(req, &params)
	var ok bool
	ok, req.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(params.OperationID, "GetReceiveMessageOpt args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := user.NewUserClient(etcdConn)
	RpcResp, err := client.GetReceiveMessageOpt(context.Background(), req)
	if err != nil {
		log.NewError(params.OperationID, "GetReceiveMessageOpt rpc failed, ", req, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": "GetReceiveMessageOpt rpc failed, " + err.Error()})
		return
	}
	optResult := make([]*api.OptResult, 0)
	for _, v := range RpcResp.ConversationOptResultList {
		temp := new(api.OptResult)
		temp.ConversationID = v.ConversationID
		temp.Result = &v.Result
		optResult = append(optResult, temp)
	}
	resp := api.GetReceiveMessageOptResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}, ConversationOptResultList: optResult}
	log.NewInfo(req.OperationID, "GetReceiveMessageOpt api return: ", resp)
	c.JSON(http.StatusOK, resp)
}

func SetReceiveMessageOpt(c *gin.Context) {
	params := api.SetReceiveMessageOptReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError(params.OperationID, utils.GetSelfFuncName(), "bind json failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}
	req := &rpc.SetReceiveMessageOptReq{}
	utils.CopyStructFields(req, &params)
	var ok bool
	ok, req.OpUserID = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"))
	if !ok {
		log.NewError(req.OperationID, "GetUserIDFromToken false ", c.Request.Header.Get("token"))
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserIDFromToken failed"})
		return
	}
	log.NewInfo(params.OperationID, "SetReceiveMessageOpt args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := user.NewUserClient(etcdConn)
	RpcResp, err := client.SetReceiveMessageOpt(context.Background(), req)
	if err != nil {
		log.NewError(params.OperationID, "SetReceiveMessageOpt rpc failed, ", req, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": "SetReceiveMessageOpt rpc failed, " + err.Error()})
		return
	}
	optResult := make([]*api.OptResult, 0)
	for _, v := range RpcResp.ConversationOptResultList {
		temp := new(api.OptResult)
		temp.ConversationID = v.ConversationID
		temp.Result = &v.Result
		optResult = append(optResult, temp)
	}
	resp := api.SetReceiveMessageOptResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}, ConversationOptResultList: optResult}
	log.NewInfo(req.OperationID, "SetReceiveMessageOpt api return: ", resp)
	c.JSON(http.StatusOK, resp)
}
