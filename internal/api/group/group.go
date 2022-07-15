package group

import (
	api "Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	rpc "Open_IM/pkg/proto/group"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"github.com/golang/protobuf/ptypes/wrappers"

	"github.com/gin-gonic/gin"

	"net/http"
	"strings"

	jsonData "Open_IM/internal/utils"
)

// @Summary 把用户踢出群组
// @Description 把用户踢出群组
// @Tags 群组相关
// @ID KickGroupMember
// @Accept json
// @Param token header string true "im token"
// @Param req body api.KickGroupMemberReq true "GroupID为要操作的群ID <br> kickedUserIDList为要踢出的群用户ID <br> reason为原因"
// @Produce json
// @Success 0 {object} api.KickGroupMemberResp "result为结果码, -1为失败, 0为成功"
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /group/kick_group [post]
func KickGroupMember(c *gin.Context) {
	params := api.KickGroupMemberReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}

	req := &rpc.KickGroupMemberReq{}
	utils.CopyStructFields(req, &params)

	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.OperationID, "KickGroupMember args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)
	RpcResp, err := client.KickGroupMember(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "GetGroupMemberList failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	var memberListResp api.KickGroupMemberResp
	memberListResp.ErrMsg = RpcResp.ErrMsg
	memberListResp.ErrCode = RpcResp.ErrCode
	for _, v := range RpcResp.Id2ResultList {
		memberListResp.UserIDResultList = append(memberListResp.UserIDResultList, &api.UserIDResult{UserID: v.UserID, Result: v.Result})
	}
	if len(memberListResp.UserIDResultList) == 0 {
		memberListResp.UserIDResultList = []*api.UserIDResult{}
	}

	log.NewInfo(req.OperationID, "KickGroupMember api return ", memberListResp)
	c.JSON(http.StatusOK, memberListResp)
}

// @Summary 获取群成员信息
// @Description 获取群成员信息
// @Tags 群组相关
// @ID GetGroupMembersInfo
// @Accept json
// @Param token header string true "im token"
// @Param req body api.GetGroupMembersInfoReq true "groupID为要获取的群ID <br> memberList为要获取群成员的群ID列表"
// @Produce json
// @Success 0 {object} api.GetGroupMembersInfoResp{data=[]open_im_sdk.GroupMemberFullInfo}
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /group/get_group_members_info [post]
func GetGroupMembersInfo(c *gin.Context) {
	params := api.GetGroupMembersInfoReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetGroupMembersInfoReq{}
	utils.CopyStructFields(req, params)

	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.OperationID, "GetGroupMembersInfo args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)

	RpcResp, err := client.GetGroupMembersInfo(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "GetGroupMemberList failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	memberListResp := api.GetGroupMembersInfoResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}, MemberList: RpcResp.MemberList}
	memberListResp.Data = jsonData.JsonDataList(RpcResp.MemberList)
	log.NewInfo(req.OperationID, "GetGroupMembersInfo api return ", memberListResp)
	c.JSON(http.StatusOK, memberListResp)
}

func GetGroupMemberList(c *gin.Context) {
	params := api.GetGroupMemberListReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetGroupMemberListReq{}
	utils.CopyStructFields(req, params)

	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.OperationID, "GetGroupMemberList args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)

	RpcResp, err := client.GetGroupMemberList(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "GetGroupMemberList failed, ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	memberListResp := api.GetGroupMemberListResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}, MemberList: RpcResp.MemberList, NextSeq: RpcResp.NextSeq}
	memberListResp.Data = jsonData.JsonDataList(memberListResp.MemberList)

	log.NewInfo(req.OperationID, "GetGroupMemberList api return ", memberListResp)
	c.JSON(http.StatusOK, memberListResp)
}

// @Summary 获取全部群成员列表
// @Description 获取全部群成员列表
// @Tags 群组相关
// @ID GetGroupAllMemberList
// @Accept json
// @Param token header string true "im token"
// @Param req body api.GetGroupAllMemberReq true "GroupID为要获取群成员的群ID"
// @Produce json
// @Success 0 {object} api.GetGroupAllMemberResp{data=[]open_im_sdk.GroupMemberFullInfo}
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /group/get_group_all_member_list [post]
func GetGroupAllMemberList(c *gin.Context) {
	params := api.GetGroupAllMemberReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetGroupAllMemberReq{}
	utils.CopyStructFields(req, &params)

	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.OperationID, "GetGroupAllMember args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)
	RpcResp, err := client.GetGroupAllMember(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "GetGroupAllMember failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	memberListResp := api.GetGroupAllMemberResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}, MemberList: RpcResp.MemberList}
	memberListResp.Data = jsonData.JsonDataList(memberListResp.MemberList)
	log.NewInfo(req.OperationID, "GetGroupAllMember api return ", memberListResp)
	c.JSON(http.StatusOK, memberListResp)
}

// @Summary 获取用户加入群列表
// @Description 获取用户加入群列表
// @Tags 群组相关
// @ID GetJoinedGroupList
// @Accept json
// @Param token header string true "im token"
// @Param req body api.GetJoinedGroupListReq true "fromUserID为要获取的用户ID"
// @Produce json
// @Success 0 {object} api.GetJoinedGroupListResp{data=[]open_im_sdk.GroupInfo}
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /group/get_joined_group_list [post]
func GetJoinedGroupList(c *gin.Context) {
	params := api.GetJoinedGroupListReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetJoinedGroupListReq{}
	utils.CopyStructFields(req, params)

	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.OperationID, "GetJoinedGroupList args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)
	RpcResp, err := client.GetJoinedGroupList(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "GetJoinedGroupList failed  ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	GroupListResp := api.GetJoinedGroupListResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}, GroupInfoList: RpcResp.GroupList}
	GroupListResp.Data = jsonData.JsonDataList(GroupListResp.GroupInfoList)
	log.NewInfo(req.OperationID, "GetJoinedGroupList api return ", GroupListResp)
	c.JSON(http.StatusOK, GroupListResp)
}

// @Summary 将用户拉入群组
// @Description 将用户拉入群组
// @Tags 群组相关
// @ID InviteUserToGroup
// @Accept json
// @Param token header string true "im token"
// @Param req body api.InviteUserToGroupReq true "groupID为要拉进的群组ID <br> invitedUserIDList为要获取群成员的群ID列表 <br> reason为原因"
// @Produce json
// @Success 0 {object} api.InviteUserToGroupResp "result为结果码, -1为失败, 0为成功""
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /group/invite_user_to_group [post]
func InviteUserToGroup(c *gin.Context) {
	params := api.InviteUserToGroupReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.InviteUserToGroupReq{}
	utils.CopyStructFields(req, &params)
	if len(req.InvitedUserIDList) > constant.MaxNotificationNum {
		errMsg := req.OperationID + " too many, Limit: " + utils.IntToString(constant.MaxNotificationNum)
		log.NewError(req.OperationID, errMsg, len(req.InvitedUserIDList))
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.OperationID, "InviteUserToGroup args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)
	RpcResp, err := client.InviteUserToGroup(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "InviteUserToGroup failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	resp := api.InviteUserToGroupResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}}
	for _, v := range RpcResp.Id2ResultList {
		resp.UserIDResultList = append(resp.UserIDResultList, &api.UserIDResult{UserID: v.UserID, Result: v.Result})
	}

	if len(resp.UserIDResultList) == 0 {
		resp.UserIDResultList = *new([]*api.UserIDResult)
	}

	log.NewInfo(req.OperationID, "InviteUserToGroup api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary 创建群组
// @Description 创建群组
// @Tags 群组相关
// @ID CreateGroup
// @Accept json
// @Param token header string true "im token"
// @Param req body api.CreateGroupReq true "groupType这里填0代表普通群 <br>groupName为群名称<br> introduction为群介绍<br> notification为群公共<br>ownerUserID为群主ID <br> ex为群扩展字段 <br> memberList中对象roleLevel为群员角色,1为普通用户 2为群主 3为管理员"
// @Produce json
// @Success 0 {object} api.CreateGroupResp{data=open_im_sdk.GroupInfo}
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /group/create_group [post]
func CreateGroup(c *gin.Context) {
	params := api.CreateGroupReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	//
	req := &rpc.CreateGroupReq{GroupInfo: &open_im_sdk.GroupInfo{}}
	utils.CopyStructFields(req.GroupInfo, &params)

	for _, v := range params.MemberList {
		req.InitMemberList = append(req.InitMemberList, &rpc.GroupAddMemberInfo{UserID: v.UserID, RoleLevel: v.RoleLevel})
	}

	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}

	req.OwnerUserID = params.OwnerUserID
	req.OperationID = params.OperationID

	log.NewInfo(req.OperationID, "CreateGroup args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)
	RpcResp, err := client.CreateGroup(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "CreateGroup failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}

	resp := api.CreateGroupResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}}
	if RpcResp.ErrCode == 0 {
		utils.CopyStructFields(&resp.GroupInfo, RpcResp.GroupInfo)
		resp.Data = jsonData.JsonDataOne(&resp.GroupInfo)
	}
	log.NewInfo(req.OperationID, "CreateGroup api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary 获取用户收到的加群信息列表
// @Description 获取用户收到的加群信息列表
// @Tags 群组相关
// @ID GetRecvGroupApplicationList
// @Accept json
// @Param token header string true "im token"
// @Param req body api.GetGroupApplicationListReq true "fromUserID为要获取的用户ID"
// @Produce json
// @Success 0 {object} api.GetGroupApplicationListResp{data=[]open_im_sdk.GroupRequest}
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /group/get_recv_group_applicationList [post]
func GetRecvGroupApplicationList(c *gin.Context) {
	params := api.GetGroupApplicationListReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetGroupApplicationListReq{}
	utils.CopyStructFields(req, params)
	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	log.NewInfo(req.OperationID, "GetGroupApplicationList args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)
	reply, err := client.GetGroupApplicationList(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "GetGroupApplicationList failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	resp := api.GetGroupApplicationListResp{CommResp: api.CommResp{ErrCode: reply.ErrCode, ErrMsg: reply.ErrMsg}, GroupRequestList: reply.GroupRequestList}
	resp.Data = jsonData.JsonDataList(resp.GroupRequestList)
	log.NewInfo(req.OperationID, "GetGroupApplicationList api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary 获取用户加群申请列表
// @Description 获取用户加群申请列表
// @Tags 群组相关
// @ID GetUserReqGroupApplicationList
// @Accept json
// @Param token header string true "im token"
// @Param req body api.GetUserReqGroupApplicationListReq true "userID为要获取的用户ID"
// @Produce json
// @Success 0 {object} api.GetGroupApplicationListResp{data=[]open_im_sdk.GroupRequest}
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /group/get_user_req_group_applicationList [post]
func GetUserReqGroupApplicationList(c *gin.Context) {
	var params api.GetUserReqGroupApplicationListReq
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", utils.GetSelfFuncName(), err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetUserReqApplicationListReq{}
	utils.CopyStructFields(req, params)
	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}
	log.NewInfo(req.OperationID, "GetGroupsInfo args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)
	RpcResp, err := client.GetUserReqApplicationList(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "GetGroupsInfo failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	log.NewInfo(req.OperationID, RpcResp)
	resp := api.GetGroupApplicationListResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}, GroupRequestList: RpcResp.GroupRequestList}
	resp.Data = jsonData.JsonDataList(resp.GroupRequestList)
	log.NewInfo(req.OperationID, "GetGroupApplicationList api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary 通过群ID列表获取群信息
// @Description 通过群ID列表获取群信息
// @Tags 群组相关
// @ID GetGroupsInfo
// @Accept json
// @Param token header string true "im token"
// @Param req body api.GetGroupInfoReq true "groupIDList为群ID列表"
// @Produce json
// @Success 0 {object} api.GetGroupInfoResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /group/get_groups_info [post]
func GetGroupsInfo(c *gin.Context) {
	params := api.GetGroupInfoReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.GetGroupsInfoReq{}
	utils.CopyStructFields(req, &params)

	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.OperationID, "GetGroupsInfo args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)
	RpcResp, err := client.GetGroupsInfo(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "GetGroupsInfo failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}

	resp := api.GetGroupInfoResp{CommResp: api.CommResp{ErrCode: RpcResp.ErrCode, ErrMsg: RpcResp.ErrMsg}, GroupInfoList: RpcResp.GroupInfoList}
	resp.Data = jsonData.JsonDataList(resp.GroupInfoList)
	log.NewInfo(req.OperationID, "GetGroupsInfo api return ", resp)
	c.JSON(http.StatusOK, resp)
}

//func transferGroupInfo(input []*open_im_sdk.GroupInfo) []*api.GroupInfoAlias {
//	var result []*api.GroupInfoAlias
//	for _, v := range input {
//		t := &api.GroupInfoAlias{}
//		utils.CopyStructFields(t, &v)
//		if v.NeedVerification != nil {
//			t.NeedVerification = v.NeedVerification.Value
//		}
//		result = append(result, t)
//	}
//	return result
//}

//process application

// @Summary 处理加群消息
// @Description 处理加群消息
// @Tags 群组相关
// @ID ApplicationGroupResponse
// @Accept json
// @Param token header string true "im token"
// @Param req body api.ApplicationGroupResponseReq true "groupID为要处理的群ID <br> fromUserID为要处理的用户ID <br> handleMsg为处理结果信息 <br> handleResult为处理结果 1为同意加群 2为拒绝加群"
// @Produce json
// @Success 0 {object} api.ApplicationGroupResponseResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /group/group_application_response [post]
func ApplicationGroupResponse(c *gin.Context) {
	params := api.ApplicationGroupResponseReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.GroupApplicationResponseReq{}
	utils.CopyStructFields(req, &params)

	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.OperationID, "ApplicationGroupResponse args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)
	reply, err := client.GroupApplicationResponse(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "GroupApplicationResponse failed ", req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	resp := api.ApplicationGroupResponseResp{CommResp: api.CommResp{ErrCode: reply.CommonResp.ErrCode, ErrMsg: reply.CommonResp.ErrMsg}}
	log.NewInfo(req.OperationID, "ApplicationGroupResponse api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary 加入群聊
// @Description 加入群聊
// @Tags 群组相关
// @ID JoinGroup
// @Accept json
// @Param token header string true "im token"
// @Param req body api.JoinGroupReq true "reqMessage为申请进群信息<br>groupID为申请的群ID"
// @Produce json
// @Success 0 {object} api.JoinGroupResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /group/join_group [post]
func JoinGroup(c *gin.Context) {
	params := api.JoinGroupReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.JoinGroupReq{}
	utils.CopyStructFields(req, params)

	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.OperationID, "JoinGroup args ", req.String())
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)

	RpcResp, err := client.JoinGroup(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "JoinGroup failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	resp := api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}
	log.NewInfo(req.OperationID, "JoinGroup api return", RpcResp.String())
	c.JSON(http.StatusOK, resp)
}

// @Summary 当前用户退出群聊
// @Description 当前用户退出群聊
// @Tags 群组相关
// @ID QuitGroup
// @Accept json
// @Param token header string true "im token"
// @Param req body api.QuitGroupReq true "groupID为要退出的群ID"
// @Produce json
// @Success 0 {object} api.QuitGroupResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /group/quit_group [post]
func QuitGroup(c *gin.Context) {
	params := api.QuitGroupReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.QuitGroupReq{}
	utils.CopyStructFields(req, &params)

	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.OperationID, "QuitGroup args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)
	RpcResp, err := client.QuitGroup(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "call quit group rpc server failed,err=%s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	resp := api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}
	log.NewInfo(req.OperationID, "QuitGroup api return", RpcResp.String())
	c.JSON(http.StatusOK, resp)
}

// @Summary 设置群信息
// @Description 设置群信息
// @Tags 群组相关
// @ID SetGroupInfo
// @Accept json
// @Param token header string true "im token"
// @Param req body api.SetGroupInfoReq true "groupID为要修改的群ID<br>groupName为新的群名称<br>notification为群介绍 <br> introduction为群公告 <br> needVerification为加群验证 0为申请需要同意 邀请直接进 1为所有人进群需要验证，除了群主管理员邀请进群 2为直接进群"
// @Produce json
// @Success 0 {object} api.SetGroupInfoResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /group/set_group_info [post]
func SetGroupInfo(c *gin.Context) {
	params := api.SetGroupInfoReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.SetGroupInfoReq{GroupInfoForSet: &open_im_sdk.GroupInfoForSet{}}
	utils.CopyStructFields(req.GroupInfoForSet, &params)
	req.OperationID = params.OperationID
	argsHandle(&params, req)
	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.OperationID, "SetGroupInfo args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)
	RpcResp, err := client.SetGroupInfo(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "SetGroupInfo failed ", err.Error(), req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": "call  rpc server failed"})
		return
	}
	resp := api.SetGroupInfoResp{CommResp: api.CommResp{ErrCode: RpcResp.CommonResp.ErrCode, ErrMsg: RpcResp.CommonResp.ErrMsg}}
	c.JSON(http.StatusOK, resp)
	log.NewInfo(req.OperationID, "SetGroupInfo api return ", resp)
}
func argsHandle(params *api.SetGroupInfoReq, req *rpc.SetGroupInfoReq) {
	if params.NeedVerification != nil {
		req.GroupInfoForSet.NeedVerification = &wrappers.Int32Value{Value: *params.NeedVerification}
		log.NewInfo(req.OperationID, "NeedVerification ", req.GroupInfoForSet.NeedVerification)
	}
	if params.LookMemberInfo != nil {
		req.GroupInfoForSet.LookMemberInfo = &wrappers.Int32Value{Value: *params.LookMemberInfo}
		log.NewInfo(req.OperationID, "LookMemberInfo ", req.GroupInfoForSet.LookMemberInfo)
	}
	if params.ApplyMemberFriend != nil {
		req.GroupInfoForSet.ApplyMemberFriend = &wrappers.Int32Value{Value: *params.ApplyMemberFriend}
		log.NewInfo(req.OperationID, "ApplyMemberFriend ", req.GroupInfoForSet.ApplyMemberFriend)
	}
}

// @Summary 转让群主
// @Description 转让群主
// @Tags 群组相关
// @ID TransferGroupOwner
// @Accept json
// @Param token header string true "im token"
// @Param req body api.TransferGroupOwnerReq true "GroupID为要操作的群ID <br> oldOwnerUserID为老群主ID <br> newOwnerUserID为新群主ID"
// @Produce json
// @Success 0 {object} api.TransferGroupOwnerResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /group/transfer_group [post]
func TransferGroupOwner(c *gin.Context) {
	params := api.TransferGroupOwnerReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.TransferGroupOwnerReq{}
	utils.CopyStructFields(req, &params)

	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.OperationID, "TransferGroupOwner args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)
	reply, err := client.TransferGroupOwner(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, "TransferGroupOwner failed ", req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	resp := api.TransferGroupOwnerResp{CommResp: api.CommResp{ErrCode: reply.CommonResp.ErrCode, ErrMsg: reply.CommonResp.ErrMsg}}
	log.NewInfo(req.OperationID, "TransferGroupOwner api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary 解散群组
// @Description 解散群组
// @Tags 群组相关
// @ID DismissGroup
// @Accept json
// @Param token header string true "im token"
// @Param req body api.DismissGroupReq true "groupID为要解散的群组ID"
// @Produce json
// @Success 0 {object} api.DismissGroupResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /group/dismiss_group [post]
func DismissGroup(c *gin.Context) {
	params := api.DismissGroupReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.DismissGroupReq{}
	utils.CopyStructFields(req, &params)

	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)
	reply, err := client.DismissGroup(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), " failed ", req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	resp := api.DismissGroupResp{CommResp: api.CommResp{ErrCode: reply.CommonResp.ErrCode, ErrMsg: reply.CommonResp.ErrMsg}}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary 禁言群成员
// @Description 禁言群成员
// @Tags 群组相关
// @ID MuteGroupMember
// @Accept json
// @Param token header string true "im token"
// @Param req body api.MuteGroupMemberReq true "groupID为群组ID <br> userID为要禁言的用户ID <br> mutedSeconds为禁言秒数"
// @Produce json
// @Success 0 {object} api.DismissGroupResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /group/mute_group_member [post]
func MuteGroupMember(c *gin.Context) {
	params := api.MuteGroupMemberReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.MuteGroupMemberReq{}
	utils.CopyStructFields(req, &params)

	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " api args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)
	reply, err := client.MuteGroupMember(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), " failed ", req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	resp := api.MuteGroupMemberResp{CommResp: api.CommResp{ErrCode: reply.CommonResp.ErrCode, ErrMsg: reply.CommonResp.ErrMsg}}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary 取消禁言群成员
// @Description 取消禁言群成员
// @Tags 群组相关
// @ID CancelMuteGroupMember
// @Accept json
// @Param token header string true "im token"
// @Param req body api.CancelMuteGroupMemberReq true "groupID为群组ID <br> userID为要取消禁言的用户ID"
// @Produce json
// @Success 0 {object} api.CancelMuteGroupMemberResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /group/cancel_mute_group_member [post]
func CancelMuteGroupMember(c *gin.Context) {
	params := api.CancelMuteGroupMemberReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.CancelMuteGroupMemberReq{}
	utils.CopyStructFields(req, &params)

	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " api args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)
	reply, err := client.CancelMuteGroupMember(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), " failed ", req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	resp := api.CancelMuteGroupMemberResp{CommResp: api.CommResp{ErrCode: reply.CommonResp.ErrCode, ErrMsg: reply.CommonResp.ErrMsg}}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary 禁言群组
// @Description 禁言群组
// @Tags 群组相关
// @ID MuteGroup
// @Accept json
// @Param token header string true "im token"
// @Param req body api.MuteGroupReq true "groupID为群组ID"
// @Produce json
// @Success 0 {object} api.MuteGroupResp
// @Failure 500 {object} api.MuteGroupResp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.MuteGroupResp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /group/mute_group [post]
func MuteGroup(c *gin.Context) {
	params := api.MuteGroupReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.MuteGroupReq{}
	utils.CopyStructFields(req, &params)

	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " api args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)
	reply, err := client.MuteGroup(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), " failed ", req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	resp := api.MuteGroupResp{CommResp: api.CommResp{ErrCode: reply.CommonResp.ErrCode, ErrMsg: reply.CommonResp.ErrMsg}}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary 取消禁言群组
// @Description 取消禁言群组
// @Tags 群组相关
// @ID CancelMuteGroup
// @Accept json
// @Param token header string true "im token"
// @Param req body api.CancelMuteGroupReq true "groupID为群组ID"
// @Produce json
// @Success 0 {object} api.CancelMuteGroupResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /group/cancel_mute_group [post]
func CancelMuteGroup(c *gin.Context) {
	params := api.CancelMuteGroupReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.CancelMuteGroupReq{}
	utils.CopyStructFields(req, &params)

	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " api args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)
	reply, err := client.CancelMuteGroup(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), " failed ", req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	resp := api.CancelMuteGroupResp{CommResp: api.CommResp{ErrCode: reply.CommonResp.ErrCode, ErrMsg: reply.CommonResp.ErrMsg}}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " api return ", resp)
	c.JSON(http.StatusOK, resp)
}

//SetGroupMemberNickname

func SetGroupMemberNickname(c *gin.Context) {
	params := api.SetGroupMemberNicknameReq{}
	if err := c.BindJSON(&params); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	req := &rpc.SetGroupMemberNicknameReq{}
	utils.CopyStructFields(req, &params)

	var ok bool
	var errInfo string
	ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " api args ", req.String())

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)
	reply, err := client.SetGroupMemberNickname(context.Background(), req)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), " failed ", req.String())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	resp := api.SetGroupMemberNicknameResp{CommResp: api.CommResp{ErrCode: reply.CommonResp.ErrCode, ErrMsg: reply.CommonResp.ErrMsg}}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " api return ", resp)
	c.JSON(http.StatusOK, resp)
}

// @Summary 修改群成员信息
// @Description 修改群成员信息
// @Tags 群组相关
// @ID SetGroupMemberInfo
// @Accept json
// @Param token header string true "im token"
// @Param req body api.SetGroupMemberInfoReq true "除了operationID, userID, groupID其他参数可选<br>ex为拓展字段<br>faceURL为群头像<br>nickName为群昵称<br>roleLevel为群员角色,1为普通用户 2为群主 3为管理员"
// @Produce json
// @Success 0 {object} api.SetGroupMemberInfoResp
// @Failure 500 {object} api.Swagger500Resp "errCode为500 一般为服务器内部错误"
// @Failure 400 {object} api.Swagger400Resp "errCode为400 一般为参数输入错误, token未带上等"
// @Router /group/set_group_member_info [post]
func SetGroupMemberInfo(c *gin.Context) {
	var (
		req  api.SetGroupMemberInfoReq
		resp api.SetGroupMemberInfoResp
	)
	if err := c.BindJSON(&req); err != nil {
		log.NewError("0", "BindJSON failed ", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 400, "errMsg": err.Error()})
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), req)
	var opUserID string
	ok, opUserID, errInfo := token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	if !ok {
		errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusBadRequest, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}

	reqPb := &rpc.SetGroupMemberInfoReq{
		GroupID:     req.GroupID,
		UserID:      req.UserID,
		OperationID: req.OperationID,
		OpUserID:    opUserID,
	}
	if req.Nickname != nil {
		reqPb.Nickname = &wrappers.StringValue{Value: *req.Nickname}
	}
	if req.FaceURL != nil {
		reqPb.FaceURL = &wrappers.StringValue{Value: *req.FaceURL}
	}
	if req.Ex != nil {
		reqPb.Ex = &wrappers.StringValue{Value: *req.Ex}
	}
	if req.RoleLevel != nil {
		reqPb.RoleLevel = &wrappers.Int32Value{Value: *req.RoleLevel}
	}

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImGroupName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
		return
	}
	client := rpc.NewGroupClient(etcdConn)
	respPb, err := client.SetGroupMemberInfo(context.Background(), reqPb)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), " failed ", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": err.Error()})
		return
	}

	resp.ErrMsg = respPb.CommonResp.ErrMsg
	resp.ErrCode = respPb.CommonResp.ErrCode
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " api args ", resp)
	c.JSON(http.StatusOK, resp)
}
