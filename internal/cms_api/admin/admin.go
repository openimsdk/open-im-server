package admin

import (
	apiStruct "Open_IM/pkg/cms_api_struct"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	openIMHttp "Open_IM/pkg/common/http"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbAdmin "Open_IM/pkg/proto/admin_cms"
	"Open_IM/pkg/utils"
	"context"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	url2 "net/url"
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
		log.NewInfo("0", utils.GetSelfFuncName(), err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.Secret = req.Secret
	reqPb.AdminID = req.AdminName
	reqPb.OperationID = utils.OperationIDGenerator()
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAdminCMSName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetConn == nil"
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
