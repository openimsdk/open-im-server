package friend

import (
	"Open_IM/pkg/cms_api_struct"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbAdmin "Open_IM/pkg/proto/admin_cms"
	pbCommon "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetUserFriends(c *gin.Context) {
	var (
		req   cms_api_struct.GetFriendsReq
		resp  cms_api_struct.GetFriendsResp
		reqPb pbAdmin.GetUserFriendsReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", req)
	reqPb.Pagination = &pbCommon.RequestPagination{}
	utils.CopyStructFields(&reqPb.Pagination, req)
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAdminCMSName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	reqPb.OperationID = req.OperationID
	reqPb.UserID = req.UserID
	reqPb.FriendUserName = req.FriendUserName
	reqPb.FriendUserID = req.FriendUserID

	client := pbAdmin.NewAdminCMSClient(etcdConn)
	respPb, err := client.GetUserFriends(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "GetUserInfo failed ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	for _, v := range respPb.FriendInfoList {
		friend := &cms_api_struct.FriendInfo{}
		utils.CopyStructFields(friend, v)
		friend.Nickname = v.FriendUser.Nickname
		friend.UserID = v.FriendUser.UserID
		resp.FriendInfoList = append(resp.FriendInfoList, friend)
	}
	resp.FriendNums = respPb.FriendNums
	resp.CurrentPage = int(respPb.Pagination.CurrentPage)
	resp.ShowNumber = int(respPb.Pagination.ShowNumber)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp)
	c.JSON(http.StatusOK, gin.H{"errCode": respPb.CommonResp.ErrCode, "errMsg": respPb.CommonResp.ErrMsg, "data": resp})
}
