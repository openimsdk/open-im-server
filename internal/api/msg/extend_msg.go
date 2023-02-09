package msg

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	rpc "Open_IM/pkg/proto/msg"
	"Open_IM/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func SetMessageReactionExtensions(c *gin.Context) {
	var (
		req   api.SetMessageReactionExtensionsReq
		resp  api.SetMessageReactionExtensionsResp
		reqPb rpc.SetMessageReactionExtensionsReq
	)

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req:", req)

	if err := utils.CopyStructFields(&reqPb, &req); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields", err.Error())
	}
	var ok bool
	var errInfo string
	ok, reqPb.OpUserID, errInfo, reqPb.OpUserIDPlatformID = token_verify.GetUserIDAndPlatformIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	grpcConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMsgName, req.OperationID)
	if grpcConn == nil {
		errMsg := req.OperationID + " getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	msgClient := rpc.NewMsgClient(grpcConn)
	respPb, err := msgClient.SetMessageReactionExtensions(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "DelMsgList failed", err.Error(), reqPb)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": constant.ErrServer.ErrCode, "errMsg": constant.ErrServer.ErrMsg + err.Error()})
		return
	}
	resp.ErrCode = respPb.ErrCode
	resp.ErrMsg = respPb.ErrMsg
	resp.Data.ResultKeyValue = respPb.Result
	resp.Data.MsgFirstModifyTime = reqPb.MsgFirstModifyTime
	resp.Data.IsReact = reqPb.IsReact
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), resp)
	c.JSON(http.StatusOK, resp)

}

func GetMessageListReactionExtensions(c *gin.Context) {
	var (
		req   api.GetMessageListReactionExtensionsReq
		resp  api.GetMessageListReactionExtensionsResp
		reqPb rpc.GetMessageListReactionExtensionsReq
	)
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusOK, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req:", req)

	if err := utils.CopyStructFields(&reqPb, &req); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields", err.Error())
	}
	var ok bool
	var errInfo string
	ok, reqPb.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusOK, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	grpcConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMsgName, req.OperationID)
	if grpcConn == nil {
		errMsg := req.OperationID + " getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	msgClient := rpc.NewMsgClient(grpcConn)
	respPb, err := msgClient.GetMessageListReactionExtensions(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "DelMsgList failed", err.Error(), reqPb)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": constant.ErrServer.ErrCode, "errMsg": constant.ErrServer.ErrMsg + err.Error()})
		return
	}
	resp.ErrCode = respPb.ErrCode
	resp.ErrMsg = respPb.ErrMsg
	resp.Data = respPb.SingleMessageResult
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), resp)
	c.JSON(http.StatusOK, resp)
}

func AddMessageReactionExtensions(c *gin.Context) {
	var (
		req   api.AddMessageReactionExtensionsReq
		resp  api.AddMessageReactionExtensionsResp
		reqPb rpc.AddMessageReactionExtensionsReq
	)

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req:", req)

	if err := utils.CopyStructFields(&reqPb, &req); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields", err.Error())
	}
	var ok bool
	var errInfo string
	ok, reqPb.OpUserID, errInfo, reqPb.OpUserIDPlatformID = token_verify.GetUserIDAndPlatformIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	grpcConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMsgName, req.OperationID)
	if grpcConn == nil {
		errMsg := req.OperationID + " getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	msgClient := rpc.NewMsgClient(grpcConn)
	respPb, err := msgClient.AddMessageReactionExtensions(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "DelMsgList failed", err.Error(), reqPb)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": constant.ErrServer.ErrCode, "errMsg": constant.ErrServer.ErrMsg + err.Error()})
		return
	}
	resp.ErrCode = respPb.ErrCode
	resp.ErrMsg = respPb.ErrMsg
	resp.Data.ResultKeyValue = respPb.Result
	resp.Data.MsgFirstModifyTime = respPb.MsgFirstModifyTime
	resp.Data.IsReact = respPb.IsReact
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), resp)
	c.JSON(http.StatusOK, resp)
}

func DeleteMessageReactionExtensions(c *gin.Context) {
	var (
		req   api.DeleteMessageReactionExtensionsReq
		resp  api.DeleteMessageReactionExtensionsResp
		reqPb rpc.DeleteMessageListReactionExtensionsReq
	)
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req:", req)

	if err := utils.CopyStructFields(&reqPb, &req); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields", err.Error())
	}
	var ok bool
	var errInfo string
	ok, reqPb.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	grpcConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMsgName, req.OperationID)
	if grpcConn == nil {
		errMsg := req.OperationID + " getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	msgClient := rpc.NewMsgClient(grpcConn)
	respPb, err := msgClient.DeleteMessageReactionExtensions(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "DelMsgList failed", err.Error(), reqPb)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": constant.ErrServer.ErrCode, "errMsg": constant.ErrServer.ErrMsg + err.Error()})
		return
	}
	resp.ErrCode = respPb.ErrCode
	resp.ErrMsg = respPb.ErrMsg
	resp.Data = respPb.Result
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), resp)
	c.JSON(http.StatusOK, resp)
}
