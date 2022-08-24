package admin

import (
	apiStruct "Open_IM/pkg/cms_api_struct"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	openIMHttp "Open_IM/pkg/common/http"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbAdmin "Open_IM/pkg/proto/admin_cms"
	pbCommon "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"net/http"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	url2 "net/url"

	"github.com/gin-gonic/gin"
)

var (
	minioClient *minio.Client
)

func init() {
	operationID := utils.OperationIDGenerator()
	log.NewInfo(operationID, utils.GetSelfFuncName(), "minio config: ", config.Config.Credential.Minio)
	var initUrl string
	if config.Config.Credential.Minio.EndpointInnerEnable {
		initUrl = config.Config.Credential.Minio.EndpointInner
	} else {
		initUrl = config.Config.Credential.Minio.Endpoint
	}
	log.NewInfo(operationID, utils.GetSelfFuncName(), "use initUrl: ", initUrl)
	minioUrl, err := url2.Parse(initUrl)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "parse failed, please check config/config.yaml", err.Error())
		return
	}
	opts := &minio.Options{
		Creds: credentials.NewStaticV4(config.Config.Credential.Minio.AccessKeyID, config.Config.Credential.Minio.SecretAccessKey, ""),
	}
	if minioUrl.Scheme == "http" {
		opts.Secure = false
	} else if minioUrl.Scheme == "https" {
		opts.Secure = true
	}
	log.NewInfo(operationID, utils.GetSelfFuncName(), "Parse ok ", config.Config.Credential.Minio)
	minioClient, err = minio.New(minioUrl.Host, opts)
	log.NewInfo(operationID, utils.GetSelfFuncName(), "new ok ", config.Config.Credential.Minio)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "init minio client failed", err.Error())
		return
	}
}

// register
func AdminLogin(c *gin.Context) {
	var (
		req   apiStruct.AdminLoginRequest
		resp  apiStruct.AdminLoginResponse
		reqPb pbAdmin.AdminLoginReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError("0", utils.GetSelfFuncName(), err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.Secret = req.Secret
	reqPb.AdminID = req.AdminName
	reqPb.OperationID = utils.OperationIDGenerator()
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAdminCMSName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbAdmin.NewAdminCMSClient(etcdConn)
	respPb, err := client.AdminLogin(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		openIMHttp.RespHttp200(c, err, nil)
		return
	}
	resp.Token = respPb.Token
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func AddUserRegisterAddFriendIDList(c *gin.Context) {
	var (
		req  apiStruct.AddUserRegisterAddFriendIDListRequest
		resp apiStruct.AddUserRegisterAddFriendIDListResponse
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError("0", utils.GetSelfFuncName(), err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), req)
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAdminCMSName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbAdmin.NewAdminCMSClient(etcdConn)
	_, err := client.AddUserRegisterAddFriendIDList(context.Background(), &pbAdmin.AddUserRegisterAddFriendIDListReq{OperationID: req.OperationID, UserIDList: req.UserIDList})
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		openIMHttp.RespHttp200(c, err, nil)
		return
	}
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func ReduceUserRegisterAddFriendIDList(c *gin.Context) {
	var (
		req  apiStruct.ReduceUserRegisterAddFriendIDListRequest
		resp apiStruct.ReduceUserRegisterAddFriendIDListResponse
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError("0", utils.GetSelfFuncName(), err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), req)
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAdminCMSName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbAdmin.NewAdminCMSClient(etcdConn)
	_, err := client.ReduceUserRegisterAddFriendIDList(context.Background(), &pbAdmin.ReduceUserRegisterAddFriendIDListReq{OperationID: req.OperationID, UserIDList: req.UserIDList, Operation: req.Operation})
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		openIMHttp.RespHttp200(c, err, nil)
		return
	}
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func GetUserRegisterAddFriendIDList(c *gin.Context) {
	var (
		req  apiStruct.GetUserRegisterAddFriendIDListRequest
		resp apiStruct.GetUserRegisterAddFriendIDListResponse
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError("0", utils.GetSelfFuncName(), err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), req)
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAdminCMSName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbAdmin.NewAdminCMSClient(etcdConn)
	respPb, err := client.GetUserRegisterAddFriendIDList(context.Background(), &pbAdmin.GetUserRegisterAddFriendIDListReq{OperationID: req.OperationID, Pagination: &pbCommon.RequestPagination{
		PageNumber: int32(req.PageNumber),
		ShowNumber: int32(req.ShowNumber),
	}})
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		openIMHttp.RespHttp200(c, err, nil)
		return
	}
	resp.Users = respPb.UserInfoList
	resp.ShowNumber = int(respPb.Pagination.ShowNumber)
	resp.CurrentPage = int(respPb.Pagination.CurrentPage)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), resp)
	openIMHttp.RespHttp200(c, constant.OK, resp)
}
