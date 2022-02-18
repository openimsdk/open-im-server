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
	"strings"

	pbGroup "Open_IM/pkg/proto/group"

	"github.com/gin-gonic/gin"
)

func GetGroupById(c *gin.Context) {
	var (
		req   cms_api_struct.GetGroupByIdRequest
		resp  cms_api_struct.GetGroupByIdResponse
		reqPb pbGroup.GetGroupByIdReq
	)
	if err := c.ShouldBindQuery(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.GroupId = req.GroupId
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pbGroup.NewGroupClient(etcdConn)
	respPb, err := client.GetGroupById(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "GetGroupById failed ", err.Error())
		openIMHttp.RespHttp200(c, err, nil)
		return
	}
	resp.GroupName = respPb.CMSGroup.GroupInfo.GroupName
	resp.GroupID = respPb.CMSGroup.GroupInfo.GroupID
	resp.CreateTime = (utils.UnixSecondToTime(int64(respPb.CMSGroup.GroupInfo.CreateTime))).String()
	resp.ProfilePhoto = respPb.CMSGroup.GroupInfo.FaceURL
	resp.GroupMasterName = respPb.CMSGroup.GroupMasterName
	resp.GroupMasterId = respPb.CMSGroup.GroupMasterId
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
	reqPb.Pagination = &commonPb.RequestPagination{}
	utils.CopyStructFields(&reqPb.Pagination, req)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pbGroup.NewGroupClient(etcdConn)
	respPb, err := client.GetGroups(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "GetUserInfo failed ", err.Error())
		openIMHttp.RespHttp200(c, err, nil)
		return
	}
	for _, v := range respPb.CMSGroups {
		resp.Groups = append(resp.Groups, cms_api_struct.GroupResponse{
			GroupName:        v.GroupInfo.GroupName,
			GroupID:          v.GroupInfo.GroupID,
			GroupMasterName:  v.GroupMasterName,
			GroupMasterId:    v.GroupMasterId,
			CreateTime:       (utils.UnixSecondToTime(int64(v.GroupInfo.CreateTime))).String(),
			IsBanChat:        false,
			IsBanPrivateChat: false,
			ProfilePhoto:     v.GroupInfo.FaceURL,
		})
	}
	resp.GroupNums = int(respPb.GroupNum)
	resp.CurrentPage = int(respPb.Pagination.PageNumber)
	resp.ShowNumber = int(respPb.Pagination.ShowNumber)
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
	reqPb.GroupName = req.GroupName
	reqPb.Pagination = &commonPb.RequestPagination{}
	utils.CopyStructFields(&reqPb.Pagination, req)
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pbGroup.NewGroupClient(etcdConn)
	respPb, err := client.GetGroup(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "GetGroup failed", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrServer, nil)
		return
	}
	for _, v := range respPb.CMSGroups {
		resp.Groups = append(resp.Groups, cms_api_struct.GroupResponse{
			GroupName:        v.GroupInfo.GroupName,
			GroupID:          v.GroupInfo.GroupID,
			GroupMasterName:  v.GroupMasterName,
			GroupMasterId:    v.GroupMasterId,
			CreateTime:       (utils.UnixSecondToTime(int64(v.GroupInfo.CreateTime))).String(),
			IsBanChat:        false,
			IsBanPrivateChat: false,
			ProfilePhoto:     v.GroupInfo.FaceURL,
		})
	}
	resp.CurrentPage = int(respPb.Pagination.PageNumber)
	resp.ShowNumber = int(respPb.Pagination.ShowNumber)
	resp.GroupNums = int(respPb.GroupNums)
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func CreateGroup(c *gin.Context) {
	var (
		req   cms_api_struct.CreateGroupRequest
		_  cms_api_struct.CreateGroupResponse
		reqPb pbGroup.CreateGroupReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "BindJSON failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
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
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pbGroup.NewGroupClient(etcdConn)
	_, err := client.CreateGroup(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "CreateGroup failed", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrServer, nil)
		return
	}
	openIMHttp.RespHttp200(c, constant.OK, nil)
}

func BanGroupChat(c *gin.Context) {
	var (
		req   cms_api_struct.BanGroupChatRequest
		reqPb pbGroup.OperateGroupStatusReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "BindJSON failed", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.GroupId = req.GroupId
	reqPb.Status = constant.GroupBanChat
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pbGroup.NewGroupClient(etcdConn)
	_, err := client.OperateGroupStatus(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "BanGroupChat failed ", err.Error())
		openIMHttp.RespHttp200(c, err, nil)
		return
	}
	openIMHttp.RespHttp200(c, constant.OK, nil)

}

func BanPrivateChat(c *gin.Context) {
	var (
		req   cms_api_struct.BanPrivateChatRequest
		reqPb pbGroup.OperateGroupStatusReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "BindJSON failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.GroupId = req.GroupId
	reqPb.Status = constant.GroupBanPrivateChat
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pbGroup.NewGroupClient(etcdConn)
	_, err := client.OperateGroupStatus(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "OperateGroupStatus failed", err.Error())
		openIMHttp.RespHttp200(c, err, nil)
		return
	}
	openIMHttp.RespHttp200(c, constant.OK, nil)
}

func OpenGroupChat(c *gin.Context) {
	var (
		req   cms_api_struct.BanPrivateChatRequest
		reqPb pbGroup.OperateGroupStatusReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "BindJSON failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.GroupId = req.GroupId
	reqPb.Status = constant.GroupOk
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pbGroup.NewGroupClient(etcdConn)
	_, err := client.OperateGroupStatus(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "OperateGroupStatus failed ", err.Error())
		openIMHttp.RespHttp200(c, err, nil)
		return
	}
	openIMHttp.RespHttp200(c, constant.OK, nil)
}

func OpenPrivateChat(c *gin.Context) {
	var (
		req   cms_api_struct.BanPrivateChatRequest
		reqPb pbGroup.OperateGroupStatusReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "OpenPrivateChat failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.GroupId = req.GroupId
	reqPb.Status = constant.GroupOk
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pbGroup.NewGroupClient(etcdConn)
	_, err := client.OperateGroupStatus(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "OperateGroupStatus failed ", err.Error())
		openIMHttp.RespHttp200(c, err, nil)
		return
	}
	openIMHttp.RespHttp200(c, constant.OK, nil)
}

func GetGroupMembers(c *gin.Context) {
	var (
		req  cms_api_struct.GetGroupMembersRequest
		reqPb pbGroup.GetGroupMembersCMSReq
		resp   cms_api_struct.GetGroupMembersResponse
	)
	if err := c.ShouldBindQuery(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "ShouldBindQuery failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.Pagination = &commonPb.RequestPagination{
		PageNumber: int32(req.PageNumber),
		ShowNumber: int32(req.ShowNumber),
	}
	reqPb.GroupId = req.GroupId
	reqPb.UserName = req.UserName
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
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
	for _, groupMembers := range respPb.Members {
		resp.GroupMembers = append(resp.GroupMembers, cms_api_struct.GroupMemberResponse{
			MemberPosition: int(groupMembers.RoleLevel),
			MemberNickName: groupMembers.Nickname,
			MemberId:       groupMembers.UserID,
			JoinTime:      utils.UnixSecondToTime(groupMembers.JoinTime).String(),
		})
	}
	openIMHttp.RespHttp200(c, constant.OK, resp)
}


func AddGroupMembers(c *gin.Context) {
	var (
		req cms_api_struct.RemoveGroupMembersRequest
		resp cms_api_struct.RemoveGroupMembersResponse
		reqPb pbGroup.AddGroupMembersCMSReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(reqPb.OperationId, utils.GetSelfFuncName(),"BindJSON failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.UserIds = req.Members
	reqPb.GroupId = req.GroupId
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pbGroup.NewGroupClient(etcdConn)
	respPb, err := client.AddGroupMembersCMS(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationId, utils.GetSelfFuncName(), "AddGroupMembersCMS failed", err.Error())
		openIMHttp.RespHttp200(c, err, nil)
		return
	}
	resp.Success = respPb.Success
	resp.Failed = respPb.Failed
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func RemoveGroupMembers(c *gin.Context) {
	var (
		req cms_api_struct.RemoveGroupMembersRequest
		resp cms_api_struct.RemoveGroupMembersResponse
		reqPb pbGroup.RemoveGroupMembersCMSReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(),"BindJSON failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.UserIds = req.Members
	reqPb.GroupId = req.GroupId
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pbGroup.NewGroupClient(etcdConn)
	respPb, err := client.RemoveGroupMembersCMS(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "RemoveGroupMembersCMS failed", err.Error())
		openIMHttp.RespHttp200(c, err, nil)
		return
	}
	resp.Success = respPb.Success
	resp.Failed = respPb.Failed
	openIMHttp.RespHttp200(c, constant.OK, resp)
}

func DeleteGroup(c *gin.Context) {
	var (
		req cms_api_struct.DeleteGroupRequest
		_ cms_api_struct.DeleteGroupResponse
		reqPb pbGroup.DeleteGroupReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(),"BindJSON failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.GroupId = req.GroupId
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pbGroup.NewGroupClient(etcdConn)
	_, err := client.DeleteGroup(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "DeleteGroup failed", err.Error())
		openIMHttp.RespHttp200(c, err, nil)
		return
	}
	openIMHttp.RespHttp200(c, constant.OK, nil)
}

func SetGroupMaster(c *gin.Context) {
	var (
		req cms_api_struct.SetGroupMasterRequest
		_ cms_api_struct.SetGroupMasterResponse
		reqPb pbGroup.OperateUserRoleReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(),"BindJSON failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.GroupId = req.GroupId
	reqPb.UserId = req.UserId
	reqPb.RoleLevel = constant.GroupOwner
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
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
		req cms_api_struct.SetGroupMemberRequest
		_ cms_api_struct.AdminLoginResponse
		reqPb pbGroup.OperateUserRoleReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "BindJSON failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.GroupId = req.GroupId
	reqPb.UserId = req.UserId
	reqPb.RoleLevel = constant.GroupOrdinaryUsers
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
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
		req cms_api_struct.AlterGroupInfoRequest
		_ cms_api_struct.SetGroupMasterResponse
		reqPb pbGroup.SetGroupInfoReq
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "BindJSON failed ", err.Error())
		openIMHttp.RespHttp200(c, constant.ErrArgs, nil)
		return
	}
	reqPb.OpUserID = c.MustGet("userID").(string)
	reqPb.GroupInfo = &commonPb.GroupInfo{
		GroupID:       req.GroupID,
		GroupName:     req.GroupName,
		Introduction:  req.Introduction,
		Notification:  req.Notification,
		FaceURL:       req.ProfilePhoto,
		GroupType:     int32(req.GroupType),
	}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName)
	client := pbGroup.NewGroupClient(etcdConn)
	_, err := client.SetGroupInfo(context.Background(), &reqPb)
	if err != nil {
		log.NewError(reqPb.OperationID, utils.GetSelfFuncName(), "DeleteGroup failed", err.Error())
		openIMHttp.RespHttp200(c, err, nil)
		return
	}
	openIMHttp.RespHttp200(c, constant.OK, nil)
}