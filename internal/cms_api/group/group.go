package group

import (
	"Open_IM/pkg/cms_api_struct"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	openIMHttp "Open_IM/pkg/common/http"
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

func GetGroupByID(c *gin.Context) {
	var (
		req   cms_api_struct.GetGroupByIDRequest
		resp  cms_api_struct.GetGroupByIDResponse
		reqPb pbGroup.GetGroupByIDReq
	)
	if err := c.ShouldBindQuery(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.OperationID = utils.OperationIDGenerator()
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", req)
	reqPb.GroupID = req.GroupID
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbGroup.NewGroupClient(etcdConn)
	respPb, err := client.GetGroupByID(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "GetGroupById failed ", err.Error())
		openIMHttp.RespHttp200(c, err, nil)
		return
	}
	utils.CopyStructFields(&resp, respPb.CMSGroup.GroupInfo)
	resp.GroupOwnerID = respPb.CMSGroup.GroupOwnerUserID
	resp.GroupOwnerName = respPb.CMSGroup.GroupOwnerUserName

	log.NewInfo("", utils.GetSelfFuncName(), "req: ", resp)
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func GetGroups(c *gin.Context) {
	var (
		req   cms_api_struct.GetGroupsRequest
		resp  cms_api_struct.GetGroupsResponse
		reqPb pbGroup.GetGroupsReq
	)
	if err := c.ShouldBindQuery(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
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
	client := pbGroup.NewGroupClient(etcdConn)
	respPb, err := client.GetGroups(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "GetUserInfo failed ", err.Error())
		openIMHttp.RespHttp200(c, err, nil)
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
	resp.CurrentPage = int(respPb.Pagination.PageNumber)
	resp.ShowNumber = int(respPb.Pagination.ShowNumber)
	log.NewInfo("", utils.GetSelfFuncName(), "resp: ", resp)
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func GetGroupByName(c *gin.Context) {
	var (
		req   cms_api_struct.GetGroupRequest
		resp  cms_api_struct.GetGroupResponse
		reqPb pbGroup.GetGroupReq
	)
	if err := c.ShouldBindQuery(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.OperationID = utils.OperationIDGenerator()
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", req)
	reqPb.GroupName = req.GroupName
	reqPb.Pagination = &commonPb.RequestPagination{}
	utils.CopyStructFields(&reqPb.Pagination, req)
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbGroup.NewGroupClient(etcdConn)
	respPb, err := client.GetGroup(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "GetGroup failed", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrServer, nil)
		return
	}
	for _, v := range respPb.CMSGroups {
		groupResp := cms_api_struct.GroupResponse{}
		utils.CopyStructFields(&groupResp, v.GroupInfo)
		groupResp.GroupOwnerName = v.GroupOwnerUserName
		groupResp.GroupOwnerID = v.GroupOwnerUserID
		resp.Groups = append(resp.Groups, groupResp)
	}
	resp.CurrentPage = int(respPb.Pagination.PageNumber)
	resp.ShowNumber = int(respPb.Pagination.ShowNumber)
	resp.GroupNums = int(respPb.GroupNums)
	log.NewInfo("", utils.GetSelfFuncName(), "resp: ", resp)
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func CreateGroup(c *gin.Context) {
	var (
		req   cms_api_struct.CreateGroupRequest
		_     cms_api_struct.CreateGroupResponse
		reqPb pbGroup.CreateGroupReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "BindJSON failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.OperationID = utils.OperationIDGenerator()
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", req)
	reqPb.GroupInfo = &commonPb.GroupInfo{}
	reqPb.GroupInfo.GroupName = req.GroupName
	reqPb.GroupInfo.CreatorUserID = req.GroupMasterId
	reqPb.OwnerUserID = req.GroupMasterId
	reqPb.OpUserID = req.GroupMasterId
	for _, v := range req.GroupMembers {
		reqPb.InitMemberList = append(reqPb.InitMemberList, &pbGroup.GroupAddMemberInfo{
			UserID:    v,
			RoleLevel: 1,
		})
	}
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbGroup.NewGroupClient(etcdConn)
	_, err := client.CreateGroup(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "CreateGroup failed", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrServer, nil)
		return
	}
	openIMHttp.RespHttp200(c, constant.OK, nil)
}

func GetGroupMembers(c *gin.Context) {
	var (
		req   cms_api_struct.GetGroupMembersRequest
		reqPb pbGroup.GetGroupMembersCMSReq
		resp  cms_api_struct.GetGroupMembersResponse
	)
	if err := c.ShouldBindQuery(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
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
		openIMHttp.RespHttp200(c, err, nil)
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
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func AddGroupMembers(c *gin.Context) {
	var (
		req   cms_api_struct.RemoveGroupMembersRequest
		resp  cms_api_struct.RemoveGroupMembersResponse
		reqPb pbGroup.AddGroupMembersCMSReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "BindJSON failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.OperationID = utils.OperationIDGenerator()
	log.NewInfo("", utils.GetSelfFuncName(), "req: ", req)
	reqPb.UserIDList = req.Members
	reqPb.GroupID = req.GroupId
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbGroup.NewGroupClient(etcdConn)
	respPb, err := client.AddGroupMembersCMS(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "AddGroupMembersCMS failed", err.Error())
		openIMHttp.RespHttp200(c, err, nil)
		return
	}
	resp.Success = respPb.Success
	resp.Failed = respPb.Failed
	log.NewInfo("", utils.GetSelfFuncName(), "resp: ", resp)
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func RemoveGroupMembers(c *gin.Context) {
	var (
		req   cms_api_struct.RemoveGroupMembersRequest
		resp  cms_api_struct.RemoveGroupMembersResponse
		reqPb pbGroup.RemoveGroupMembersCMSReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "BindJSON failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.OperationID = utils.OperationIDGenerator()
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", req)
	reqPb.UserIDList = req.Members
	reqPb.GroupID = req.GroupId
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbGroup.NewGroupClient(etcdConn)
	respPb, err := client.RemoveGroupMembersCMS(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "RemoveGroupMembersCMS failed", err.Error())
		openIMHttp.RespHttp200(c, err, nil)
		return
	}
	resp.Success = respPb.Success
	resp.Failed = respPb.Failed
	log.NewInfo("", utils.GetSelfFuncName(), "req: ", resp)
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func SetGroupOwner(c *gin.Context) {
	var (
		req   cms_api_struct.SetGroupMasterRequest
		_     cms_api_struct.SetGroupMasterResponse
		reqPb pbGroup.OperateUserRoleReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "BindJSON failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.OperationID = utils.OperationIDGenerator()
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", req)
	reqPb.GroupID = req.GroupId
	reqPb.UserID = req.UserId
	reqPb.RoleLevel = constant.GroupOwner
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbGroup.NewGroupClient(etcdConn)
	_, err := client.OperateUserRole(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "DeleteGroup failed", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrServer, nil)
		return
	}
	openIMHttp.RespHttp200(c, constant.OK, nil)
}

func SetGroupOrdinaryUsers(c *gin.Context) {
	var (
		req   cms_api_struct.SetGroupMemberRequest
		_     cms_api_struct.AdminLoginResponse
		reqPb pbGroup.OperateUserRoleReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "BindJSON failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.OperationID = utils.OperationIDGenerator()
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", req)
	reqPb.GroupID = req.GroupId
	reqPb.UserID = req.UserId
	reqPb.RoleLevel = constant.GroupOrdinaryUsers
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbGroup.NewGroupClient(etcdConn)
	_, err := client.OperateUserRole(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "DeleteGroup failed", err.Error())
		openIMHttp.RespHttp200(c, err, nil)
		return
	}
	openIMHttp.RespHttp200(c, constant.OK, nil)
}

func AlterGroupInfo(c *gin.Context) {
	var (
		req   cms_api_struct.AlterGroupInfoRequest
		_     cms_api_struct.SetGroupMasterResponse
		reqPb pbGroup.SetGroupInfoReq
	)
	reqPb.OperationID = utils.OperationIDGenerator()
	log.NewInfo(reqPb.OperationID, utils.GetSelfFuncName(), "req: ", req)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "BindJSON failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.OpUserID = c.MustGet("userID").(string)
	reqPb.GroupInfoForSet = &commonPb.GroupInfoForSet{
		GroupID:      req.GroupID,
		GroupName:    req.GroupName,
		Introduction: req.Introduction,
		Notification: req.Notification,
		FaceURL:      req.ProfilePhoto,
	}
	etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, reqPb.OperationID)
	if etcdConn == nil {
		errMsg := reqPb.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(reqPb.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := pbGroup.NewGroupClient(etcdConn)
	_, err := client.SetGroupInfo(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "DeleteGroup failed", err.Error())
		openIMHttp.RespHttp200(c, err, nil)
		return
	}
	openIMHttp.RespHttp200(c, constant.OK, nil)
}
