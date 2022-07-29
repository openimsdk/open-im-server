package msg

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	rpc "Open_IM/pkg/proto/msg"
	pbCommon "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"net/http"
	"strings"
)

// @Summary 根据seq列表删除消息
// @Description 根据seq列表删除消息
// @Tags 消息相关
// @ID DelMsg
// @Accept json
// @Param token header string true "im token"
// @Param req body api.DelMsgReq true "userID为要删除的用户ID <br> seqList为seq列表"
// @Produce json
// @Success 0 {object} api.DelMsgResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /msg/del_msg [post]
func DelMsg(c *gin.Context) {
	var (
		req   api.DelMsgReq
		resp  api.DelMsgResp
		reqPb pbCommon.DelMsgListReq
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

	grpcConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMsgName, req.OperationID)
	if grpcConn == nil {
		errMsg := req.OperationID + " getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	msgClient := rpc.NewMsgClient(grpcConn)
	respPb, err := msgClient.DelMsgList(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "DelMsgList failed", err.Error(), reqPb)
		c.JSON(http.StatusOK, gin.H{"errCode": constant.ErrServer.ErrCode, "errMsg": constant.ErrServer.ErrMsg + err.Error()})
		return
	}
	resp.ErrCode = respPb.ErrCode
	resp.ErrMsg = respPb.ErrMsg
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), resp)
	c.JSON(http.StatusOK, resp)
}
func DelSuperGroupMsg(c *gin.Context) {
	var (
		req  api.DelSuperGroupMsgReq
		resp api.DelSuperGroupMsgResp
	)
	rpcReq := &rpc.DelSuperGroupMsgReq{}
	utils.CopyStructFields(req, &req)
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req:", req)
	var ok bool
	var errInfo string
	ok, rpcReq.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + rpcReq.OpUserID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	options := make(map[string]bool, 5)
	utils.SetSwitchFromOptions(options, constant.IsConversationUpdate, false)
	utils.SetSwitchFromOptions(options, constant.IsSenderConversationUpdate, false)
	utils.SetSwitchFromOptions(options, constant.IsUnreadCount, false)
	utils.SetSwitchFromOptions(options, constant.IsOfflinePush, false)
	pbData := rpc.SendMsgReq{
		OperationID: req.OperationID,
		MsgData: &pbCommon.MsgData{
			SendID:      req.UserID,
			RecvID:      req.UserID,
			ClientMsgID: utils.GetMsgID(req.UserID),
			SessionType: constant.SingleChatType,
			MsgFrom:     constant.SysMsgType,
			ContentType: constant.MsgDeleteNotification,
			//	ForceList:        params.ForceList,
			CreateTime: utils.GetCurrentTimestampByMill(),
			Options:    options,
		},
	}
	var tips pbCommon.TipsComm
	deleteMsg := api.MsgDeleteNotificationElem{
		GroupID:     req.GroupID,
		IsAllDelete: req.IsAllDelete,
		SeqList:     req.SeqList,
	}
	tips.JsonDetail = utils.StructToJsonString(deleteMsg)
	var err error
	pbData.MsgData.Content, err = proto.Marshal(&tips)
	if err != nil {
		log.Error(req.OperationID, "Marshal failed ", err.Error(), tips.String())
		resp.ErrCode = 400
		resp.ErrMsg = err.Error()
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), resp)
		c.JSON(http.StatusOK, resp)
	}
	log.Info(req.OperationID, "", "api DelSuperGroupMsg call start..., [data: %s]", pbData.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMsgName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewMsgClient(etcdConn)

	log.Info(req.OperationID, "", "api DelSuperGroupMsg call, api call rpc...")
	if req.IsAllDelete  {
		RpcResp, err := client.DelSuperGroupMsg(context.Background(),rpcReq)
		if err != nil {
			log.NewError(req.OperationID, "call delete DelSuperGroupMsg rpc server failed", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call DelSuperGroupMsg  rpc server failed"})
			return
		}
		log.Info(req.OperationID, "", "api DelSuperGroupMsg call end..., [data: %s] [reply: %s]", pbData.String(), RpcResp.String())
		resp.ErrCode = RpcResp.ErrCode
		resp.ErrMsg = RpcResp.ErrMsg
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), resp)
		c.JSON(http.StatusOK, resp)
	}else{
		RpcResp, err := client.SendMsg(context.Background(), &pbData)
		if err != nil {
			log.NewError(req.OperationID, "call delete UserSendMsg rpc server failed", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call UserSendMsg  rpc server failed"})
			return
		}
		log.Info(req.OperationID, "", "api DelSuperGroupMsg call end..., [data: %s] [reply: %s]", pbData.String(), RpcResp.String())
		resp.ErrCode = RpcResp.ErrCode
		resp.ErrMsg = RpcResp.ErrMsg
		log.NewInfo(req.OperationID, utils.GetSelfFuncName(), resp)
		c.JSON(http.StatusOK, resp)
	}

}

// @Summary 清空用户消息
// @Description 清空用户消息
// @Tags 消息相关
// @ID ClearMsg
// @Accept json
// @Param token header string true "im token"
// @Param req body api.CleanUpMsgReq true "userID为要清空的用户ID"
// @Produce json
// @Success 0 {object} api.CleanUpMsgResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /msg/clear_msg [post]
func ClearMsg(c *gin.Context) {
	params := api.CleanUpMsgReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	//
	req := &rpc.ClearMsgReq{}
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

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " api args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMsgName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + " getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewMsgClient(etcdConn)
	RpcResp, err := client.ClearMsg(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, " CleanUpMsg failed ", err.Error(), req.String(), RpcResp.ErrMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": RpcResp.ErrMsg})
		return
	}

	resp := api.CleanUpMsgResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary 设置用户最小seq
// @Description 设置用户最小seq，以及用户相关读扩散群组最小seq
// @Tags 消息相关
// @ID SetMsgMinSeq
// @Accept json
// @Param token header string true "im token"
func SetMsgMinSeq(c *gin.Context) {
	params := api.SetMsgMinSeqReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	//
	req := &rpc.SetMsgMinSeqReq{}
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
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " api args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMsgName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + " getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewMsgClient(etcdConn)
	RpcResp, err := client.SetMsgMinSeq(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, " SetMsgMinSeq failed ", err.Error(), req.String(), RpcResp.ErrMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": RpcResp.ErrMsg})
		return
	}

	resp := api.SetMsgMinSeqResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " api return ", resp)
	c.JSON(http.StatusOK, resp)
}
