package admin

import (
	"Open_IM/pkg/cms_api_struct"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	openIMHttp "Open_IM/pkg/common/http"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbAdmin "Open_IM/pkg/proto/admin_cms"
	"Open_IM/pkg/utils"
	"context"
	"strings"

	"github.com/gin-gonic/gin"
)

// register
func AdminLogin(c *gin.Context) {
	var (
		req cms_api_struct.AdminLoginRequest
		resp cms_api_struct.AdminLoginResponse
		reqPb pbAdmin.AdminLoginReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewInfo("0", utils.GetSelfFuncName(), err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.Secret = req.Secret
	reqPb.AdminID = req.AdminName
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAdminCMSName)
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
