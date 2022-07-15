package office

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbOffice "Open_IM/pkg/proto/office"
	pbCommon "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"net/http"
	"strings"
)

// @Summary 获取用户标签信息
// @Description 用户获取自己的所有的标签
// @Tags 标签
// @ID GetUserTags
// @Accept json
// @Param token header string true "im token"
// @Param req body api.GetUserTagsReq true "请求"
// @Produce json
// @Success 0 {object} api.GetUserTagsResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /office/get_user_tags [post]
func GetUserTags(c *gin.Context) {
	var (
		req    api.GetUserTagsReq
		resp   api.GetUserTagsResp
		reqPb  pbOffice.GetUserTagsReq
		respPb *pbOffice.GetUserTagsResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "bind json failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}

	var ok bool
	var errInfo string
	var userID string
	ok, userID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	reqPb.UserID = userID
	reqPb.OperationID = req.OperationID
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfficeName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbOffice.NewOfficeServiceClient(etcdConn)
	respPb, err := client.GetUserTags(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUserTags failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetUserTags rpc server failed" + err.Error()})
		return
	}
	if err := utils.CopyStructFields(&resp.CommResp, respPb.CommonResp); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	if respPb.Tags != nil {
		resp.Data.Tags = respPb.Tags
	} else {
		resp.Data.Tags = []*pbOffice.Tag{}
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary 创建标签
// @Description 创建标签
// @Tags 标签
// @ID CreateTag
// @Accept json
// @Param token header string true "im token"
// @Param req body api.CreateTagReq true "请求"
// @Produce json
// @Success 0 {object} api.CreateTagResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /office/create_tag [post]
func CreateTag(c *gin.Context) {
	var (
		req    api.CreateTagReq
		resp   api.CreateTagResp
		reqPb  pbOffice.CreateTagReq
		respPb *pbOffice.CreateTagResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "bind json failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}
	if err := utils.CopyStructFields(&reqPb, req); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}

	var ok bool
	var errInfo string
	var userID string
	ok, userID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	reqPb.UserID = userID
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfficeName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbOffice.NewOfficeServiceClient(etcdConn)
	respPb, err := client.CreateTag(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUserTags failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "CreateTag rpc server failed" + err.Error()})
		return
	}
	if err := utils.CopyStructFields(&resp.CommResp, respPb.CommonResp); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary 删除标签
// @Description 根据标签ID创建标签
// @Tags 标签
// @ID DeleteTag
// @Accept json
// @Param token header string true "im token"
// @Param req body api.DeleteTagReq true "请求"
// @Produce json
// @Success 0 {object} api.DeleteTagResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /office/delete_tag [post]
func DeleteTag(c *gin.Context) {
	var (
		req    api.DeleteTagReq
		resp   api.DeleteTagResp
		reqPb  pbOffice.DeleteTagReq
		respPb *pbOffice.DeleteTagResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "bind json failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}
	if err := utils.CopyStructFields(&reqPb, req); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}

	var ok bool
	var errInfo string
	var userID string
	ok, userID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	reqPb.UserID = userID
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfficeName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbOffice.NewOfficeServiceClient(etcdConn)
	respPb, err := client.DeleteTag(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUserTags failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "CreateTag rpc server failed" + err.Error()})
		return
	}
	if err := utils.CopyStructFields(&resp.CommResp, respPb.CommonResp); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary 修改标签
// @Description 根据标签ID修改标签用户列表, 名称
// @Tags 标签
// @ID SetTag
// @Accept json
// @Param token header string true "im token"
// @Param req body api.SetTagReq true "请求"
// @Produce json
// @Success 0 {object} api.SetTagResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /office/set_tag [post]
func SetTag(c *gin.Context) {
	var (
		req    api.SetTagReq
		resp   api.SetTagResp
		reqPb  pbOffice.SetTagReq
		respPb *pbOffice.SetTagResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "bind json failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}
	if err := utils.CopyStructFields(&reqPb, req); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}

	var ok bool
	var errInfo string
	var userID string
	ok, userID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	reqPb.UserID = userID
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfficeName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbOffice.NewOfficeServiceClient(etcdConn)
	respPb, err := client.SetTag(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUserTags failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "CreateTag rpc server failed" + err.Error()})
		return
	}
	if err := utils.CopyStructFields(&resp.CommResp, respPb.CommonResp); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary 发送标签消息
// @Description 对标签用户发送消息
// @Tags 标签
// @ID SendMsg2Tag
// @Accept json
// @Param token header string true "im token"
// @Param req body api.SendMsg2TagReq true "请求"
// @Produce json
// @Success 0 {object} api.SendMsg2TagResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /office/send_msg_to_tag [post]
func SendMsg2Tag(c *gin.Context) {
	var (
		req    api.SendMsg2TagReq
		resp   api.SendMsg2TagResp
		reqPb  pbOffice.SendMsg2TagReq
		respPb *pbOffice.SendMsg2TagResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "bind json failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}
	if err := utils.CopyStructFields(&reqPb, req); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}

	var ok bool
	var errInfo string
	var userID string
	ok, userID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	reqPb.SendID = userID
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfficeName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbOffice.NewOfficeServiceClient(etcdConn)
	respPb, err := client.SendMsg2Tag(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUserTags failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "CreateTag rpc server failed" + err.Error()})
		return
	}
	if err := utils.CopyStructFields(&resp.CommResp, respPb.CommonResp); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	c.JSON(http.StatusOK, resp)
}

// @Summary 获取发送历史记录
// @Description 分页获取发送历史记录
// @Tags 标签
// @ID GetTagSendLogs
// @Accept json
// @Param token header string true "im token"
// @Param req body api.GetTagSendLogsReq true "请求"
// @Produce json
// @Success 0 {object} api.GetTagSendLogsResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /office/get_send_tag_log [post]
func GetTagSendLogs(c *gin.Context) {
	var (
		req    api.GetTagSendLogsReq
		resp   api.GetTagSendLogsResp
		reqPb  pbOffice.GetTagSendLogsReq
		respPb *pbOffice.GetTagSendLogsResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "bind json failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}

	var ok bool
	var errInfo string
	var userID string
	ok, userID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	reqPb.UserID = userID
	reqPb.OperationID = req.OperationID
	reqPb.Pagination = &pbCommon.RequestPagination{
		PageNumber: req.PageNumber,
		ShowNumber: req.ShowNumber,
	}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfficeName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbOffice.NewOfficeServiceClient(etcdConn)
	maxSizeOption := grpc.MaxCallRecvMsgSize(1024 * 1024 * 20)
	respPb, err := client.GetTagSendLogs(context.Background(), &reqPb, maxSizeOption)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetTagSendLogs failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "GetTagSendLogs rpc server failed" + err.Error()})
		return
	}
	if err := utils.CopyStructFields(&resp.CommResp, respPb.CommonResp); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	if respPb.TagSendLogs != nil {
		resp.Data.Logs = respPb.TagSendLogs
	} else {
		resp.Data.Logs = []*pbOffice.TagSendLog{}
	}
	resp.Data.ShowNumber = respPb.Pagination.ShowNumber
	resp.Data.CurrentPage = respPb.Pagination.CurrentPage
	c.JSON(http.StatusOK, resp)
}

// @Summary 获取该用户的标签信息
// @Description 通过标签id获取该用户的标签信息
// @Tags 标签
// @ID GetUserTagByID
// @Accept json
// @Param token header string true "im token"
// @Param req body api.GetUserTagByIDReq true "请求"
// @Produce json
// @Success 0 {object} api.GetUserTagByIDResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /office/get_user_tag_by_id [post]
func GetUserTagByID(c *gin.Context) {
	var (
		req    api.GetUserTagByIDReq
		resp   api.GetUserTagByIDResp
		reqPb  pbOffice.GetUserTagByIDReq
		respPb *pbOffice.GetUserTagByIDResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "bind json failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": "bind json failed " + err.Error()})
		return
	}

	var ok bool
	var errInfo string
	var userID string
	ok, userID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	reqPb.UserID = userID
	reqPb.OperationID = req.OperationID
	reqPb.TagID = req.TagID
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfficeName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbOffice.NewOfficeServiceClient(etcdConn)
	respPb, err := client.GetUserTagByID(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUserTagByID failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "CreateTag rpc server failed" + err.Error()})
		return
	}
	if err := utils.CopyStructFields(&resp.CommResp, respPb.CommonResp); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	resp.Data.Tag = respPb.Tag
	c.JSON(http.StatusOK, resp)
}
