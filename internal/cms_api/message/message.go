package message

import (
	"Open_IM/pkg/cms_api_struct"
	"Open_IM/pkg/common/config"
	openIMHttp "Open_IM/pkg/common/http"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbMessage "Open_IM/pkg/proto/message_cms"
	"Open_IM/pkg/utils"
	"context"
	"strings"

	"Open_IM/pkg/common/constant"

	"github.com/gin-gonic/gin"
)

func BroadcastMessage(c *gin.Context) {
	openIMHttp.RespHttp200(c, constant.OK, nil)
}

func MassSendMassage(c *gin.Context) {
	openIMHttp.RespHttp200(c, constant.OK, nil)
}

func WithdrawMessage(c *gin.Context) {
	openIMHttp.RespHttp200(c, constant.OK, nil)
}

func GetChatLogs(c *gin.Context) {
	var (
		req   cms_api_struct.GetChatLogsRequest
		resp  cms_api_struct.GetChatLogsResponse
		reqPb pbMessage.GetChatLogsReq
	)
	if err := c.ShouldBindQuery(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, resp)
		return
	}
	utils.CopyStructFields(&reqPb, &req)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImStatisticsName)
	client := pbMessage.NewMessageClient(etcdConn)
	respPb, err := client.GetChatLogs(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "GetChatLogs rpc failed", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrServer, resp)
		return
	}
	utils.CopyStructFields(&resp, &respPb)
	openIMHttp.RespHttp200(c, constant.OK, resp)
}
