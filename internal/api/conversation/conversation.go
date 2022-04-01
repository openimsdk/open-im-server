package conversation

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbUser "Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func SetConversation(c *gin.Context) {
	var (
		req   api.SetConversationReq
		resp  api.SetConversationResp
		reqPb pbUser.SetConversationReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "bind json failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}
	reqPb.Conversation = &pbUser.Conversation{}
	err := utils.CopyStructFields(&reqPb, req)
	err = utils.CopyStructFields(reqPb.Conversation, req.Conversation)
	if err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", reqPb.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := pbUser.NewUserClient(etcdConn)
	respPb, err := client.SetConversation(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "SetConversation rpc failed, ", reqPb.String(), err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": "GetAllConversationMsgOpt rpc failed, " + err.Error()})
		return
	}
	resp.ErrMsg = respPb.CommonResp.ErrMsg
	resp.ErrCode = respPb.CommonResp.ErrCode
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp)
	c.JSON(http.StatusOK, resp)
}

func BatchSetConversations(c *gin.Context) {
	var (
		req   api.BatchSetConversationsReq
		resp  api.BatchSetConversationsResp
		reqPb pbUser.BatchSetConversationsReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "bind json failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}
	if err := utils.CopyStructFields(&reqPb, req); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", reqPb.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := pbUser.NewUserClient(etcdConn)
	respPb, err := client.BatchSetConversations(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "SetConversation rpc failed, ", reqPb.String(), err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": "GetAllConversationMsgOpt rpc failed, " + err.Error()})
		return
	}
	if err := utils.CopyStructFields(&resp.Data, respPb); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	resp.ErrMsg = respPb.CommonResp.ErrMsg
	resp.ErrCode = respPb.CommonResp.ErrCode
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp)
	c.JSON(http.StatusOK, resp)
}

func GetAllConversations(c *gin.Context) {
	var (
		req   api.GetAllConversationsReq
		resp  api.GetAllConversationsResp
		reqPb pbUser.GetAllConversationsReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "bind json failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}
	if err := utils.CopyStructFields(&reqPb, req); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", reqPb.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := pbUser.NewUserClient(etcdConn)
	respPb, err := client.GetAllConversations(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "SetConversation rpc failed, ", reqPb.String(), err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": "GetAllConversationMsgOpt rpc failed, " + err.Error()})
		return
	}
	resp.ErrMsg = respPb.CommonResp.ErrMsg
	resp.ErrCode = respPb.CommonResp.ErrCode
	if err := utils.CopyStructFields(&resp.Conversations, respPb.Conversations); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed, ", err.Error())
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp)
	c.JSON(http.StatusOK, resp)
}

func GetConversation(c *gin.Context) {
	var (
		req   api.GetConversationReq
		resp  api.GetConversationResp
		reqPb pbUser.GetConversationReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "bind json failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}
	if err := utils.CopyStructFields(&reqPb, req); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", reqPb.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := pbUser.NewUserClient(etcdConn)
	respPb, err := client.GetConversation(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetConversation rpc failed, ", reqPb.String(), err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": "GetAllConversationMsgOpt rpc failed, " + err.Error()})
		return
	}
	resp.ErrMsg = respPb.CommonResp.ErrMsg
	resp.ErrCode = respPb.CommonResp.ErrCode
	if err := utils.CopyStructFields(&resp.Conversation, respPb.Conversation); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp)
	c.JSON(http.StatusOK, resp)
}

func GetConversations(c *gin.Context) {
	var (
		req   api.GetConversationsReq
		resp  api.GetConversationsResp
		reqPb pbUser.GetConversationsReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "bind json failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}
	if err := utils.CopyStructFields(&reqPb, req); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", reqPb.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := pbUser.NewUserClient(etcdConn)
	respPb, err := client.GetConversations(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "SetConversation rpc failed, ", reqPb.String(), err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": "GetAllConversationMsgOpt rpc failed, " + err.Error()})
		return
	}
	resp.ErrMsg = respPb.CommonResp.ErrMsg
	resp.ErrCode = respPb.CommonResp.ErrCode
	if err := utils.CopyStructFields(&resp.Conversations, respPb.Conversations); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp)
	c.JSON(http.StatusOK, resp)
}

func SetRecvMsgOpt(c *gin.Context) {
	var (
		req   api.SetRecvMsgOptReq
		resp  api.SetRecvMsgOptResp
		reqPb pbUser.SetRecvMsgOptReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "bind json failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}
	if err := utils.CopyStructFields(&reqPb, req); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", reqPb.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := pbUser.NewUserClient(etcdConn)
	respPb, err := client.SetRecvMsgOpt(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "SetRecvMsgOpt rpc failed, ", reqPb.String(), err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": "GetAllConversationMsgOpt rpc failed, " + err.Error()})
		return
	}
	resp.ErrMsg = respPb.CommonResp.ErrMsg
	resp.ErrCode = respPb.CommonResp.ErrCode
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp)
	c.JSON(http.StatusOK, resp)
}

//Deprecated
func SetReceiveMessageOpt(c *gin.Context) {

}

//Deprecated
func GetReceiveMessageOpt(c *gin.Context) {

}

//Deprecated
func GetAllConversationMessageOpt(c *gin.Context) {

}
