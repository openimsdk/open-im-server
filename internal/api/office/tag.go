package office

import (
	apistruct "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbOffice "Open_IM/pkg/proto/office"
	"Open_IM/pkg/utils"
	"context"
	"github.com/gin-gonic/gin"
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
	if err := utils.CopyStructFields(&reqPb, req); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfficeName)
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
	resp.Data.Tags = respPb.Tags
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
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfficeName)
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
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfficeName)
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
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfficeName)
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
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfficeName)
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
	if err := utils.CopyStructFields(&reqPb, req); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImOfficeName)
	client := pbOffice.NewOfficeServiceClient(etcdConn)
	respPb, err := client.GetTagSendLogs(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUserTags failed", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "CreateTag rpc server failed" + err.Error()})
		return
	}
	if err := utils.CopyStructFields(&resp.CommResp, respPb.CommonResp); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	resp.Data.Logs = respPb.TagSendLogs
	resp.Data.ShowNumber = respPb.Pagination.ShowNumber
	resp.Data.CurrentPage = respPb.Pagination.CurrentPage
	c.JSON(http.StatusOK, resp)
}
