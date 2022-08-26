package messageCMS

import (
	"Open_IM/pkg/cms_api_struct"
	"Open_IM/pkg/common/config"
	openIMHttp "Open_IM/pkg/common/http"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbAdminCMS "Open_IM/pkg/proto/admin_cms"
	pbCommon "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"net/http"
	"strings"

	"Open_IM/pkg/common/constant"

	"github.com/gin-gonic/gin"
)

func GetChatLogs(c *gin.Context) {
	var (
		req   cms_api_struct.GetChatLogsReq
		resp  cms_api_struct.GetChatLogsResp
		reqPb pbAdminCMS.GetChatLogsReq
	)
	if err := c.ShouldBindQuery(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, resp)
		return
	}
	reqPb.Pagination = &pbCommon.RequestPagination{
		PageNumber: int32(req.PageNumber),
		ShowNumber: int32(req.ShowNumber),
	}
	utils.CopyStructFields(&reqPb, &req)
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", req)
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMessageCMSName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbAdminCMS.NewAdminCMSClient(etcdConn)
	respPb, err := client.GetChatLogs(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "GetChatLogs rpc failed", err.Error())
		openIMHttp.RespHttp200(c, err, resp)
		return
	}
	for _, v := range respPb.ChatLogs {
		chatLog := pbCommon.MsgData{}
		utils.CopyStructFields(&chatLog, v)
		resp.ChatLogs = append(resp.ChatLogs, &chatLog)
	}
	resp.ShowNumber = int(respPb.Pagination.ShowNumber)
	resp.CurrentPage = int(respPb.Pagination.CurrentPage)
	resp.ChatLogsNum = int(respPb.ChatLogsNum)
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "resp", resp)
	openIMHttp.RespHttp200(c, constant.OK, resp)
}
