package conversation

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	"Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type paramsSetReceiveMessageOpt struct {
	OperationID        string   `json:"operationID" binding:"required"`
	Option             int32    `json:"option" binding:"required"`
	ConversationIdList []string `json:"conversationIdList" binding:"required"`
}

type OptResult struct {
	ConversationId string `json:"conversationId" binding:"required"`
	Result         int32  `json:"result" binding:"required"`
}

type SetReceiveMessageOptResp struct {
	ErrCode int32       `json:"errCode"`
	ErrMsg  string      `json:"errMsg"`
	Data    []OptResult `json:"data"`
}

type paramGetReceiveMessageOpt struct {
	ConversationIdList []string `json:"ConversationIdList" binding:"required"`
	OperationID        string   `json:"operationID" binding:"required"`
}

type GetReceiveMessageOptResp struct {
	SetReceiveMessageOptResp
}

type paramGetAllConversationMessageOpt struct {
	OperationID string `json:"operationID" binding:"required"`
}

type GetAllConversationMessageOptResp struct {
	SetReceiveMessageOptResp
}

//CopyStructFields

func GetAllConversationMessageOpt(c *gin.Context) {
	params := paramGetAllConversationMessageOpt{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError(params.OperationID, "bind json failed ", err.Error(), c)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}

	claims, err := token_verify.ParseToken(c.Request.Header.Get("token"))
	if err != nil {
		log.NewError(params.OperationID, "ParseToken failed, ", err.Error(), c.Request.Header.Get("token"))
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 401, "errMsg": "ParseToken failed, " + err.Error()})
		return
	}

	req := &user.GetAllConversationMsgOptReq{
		UId:         claims.UID,
		OperationID: params.OperationID,
	}
	log.NewInfo(req.OperationID, "GetAllConversationMsgOpt req: ", req)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := user.NewUserClient(etcdConn)
	resp, err := client.GetAllConversationMsgOpt(context.Background(), req)
	if err != nil {
		log.NewError(params.OperationID, "GetAllConversationMsgOpt rpc failed, ", req, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 401, "errMsg": "GetAllConversationMsgOpt rpc failed, " + err.Error()})
		return
	}
	var ginResp GetAllConversationMessageOptResp
	ginResp.ErrCode = resp.ErrCode
	ginResp.ErrMsg = resp.ErrMsg
	for _, v := range resp.ConversationOptResult {
		var opt OptResult
		err := utils.CopyStructFields(&opt, v)
		if err != nil {
			log.NewError(req.OperationID, "CopyStructFields failed ", err.Error())
			continue
		}
		ginResp.Data = append(ginResp.Data, opt)
	}
	log.NewInfo(req.OperationID, "GetAllConversationMsgOpt resp: ", ginResp, req)
	c.JSON(http.StatusOK, ginResp)
}

func GetReceiveMessageOpt(c *gin.Context) {
	params := paramGetReceiveMessageOpt{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError(params.OperationID, "bind json failed ", err.Error(), c)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}

	claims, err := token_verify.ParseToken(c.Request.Header.Get("token"))
	if err != nil {
		log.NewError(params.OperationID, "ParseToken failed, ", err.Error(), c.Request.Header.Get("token"))
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 401, "errMsg": "ParseToken failed, " + err.Error()})
		return
	}

	req := &user.GetReceiveMessageOptReq{
		UId:            claims.UID,
		ConversationId: params.ConversationIdList,
		OperationID:    params.OperationID,
	}
	log.NewInfo(req.OperationID, "GetReceiveMessageOptReq req: ", req)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := user.NewUserClient(etcdConn)
	resp, err := client.GetReceiveMessageOpt(context.Background(), req)
	if err != nil {
		log.NewError(params.OperationID, "GetReceiveMessageOpt rpc failed, ", req, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 401, "errMsg": "GetReceiveMessageOpt rpc failed, " + err.Error()})
		return
	}
	log.NewInfo(req.OperationID, "GetReceiveMessageOptReq req: ", req, resp)
	var ginResp GetReceiveMessageOptResp
	ginResp.ErrCode = resp.ErrCode
	ginResp.ErrMsg = resp.ErrMsg

	for _, v := range resp.ConversationOptResult {
		var opt OptResult
		log.NewInfo("CopyStructFields begin ", v, req.OperationID)
		err := utils.CopyStructFields(&opt, v)
		log.NewInfo("CopyStructFields end ", v, req.OperationID)
		if err != nil {
			log.NewError(req.OperationID, "CopyStructFields failed ", err.Error())
			continue
		}
		ginResp.Data = append(ginResp.Data, opt)
	}
	log.NewInfo(req.OperationID, "GetReceiveMessageOpt resp: ", ginResp)
	c.JSON(http.StatusOK, ginResp)
}

func SetReceiveMessageOpt(c *gin.Context) {
	params := paramsSetReceiveMessageOpt{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError(params.OperationID, "bind json failed ", err.Error(), c)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}

	claims, err := token_verify.ParseToken(c.Request.Header.Get("token"))
	if err != nil {
		log.NewError(params.OperationID, "ParseToken failed, ", err.Error(), c.Request.Header.Get("token"))
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 401, "errMsg": "ParseToken failed, " + err.Error()})
		return
	}

	req := &user.SetReceiveMessageOptReq{
		UId:            claims.UID,
		Opt:            params.Option,
		ConversationId: params.ConversationIdList,
		OperationID:    params.OperationID,
	}
	log.NewInfo(req.OperationID, "SetReceiveMessageOpt req: ", req)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := user.NewUserClient(etcdConn)
	resp, err := client.SetReceiveMessageOpt(context.Background(), req)
	if err != nil {
		log.NewError(params.OperationID, "SetReceiveMessageOpt rpc failed, ", req, err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 401, "errMsg": "SetReceiveMessageOpt rpc failed, " + err.Error()})
		return
	}
	log.NewInfo(req.OperationID, "SetReceiveMessageOpt req: ", req, resp)
	ginResp := SetReceiveMessageOptResp{
		ErrCode: resp.ErrCode,
		ErrMsg:  resp.ErrMsg,
	}

	for _, v := range resp.OptResult {
		var opt OptResult
		log.NewDebug("CopyStructFields begin ", v, req.OperationID)
		err := utils.CopyStructFields(&opt, *v, "ConversationId", "Result")
		log.NewDebug("CopyStructFields end ", v, req.OperationID)
		if err != nil {
			log.NewError(req.OperationID, "CopyStructFields failed ", err.Error())
			continue
		}
		ginResp.Data = append(ginResp.Data, opt)
	}
	log.NewInfo(req.OperationID, "SetReceiveMessageOpt resp: ", ginResp)
	c.JSON(http.StatusOK, ginResp)
}
