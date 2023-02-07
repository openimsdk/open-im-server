package user

import (
	"Open_IM/pkg/cms_api_struct"
	"Open_IM/pkg/cms_struct"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/getcdv3"
	pbAdminCms "Open_IM/pkg/proto/admin_cms"
	commonPb "Open_IM/pkg/proto/sdk_ws"
	pb "Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetUserIDByEmailAndPhoneNumber(c *gin.Context) {
	var (
		req    cms_struct.GetUserIDByEmailAndPhoneNumberRequest
		resp   cms_struct.GetUserIDByEmailAndPhoneNumberResponse
		reqPb  pbAdminCms.GetUserIDByEmailAndPhoneNumberReq
		respPb *pbAdminCms.GetUserIDByEmailAndPhoneNumberResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req)
	reqPb.OperationID = req.OperationID
	reqPb.Email = req.Email
	reqPb.PhoneNumber = req.PhoneNumber
	etcdConn := rpc.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAdminCMSName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbAdminCms.NewAdminCMSClient(etcdConn)
	respPb, err := client.GetUserIDByEmailAndPhoneNumber(context.Background(), &reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "rpc failed", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	resp.UserIDList = respPb.UserIDList
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", resp)
	c.JSON(http.StatusOK, gin.H{"errCode": respPb.CommonResp.ErrCode, "errMsg": respPb.CommonResp.ErrMsg, "data": resp})
}
