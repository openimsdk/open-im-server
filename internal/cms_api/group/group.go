package group

import (
	"Open_IM/pkg/cms_api_struct"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	commonPb "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"net/http"
	"strings"

	pbGroup "Open_IM/pkg/proto/group"

	"github.com/gin-gonic/gin"
)

func GetGroups(c *gin.Context) {
	var (
		req   cms_api_struct.GetGroupsRequest
		resp  cms_api_struct.GetGroupsResponse
		reqPb pbGroup.GetGroupsReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	reqPb.OperationID = utils.OperationIDGenerator()
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", req)
	reqPb.Pagination = &commonPb.RequestPagination{}
	utils.CopyStructFields(&reqPb.Pagination, req)
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	reqPb.GroupID = req.GroupID
	reqPb.GroupName = req.GroupName
	client := pbGroup.NewGroupClient(etcdConn)
	respPb, err := client.GetGroups(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "GetUserInfo failed ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	for _, v := range respPb.CMSGroups {
		groupResp := cms_api_struct.GroupResponse{}
		utils.CopyStructFields(&groupResp, v.GroupInfo)
		groupResp.GroupOwnerName = v.GroupOwnerUserName
		groupResp.GroupOwnerID = v.GroupOwnerUserID
		resp.Groups = append(resp.Groups, groupResp)
	}
	resp.GroupNums = int(respPb.GroupNum)
	resp.CurrentPage = int(respPb.Pagination.CurrentPage)
	resp.ShowNumber = int(respPb.Pagination.ShowNumber)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp)
	c.JSON(http.StatusOK, gin.H{"errCode": respPb.CommonResp.ErrCode, "errMsg": respPb.CommonResp.ErrMsg, "data": resp})
}

func GetGroupMembers(c *gin.Context) {
	var (
		req   cms_api_struct.GetGroupMembersRequest
		reqPb pbGroup.GetGroupMembersCMSReq
		resp  cms_api_struct.GetGroupMembersResponse
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	reqPb.OperationID = utils.OperationIDGenerator()
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", req)
	reqPb.Pagination = &commonPb.RequestPagination{
		PageNumber: int32(req.PageNumber),
		ShowNumber: int32(req.ShowNumber),
	}
	reqPb.GroupID = req.GroupID
	reqPb.UserName = req.UserName
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbGroup.NewGroupClient(etcdConn)
	respPb, err := client.GetGroupMembersCMS(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "GetGroupMembersCMS failed:", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}
	resp.ResponsePagination = cms_api_struct.ResponsePagination{
		CurrentPage: int(respPb.Pagination.CurrentPage),
		ShowNumber:  int(respPb.Pagination.ShowNumber),
	}
	resp.MemberNums = int(respPb.MemberNums)
	for _, groupMember := range respPb.Members {
		memberResp := cms_api_struct.GroupMemberResponse{}
		utils.CopyStructFields(&memberResp, groupMember)
		resp.GroupMembers = append(resp.GroupMembers, memberResp)
	}
	log.NewInfo("", utils.GetSelfFuncName(), "req: ", resp)
	c.JSON(http.StatusOK, gin.H{"errCode": respPb.CommonResp.ErrCode, "errMsg": respPb.CommonResp.ErrMsg, "data": resp})
}
