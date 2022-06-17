package office

import (
	apistruct "Open_IM/pkg/base_info"
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

func GetUserTags(c *gin.Context) {
	var (
		req    apistruct.GetUserTagsReq
		resp   apistruct.GetUserTagsResp
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
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
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

func CreateTag(c *gin.Context) {
	var (
		req    apistruct.CreateTagReq
		resp   apistruct.CreateTagResp
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

func DeleteTag(c *gin.Context) {
	var (
		req    apistruct.DeleteTagReq
		resp   apistruct.DeleteTagResp
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

func SetTag(c *gin.Context) {
	var (
		req    apistruct.SetTagReq
		resp   apistruct.SetTagResp
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

func SendMsg2Tag(c *gin.Context) {
	var (
		req    apistruct.SendMsg2TagReq
		resp   apistruct.SendMsg2TagResp
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
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
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

func GetTagSendLogs(c *gin.Context) {
	var (
		req    apistruct.GetTagSendLogsReq
		resp   apistruct.GetTagSendLogsResp
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
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
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
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
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

func GetUserTagByID(c *gin.Context) {
	var (
		req    apistruct.GetUserTagByIDReq
		resp   apistruct.GetUserTagByIDResp
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
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
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
