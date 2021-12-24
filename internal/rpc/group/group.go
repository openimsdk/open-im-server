package group

import (
	chat "Open_IM/internal/rpc/msg"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbGroup "Open_IM/pkg/proto/group"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"strings"
	"time"
)

type groupServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func NewGroupServer(port int) *groupServer {
	log.NewPrivateLog("group")
	return &groupServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImGroupName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (s *groupServer) Run() {
	log.Info("", "", "rpc group init....")

	ip := utils.ServerIP
	registerAddress := ip + ":" + strconv.Itoa(s.rpcPort)
	//listener network
	listener, err := net.Listen("tcp", registerAddress)
	if err != nil {
		log.InfoByArgs("listen network failed,err=%s", err.Error())
		return
	}
	log.Info("", "", "listen network success, address = %s", registerAddress)
	defer listener.Close()
	//grpc server
	srv := grpc.NewServer()
	defer srv.GracefulStop()
	//Service registers with etcd
	pbGroup.RegisterGroupServer(srv, s)
	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), ip, s.rpcPort, s.rpcRegisterName, 10)
	if err != nil {
		log.ErrorByArgs("get etcd failed,err=%s", err.Error())
		return
	}
	err = srv.Serve(listener)
	if err != nil {
		log.ErrorByArgs("listen rpc_group error,err=%s", err.Error())
		return
	}
	log.Info("", "", "rpc create group init success")
}

func (s *groupServer) CreateGroup(ctx context.Context, req *pbGroup.CreateGroupReq) (*pbGroup.CreateGroupResp, error) {
	log.NewInfo(req.OperationID, "CreateGroup, args ", req.String())
	if !token_verify.CheckAccess(req.OpUserID, req.FromUserID) {
		log.NewError(req.OperationID, "CheckAccess false ", req.OpUserID, req.FromUserID)
		return &pbGroup.CreateGroupResp{ErrCode: constant.ErrParseToken.ErrCode, ErrMsg: constant.ErrParseToken.ErrMsg}, nil
	}
	//Time stamp + MD5 to generate group chat id
	groupId := utils.Md5(strconv.FormatInt(time.Now().UnixNano(), 10))
	//to group
	err := im_mysql_model.InsertIntoGroup(groupId, req.GroupName, req.Introduction, req.Notification, req.FaceUrl, req.Ext)
	if err != nil {
		log.NewError(req.OperationID, "InsertIntoGroup failed, ", err.Error(), groupId, req.GroupName, req.Introduction, req.Notification, req.FaceUrl, req.Ext)
		return &pbGroup.CreateGroupResp{ErrCode: constant.ErrCreateGroup.ErrCode, ErrMsg: constant.ErrCreateGroup.ErrMsg}, nil
	}

	us, err := im_mysql_model.FindUserByUID(req.FromUserID)
	if err != nil {
		log.NewError(req.OperationID, "FindUserByUID failed ", err.Error(), req.FromUserID)
		return &pbGroup.CreateGroupResp{ErrCode: constant.ErrCreateGroup.ErrCode, ErrMsg: constant.ErrCreateGroup.ErrMsg}, nil
	}

	//to group member
	err = im_mysql_model.InsertIntoGroupMember(groupId, us.UserID, us.Nickname, us.FaceUrl, constant.GroupOwner)
	if err != nil {
		log.NewError(req.OperationID, "InsertIntoGroupMember failed ", err.Error())
		return &pbGroup.CreateGroupResp{ErrCode: constant.ErrCreateGroup.ErrCode, ErrMsg: constant.ErrCreateGroup.ErrMsg}, nil
	}

	err = db.DB.AddGroupMember(groupId, req.FromUserID)
	if err != nil {
		log.NewError(req.OperationID, "AddGroupMember failed ", err.Error(), groupId, req.FromUserID)
		//	return &pbGroup.CreateGroupResp{ErrCode: constant.ErrCreateGroup.ErrCode, ErrMsg: constant.ErrCreateGroup.ErrMsg}, nil
	}

	//to group member
	for _, user := range req.InitMemberList {
		us, err := im_mysql_model.FindUserByUID(user.UserID)
		if err != nil {
			log.NewError(req.OperationID, "FindUserByUID failed ", err.Error(), user.UserID)
			continue
		}
		if user.Role == 1 {
			log.NewError(req.OperationID, "only one owner, failed ", user)
			continue
		}
		err = im_mysql_model.InsertIntoGroupMember(groupId, user.UserID, us.Nickname, us.FaceUrl, user.Role)
		if err != nil {
			log.NewError(req.OperationID, "InsertIntoGroupMember failed ", groupId, user.UserID, us.Nickname, us.FaceUrl, user.Role)
		}
		err = db.DB.AddGroupMember(groupId, user.UserID)
		if err != nil {
			log.NewError(req.OperationID, "add mongo group member failed, db.DB.AddGroupMember failed ", err.Error())
		}
	}

	resp := &pbGroup.CreateGroupResp{}
	group, err := im_mysql_model.FindGroupInfoByGroupId(groupId)
	if err != nil {
		log.NewError(req.OperationID, "FindGroupInfoByGroupId failed ", err.Error(), groupId)
		resp.ErrCode = constant.ErrCreateGroup.ErrCode
		resp.ErrMsg = constant.ErrCreateGroup.ErrMsg
		return resp, nil
	}
	chat.GroupCreatedNotification(req, groupId)
	utils.CopyStructFields(resp.GroupInfo, group)
	log.NewInfo(req.OperationID, "rpc CreateGroup return ", resp.String())
	return resp, nil
}

func (s *groupServer) GetJoinedGroupList(ctx context.Context, req *pbGroup.GetJoinedGroupListReq) (*pbGroup.GetJoinedGroupListResp, error) {
	log.NewInfo(req.OperationID, "GetJoinedGroupList, args ", req.String())
	if !token_verify.CheckAccess(req.OpUserID, req.FromUserID) {
		log.NewError(req.OperationID, "CheckAccess false ", req.OpUserID, req.FromUserID)
		return &pbGroup.GetJoinedGroupListResp{ErrCode: constant.ErrParseToken.ErrCode, ErrMsg: constant.ErrParseToken.ErrMsg}, nil
	}

	//group list
	joinedGroupList, err := imdb.GetJoinedGroupIdListByMemberId(req.FromUserID)
	if err != nil {
		log.NewError(req.OperationID, "GetJoinedGroupIdListByMemberId failed ", err.Error(), req.FromUserID)
		return &pbGroup.GetJoinedGroupListResp{ErrCode: constant.ErrParam.ErrCode, ErrMsg: constant.ErrParam.ErrMsg}, nil
	}

	var resp pbGroup.GetJoinedGroupListResp
	for _, v := range joinedGroupList {
		var groupNode open_im_sdk.GroupInfo
		num := imdb.GetGroupMemberNumByGroupId(v.GroupID)
		owner, err2 := imdb.GetGroupOwnerInfoByGroupId(v.GroupID)
		group, err := imdb.FindGroupInfoByGroupId(v.GroupID)
		if num > 0 && owner != nil && err2 == nil && group != nil && err == nil {
			utils.CopyStructFields(&groupNode, group)
			groupNode.CreateTime = group.CreateTime
			utils.CopyStructFields(groupNode.Owner, owner)
			groupNode.MemberCount = uint32(num)
			resp.GroupList = append(resp.GroupList, &groupNode)
		} else {
			log.NewError(req.OperationID, "check nil ", num, owner, err, group)
			continue
		}
		log.NewDebug(req.OperationID, "joinedGroup ", groupNode)
	}
	resp.ErrCode = 0
	log.NewInfo(req.OperationID, "GetJoinedGroupList return ", resp.String())
	return &resp, nil
}

func (s *groupServer) InviteUserToGroup(ctx context.Context, req *pbGroup.InviteUserToGroupReq) (*pbGroup.InviteUserToGroupResp, error) {
	log.NewInfo(req.OperationID, "InviteUserToGroup args ", req.String())

	if !imdb.IsExistGroupMember(req.GroupID, req.OpUserID) && !token_verify.IsMangerUserID(req.OpUserID) {
		log.NewError(req.OperationID, "no permission InviteUserToGroup ", req.GroupID, req.OpUserID)
		return &pbGroup.InviteUserToGroupResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}, nil
	}

	groupInfoFromMysql, err := imdb.FindGroupInfoByGroupId(req.GroupID)
	if err != nil || groupInfoFromMysql == nil {
		log.NewError(req.OperationID, "FindGroupInfoByGroupId failed ", req.GroupID, err)
		return &pbGroup.InviteUserToGroupResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}, nil
	}
	//
	//from User:  invite: applicant
	//to user:  invite: invited
	var resp pbGroup.InviteUserToGroupResp
	var okUserIDList []string
	for _, v := range req.InvitedUserIDList {
		var resultNode pbGroup.Id2Result
		resultNode.UserID = v
		resultNode.Result = 0
		toUserInfo, err := imdb.FindUserByUID(v)
		if err != nil {
			log.NewError(req.OperationID, "FindUserByUID failed ", err.Error(), v)
			resultNode.Result = -1
			resp.Id2ResultList = append(resp.Id2ResultList, &resultNode)
			continue
		}

		if imdb.IsExistGroupMember(req.GroupID, v) {
			log.NewError(req.OperationID, "IsExistGroupMember ", req.GroupID, v)
			resultNode.Result = -1
			resp.Id2ResultList = append(resp.Id2ResultList, &resultNode)
			continue
		}

		err = imdb.InsertGroupMember(req.GroupID, toUserInfo.UserID, toUserInfo.Nickname, toUserInfo.FaceUrl, 0)
		if err != nil {
			log.NewError(req.OperationID, "InsertGroupMember failed ", req.GroupID, toUserInfo.UserID, toUserInfo.Nickname, toUserInfo.FaceUrl)
			resultNode.Result = -1
			resp.Id2ResultList = append(resp.Id2ResultList, &resultNode)
			continue
		}
		okUserIDList = append(okUserIDList, v)
		err = db.DB.AddGroupMember(req.GroupID, toUserInfo.UserID)
		if err != nil {
			log.NewError(req.OperationID, "AddGroupMember failed ", err.Error(), req.GroupID, toUserInfo.UserID)
		}
		resp.Id2ResultList = append(resp.Id2ResultList, &resultNode)
	}

	chat.MemberInvitedNotification(req.OperationID, req.GroupID, req.OpUserID, req.Reason, okUserIDList)
	resp.ErrCode = 0
	log.NewInfo(req.OperationID, "InviteUserToGroup rpc return ", resp.String())
	return &resp, nil
}

func (s *groupServer) GetGroupAllMember(ctx context.Context, req *pbGroup.GetGroupAllMemberReq) (*pbGroup.GetGroupAllMemberResp, error) {
	log.NewInfo(req.OperationID, "GetGroupAllMember, args ", req.String())
	var resp pbGroup.GetGroupAllMemberResp
	resp.ErrCode = 0
	memberList, err := imdb.FindGroupMemberListByGroupId(req.GroupID)
	if err != nil {
		resp.ErrCode = constant.ErrDb.ErrCode
		resp.ErrMsg = constant.ErrDb.ErrMsg
		log.NewError(req.OperationID, "FindGroupMemberListByGroupId failed,", err.Error(), req.GroupID)
		return &resp, nil
	}
	m := token_verify.IsMangerUserID(req.OpUserID)
	in := false
	if m {
		in = true
	}
	for _, v := range memberList {
		var node open_im_sdk.GroupMemberFullInfo
		utils.CopyStructFields(node, v)
		resp.MemberList = append(resp.MemberList, &node)
		if !m && req.OpUserID == v.UserID {
			in = true
		}
	}
	if !in {

	}
	resp.ErrCode = 0
	return &resp, nil
}

func (s *groupServer) GetGroupMemberList(ctx context.Context, req *pbGroup.GetGroupMemberListReq) (*pbGroup.GetGroupMemberListResp, error) {
	log.NewInfo(req.OperationID, "GetGroupMemberList, args ", req.String())

	var resp pbGroup.GetGroupMemberListResp
	resp.ErrCode = 0
	memberList, err := imdb.GetGroupMemberByGroupId(req.GroupID, req.Filter, req.NextSeq, 30)
	if err != nil {
		resp.ErrCode = constant.ErrDb.ErrCode
		resp.ErrMsg = err.Error()
		log.NewError(req.OperationID, "GetGroupMemberByGroupId failed,", req.GroupID, req.Filter, req.NextSeq, 30)
		return &resp, nil
	}

	for _, v := range memberList {
		var node open_im_sdk.GroupMemberFullInfo
		utils.CopyStructFields(&node, v)
		resp.MemberList = append(resp.MemberList, &node)
	}
	//db operate  get db sorted by join time
	if int32(len(memberList)) < 30 {
		resp.NextSeq = 0
	} else {
		resp.NextSeq = req.NextSeq + int32(len(memberList))
	}

	resp.ErrCode = 0
	log.NewInfo(req.OperationID, "GetGroupMemberList rpc return ", resp)
	return &resp, nil
}

func (s *groupServer) KickGroupMember(ctx context.Context, req *pbGroup.KickGroupMemberReq) (*pbGroup.KickGroupMemberResp, error) {
	log.NewInfo(req.OperationID, "KickGroupMember args ", req.String())
	ownerList, err := imdb.GetOwnerManagerByGroupId(req.GroupID)
	if err != nil {
		log.NewError(req.OperationID, "GetOwnerManagerByGroupId failed ", err.Error(), req.GroupID)
		return &pbGroup.KickGroupMemberResp{ErrCode: constant.ErrParam.ErrCode, ErrMsg: constant.ErrParam.ErrMsg}, nil
	}
	//op is group owner?
	var flag = 0
	for _, v := range ownerList {
		if v.UserID == req.OpUserID {
			flag = 1
			log.NewDebug(req.OperationID, "is group owner ", req.OpUserID, req.GroupID)
			break
		}
	}

	//op is app manager
	if flag != 1 {
		if token_verify.IsMangerUserID(req.OpUserID) {
			flag = 1
			log.NewDebug(req.OperationID, "is app manager ", req.OpUserID)
		}
	}

	if flag != 1 {
		log.NewError(req.OperationID, "failed, no access kick ")
		return &pbGroup.KickGroupMemberResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}, nil
	}

	if len(req.KickedUserIDList) == 0 {
		log.NewError(req.OperationID, "failed, kick list 0")
		return &pbGroup.KickGroupMemberResp{ErrCode: constant.ErrParam.ErrCode, ErrMsg: constant.ErrParam.ErrMsg}, nil
	}

	groupOwnerUserID := ""
	for _, v := range ownerList {
		if v.AdministratorLevel == 1 {
			groupOwnerUserID = v.UserID
		}
	}

	var okUserIDList []string
	//remove
	var resp pbGroup.KickGroupMemberResp
	for _, v := range req.KickedUserIDList {
		//owner cant kicked
		if v == groupOwnerUserID {
			log.NewError(req.OperationID, "failed, can't kick owner ", v)
			resp.Id2ResultList = append(resp.Id2ResultList, &pbGroup.Id2Result{UserID: v, Result: -1})
			continue
		}
		err := imdb.RemoveGroupMember(req.GroupID, v)
		if err != nil {
			log.NewError(req.OperationID, "RemoveGroupMember failed ", err.Error(), req.GroupID, v)
			resp.Id2ResultList = append(resp.Id2ResultList, &pbGroup.Id2Result{UserID: v, Result: -1})
		} else {
			resp.Id2ResultList = append(resp.Id2ResultList, &pbGroup.Id2Result{UserID: v, Result: 0})
			okUserIDList = append(okUserIDList, v)
		}

		err = db.DB.DelGroupMember(req.GroupID, v)
		if err != nil {
			log.NewError(req.OperationID, "DelGroupMember failed ", err.Error(), req.GroupID, v)
		}
	}
	chat.MemberKickedNotification(req, okUserIDList)
	resp.ErrCode = 0
	log.NewInfo(req.OperationID, "GetGroupMemberList rpc return ", resp.String())
	return &resp, nil
}

func (s *groupServer) GetGroupMembersInfo(ctx context.Context, req *pbGroup.GetGroupMembersInfoReq) (*pbGroup.GetGroupMembersInfoResp, error) {
	log.NewInfo(req.OperationID, "GetGroupMembersInfo args ", req.String())

	var resp pbGroup.GetGroupMembersInfoResp

	for _, v := range req.MemberList {
		var memberNode open_im_sdk.GroupMemberFullInfo
		memberInfo, err := imdb.GetMemberInfoById(req.GroupID, v)
		memberNode.UserID = v
		if err != nil {
			log.NewError(req.OperationID, "GetMemberInfoById failed ", err.Error(), req.GroupID, v)
			continue
		} else {
			utils.CopyStructFields(&memberNode, memberInfo)
			memberNode.JoinTime = memberInfo.JoinTime
			resp.MemberList = append(resp.MemberList, &memberNode)
		}
	}
	resp.ErrCode = 0
	log.NewInfo(req.OperationID, "GetGroupMembersInfo rpc return ", resp.String())
	return &resp, nil
}

func (s *groupServer) GetGroupApplicationList(_ context.Context, req *pbGroup.GetGroupApplicationListReq) (*pbGroup.GetGroupApplicationListResp, error) {
	log.NewInfo(req.OperationID, "GetGroupMembersInfo args ", req.String())
	reply, err := im_mysql_model.GetGroupApplicationList(req.OpUserID)
	if err != nil {
		log.NewError(req.OperationID, "GetGroupApplicationList failed ", err.Error(), req.OpUserID)
		return &pbGroup.GetGroupApplicationListResp{ErrCode: 701, ErrMsg: "GetGroupApplicationList failed"}, nil
	}
	log.NewInfo(req.OperationID, "GetGroupMembersInfo rpc return ", reply)
	return reply, nil
}

func (s *groupServer) GetGroupsInfo(ctx context.Context, req *pbGroup.GetGroupsInfoReq) (*pbGroup.GetGroupsInfoResp, error) {
	log.NewInfo(req.OperationID, "GetGroupsInfo args ", req.String())
	groupsInfoList := make([]*open_im_sdk.GroupInfo, 0)
	for _, groupID := range req.GroupIDList {
		groupInfoFromMysql, err := im_mysql_model.FindGroupInfoByGroupId(groupID)
		if err != nil {
			log.NewError(req.OperationID, "FindGroupInfoByGroupId failed ", err.Error(), groupID)
			continue
		}
		var groupInfo open_im_sdk.GroupInfo
		utils.CopyStructFields(&groupInfo, groupInfoFromMysql)
		groupInfo.CreateTime = groupInfoFromMysql.CreateTime
		groupsInfoList = append(groupsInfoList, &groupInfo)
	}

	resp := pbGroup.GetGroupsInfoResp{GroupInfoList: groupsInfoList}
	log.NewInfo(req.OperationID, "GetGroupsInfo rpc return ", resp)
	return &resp, nil
}

func (s *groupServer) GroupApplicationResponse(_ context.Context, req *pbGroup.GroupApplicationResponseReq) (*pbGroup.CommonResp, error) {
	log.NewInfo(req.OperationID, "GroupApplicationResponse args ", req.String())
	reply, err := imdb.GroupApplicationResponse(req)
	if err != nil {
		log.NewError(req.OperationID, "GroupApplicationResponse failed ", err.Error(), req.String())
		return &pbGroup.CommonResp{ErrCode: 702, ErrMsg: err.Error()}, nil
	}

	if req.HandleResult == 1 {
		if req.ToUserID == "0" {
			err = db.DB.AddGroupMember(req.GroupID, req.FromUserID)
			if err != nil {
				log.NewError(req.OperationID, "AddGroupMember failed ", err.Error(), req.GroupID, req.FromUserID)
			}
		} else {
			err = db.DB.AddGroupMember(req.GroupID, req.ToUserID)
			if err != nil {
				log.NewError(req.OperationID, "AddGroupMember failed ", err.Error(), req.GroupID, req.ToUserID)
			}
		}
	}
	if req.ToUserID == "0" {
		//group, err := imdb.FindGroupInfoByGroupId(req.GroupID)
		//if err != nil {
		//	log.NewError(req.OperationID, "FindGroupInfoByGroupId failed ", req.GroupID)
		//	return reply, nil
		//}
		//member, err := imdb.FindGroupMemberInfoByGroupIdAndUserId(req.GroupID, req.OpUserID)
		//if err != nil {
		//	log.NewError(req.OperationID, "FindGroupMemberInfoByGroupIdAndUserId failed ", req.GroupID, req.OpUserID)
		//	return reply, nil
		//}
		chat.ApplicationProcessedNotification(req)
		if req.HandleResult == 1 {
			//	entrantUser, err := imdb.FindGroupMemberInfoByGroupIdAndUserId(req.GroupID, req.FromUserID)
			//	if err != nil {
			//		log.NewError(req.OperationID, "FindGroupMemberInfoByGroupIdAndUserId failed ", err.Error(), req.GroupID, req.FromUserID)
			//	return reply, nil
			//	}
			chat.MemberEnterNotification(req)
		}
	} else {
		log.NewError(req.OperationID, "args failed ", req.String())
	}

	log.NewInfo(req.OperationID, "rpc GroupApplicationResponse ok ", reply)
	return reply, nil
}

func (s *groupServer) JoinGroup(ctx context.Context, req *pbGroup.JoinGroupReq) (*pbGroup.CommonResp, error) {
	log.NewInfo(req.OperationID, "JoinGroup args ", req.String())

	applicationUserInfo, err := im_mysql_model.FindUserByUID(req.OpUserID)
	if err != nil {
		log.NewError(req.OperationID, "FindUserByUID failed ", err.Error(), req.OpUserID)
		return &pbGroup.CommonResp{ErrCode: constant.ErrSearchUserInfo.ErrCode, ErrMsg: constant.ErrSearchUserInfo.ErrMsg}, nil
	}

	_, err = im_mysql_model.FindGroupRequestUserInfoByGroupIDAndUid(req.GroupID, req.OpUserID)
	if err == nil {
		err = im_mysql_model.DelGroupRequest(req.GroupID, req.OpUserID, "0")
	}

	if err = im_mysql_model.InsertIntoGroupRequest(req.GroupID, req.OpUserID, "0", req.ReqMessage, applicationUserInfo.Nickname, applicationUserInfo.FaceUrl); err != nil {
		log.NewError(req.OperationID, "InsertIntoGroupRequest ", err.Error(), req.GroupID, req.OpUserID, "0", req.ReqMessage, applicationUserInfo.Nickname, applicationUserInfo.FaceUrl)
		return &pbGroup.CommonResp{ErrCode: constant.ErrJoinGroupApplication.ErrCode, ErrMsg: constant.ErrJoinGroupApplication.ErrMsg}, nil
	}

	memberList, err := im_mysql_model.FindGroupMemberListByGroupIdAndFilterInfo(req.GroupID, constant.GroupOwner)
	if len(memberList) == 0 {
		log.NewError(req.OperationID, "FindGroupMemberListByGroupIdAndFilterInfo failed ", req.GroupID, constant.GroupOwner, err)
		return &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}, nil
	}

	chat.ReceiveJoinApplicationNotification(req)

	log.NewInfo(req.OperationID, "ReceiveJoinApplicationNotification rpc JoinGroup success return")
	return &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}, nil
}

func (s *groupServer) QuitGroup(ctx context.Context, req *pbGroup.QuitGroupReq) (*pbGroup.CommonResp, error) {
	log.NewError("QuitGroup args ", req.String())

	_, err := im_mysql_model.FindGroupMemberInfoByGroupIdAndUserId(req.GroupID, req.OpUserID)
	if err != nil {
		log.NewError(req.OperationID, "FindGroupMemberInfoByGroupIdAndUserId failed", err.Error(), req.GroupID, req.OpUserID)
		return &pbGroup.CommonResp{ErrCode: constant.ErrQuitGroup.ErrCode, ErrMsg: constant.ErrQuitGroup.ErrMsg}, nil
	}

	err = im_mysql_model.DeleteGroupMemberByGroupIdAndUserId(req.GroupID, req.OpUserID)
	if err != nil {
		log.NewError(req.OperationID, "DeleteGroupMemberByGroupIdAndUserId failed ", err.Error(), req.GroupID, req.OpUserID)
		return &pbGroup.CommonResp{ErrCode: constant.ErrQuitGroup.ErrCode, ErrMsg: constant.ErrQuitGroup.ErrMsg}, nil
	}

	err = db.DB.DelGroupMember(req.GroupID, req.OpUserID)
	if err != nil {
		log.NewError(req.OperationID, "DelGroupMember failed ", req.GroupID, req.OpUserID)
		//	return &pbGroup.CommonResp{ErrorCode: constant.ErrQuitGroup.ErrCode, ErrorMsg: constant.ErrQuitGroup.ErrMsg}, nil
	}

	chat.MemberLeaveNotification(req)
	log.NewInfo(req.OperationID, "rpc quit group is success return")
	return &pbGroup.CommonResp{}, nil
}

func hasAccess(req *pbGroup.SetGroupInfoReq) bool {
	if utils.IsContain(req.OpUserID, config.Config.Manager.AppManagerUid) {
		return true
	}
	groupUserInfo, err := im_mysql_model.FindGroupMemberInfoByGroupIdAndUserId(req.GroupInfo.GroupID, req.OpUserID)
	if err != nil {
		log.NewError(req.OperationID, "FindGroupMemberInfoByGroupIdAndUserId failed, ", err.Error(), req.GroupInfo.GroupID, req.OpUserID)
		return false

	}
	if groupUserInfo.AdministratorLevel == constant.OrdinaryMember {
		return true
	}
}

func (s *groupServer) SetGroupInfo(ctx context.Context, req *pbGroup.SetGroupInfoReq) (*pbGroup.CommonResp, error) {
	log.NewInfo(req.OperationID, "SetGroupInfo args ", req.String())
	if !hasAccess(req) {
		log.NewError(req.OperationID, "no access ")
		return &pbGroup.CommonResp{ErrCode: constant.ErrSetGroupInfo.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}, nil
	}

	group, err := im_mysql_model.FindGroupInfoByGroupId(req.GroupInfo.GroupID)
	if err != nil {
		log.NewError(req.OperationID, "FindGroupInfoByGroupId failed, ", err.Error(), req.GroupInfo.GroupID)
		return &pbGroup.CommonResp{ErrCode: constant.ErrSetGroupInfo.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}, nil
	}

	////bitwise operators: 1:groupName; 10:Notification  100:Introduction; 1000:FaceUrl
	var changedType int32
	if group.GroupName != req.GroupInfo.GroupName && req.GroupInfo.GroupName != "" {
		changedType = 1
	}
	if group.Notification != req.GroupInfo.Notification && req.GroupInfo.Notification != "" {
		changedType = changedType | (1 << 1)
	}
	if group.Introduction != req.GroupInfo.Introduction && req.GroupInfo.Introduction != "" {
		changedType = changedType | (1 << 2)
	}
	if group.FaceUrl != req.GroupInfo.FaceUrl && req.GroupInfo.FaceUrl != "" {
		changedType = changedType | (1 << 3)
	}
	//only administrators can set group information
	if err = im_mysql_model.SetGroupInfo(req.GroupInfo.GroupID, req.GroupInfo.GroupName, req.GroupInfo.Introduction, req.GroupInfo.Notification, req.GroupInfo.FaceUrl, ""); err != nil {
		return &pbGroup.CommonResp{ErrCode: constant.ErrSetGroupInfo.ErrCode, ErrMsg: constant.ErrSetGroupInfo.ErrMsg}, nil
	}

	if changedType != 0 {
		chat.GroupInfoChangedNotification(req)
	}

	return &pbGroup.CommonResp{}, nil
}

func (s *groupServer) TransferGroupOwner(_ context.Context, pb *pbGroup.TransferGroupOwnerReq) (*pbGroup.CommonResp, error) {
	log.Info("", "", "rpc TransferGroupOwner call start..., [pb: %s]", pb.String())

	reply, err := im_mysql_model.TransferGroupOwner(pb)
	if err != nil {
		log.Error("", "", "rpc TransferGroupOwner call..., im_mysql_model.TransferGroupOwner fail [pb: %s] [err: %s]", pb.String(), err.Error())
		return nil, err
	}
	log.Info("", "", "rpc TransferGroupOwner call..., im_mysql_model.TransferGroupOwner")

	return reply, nil
}
