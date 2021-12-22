package group

import (
	"Open_IM/internal/rpc/chat"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbGroup "Open_IM/pkg/proto/group"
	"Open_IM/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
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
	log.NewInfo(req.OperationID, "CreateGroup, args=%s", req.String())
	var (
		groupId string
	)
	//Parse token, to find current user information
	claims, err := token_verify.ParseToken(req.Token)
	if err != nil {
		log.NewError(req.OperationID, "ParseToken failed, ", err.Error(), req.String())
		return &pbGroup.CreateGroupResp{ErrCode: constant.ErrParseToken.ErrCode, ErrMsg: constant.ErrParseToken.ErrMsg}, nil
	}
	//Time stamp + MD5 to generate group chat id
	groupId = utils.Md5(strconv.FormatInt(time.Now().UnixNano(), 10))
	err = im_mysql_model.InsertIntoGroup(groupId, req.GroupName, req.Introduction, req.Notification, req.FaceUrl, req.Ext)
	if err != nil {
		log.NewError(req.OperationID, "InsertIntoGroup failed, ", err.Error(), req.String())
		return &pbGroup.CreateGroupResp{ErrCode: constant.ErrCreateGroup.ErrCode, ErrMsg: constant.ErrCreateGroup.ErrMsg}, nil
	}

	isManagerFlag := 0
	tokenUid := claims.UID

	if utils.IsContain(tokenUid, config.Config.Manager.AppManagerUid) {
		isManagerFlag = 1
	}

	us, err := im_mysql_model.FindUserByUID(claims.UID)
	if err != nil {
		log.Error("", req.OperationID, "find userInfo failed", err.Error())
		return &pbGroup.CreateGroupResp{ErrCode: constant.ErrCreateGroup.ErrCode, ErrMsg: constant.ErrCreateGroup.ErrMsg}, nil
	}

	if isManagerFlag == 0 {
		//Add the group owner to the group first, otherwise the group creation will fail
		err = im_mysql_model.InsertIntoGroupMember(groupId, claims.UID, us.Nickname, us.FaceUrl, constant.GroupOwner)
		if err != nil {
			log.Error("", req.OperationID, "create group chat failed,err=%s", err.Error())
			return &pbGroup.CreateGroupResp{ErrCode: constant.ErrCreateGroup.ErrCode, ErrMsg: constant.ErrCreateGroup.ErrMsg}, nil
		}

		err = db.DB.AddGroupMember(groupId, claims.UID)
		if err != nil {
			log.NewError(req.OperationID, "AddGroupMember failed ", err.Error(), groupId, claims.UID)
			return &pbGroup.CreateGroupResp{ErrCode: constant.ErrCreateGroup.ErrCode, ErrMsg: constant.ErrCreateGroup.ErrMsg}, nil
		}
	}

	//Binding group id and member id
	for _, user := range req.MemberList {
		us, err := im_mysql_model.FindUserByUID(user.Uid)
		if err != nil {
			log.NewError(req.OperationID, "FindUserByUID failed ", err.Error(), user.Uid)
			continue
		}
		err = im_mysql_model.InsertIntoGroupMember(groupId, user.Uid, us.Nickname, us.FaceUrl, user.SetRole)
		if err != nil {
			log.ErrorByArgs("InsertIntoGroupMember failed", user.Uid, groupId, err.Error())
		}
		err = db.DB.AddGroupMember(groupId, user.Uid)
		if err != nil {
			log.Error("", "", "add mongo group member failed, db.DB.AddGroupMember fail [err: %s]", err.Error())
		}
	}

	if isManagerFlag == 1 {

	}
	group, err := im_mysql_model.FindGroupInfoByGroupId(groupId)
	if err != nil {
		log.NewError(req.OperationID, "FindGroupInfoByGroupId failed ", err.Error(), groupId)
		return &pbGroup.CreateGroupResp{GroupID: groupId}, nil
	}
	memberList, err := im_mysql_model.FindGroupMemberListByGroupId(groupId)
	if err != nil {
		log.NewError(req.OperationID, "FindGroupMemberListByGroupId failed ", err.Error(), groupId)
		return &pbGroup.CreateGroupResp{GroupID: groupId}, nil
	}
	chat.GroupCreatedNotification(req.OperationID, us, group, memberList)
	log.NewInfo(req.OperationID, "GroupCreatedNotification, rpc CreateGroup success return ", groupId)

	return &pbGroup.CreateGroupResp{GroupID: groupId}, nil
}

func (s *groupServer) GetJoinedGroupList(ctx context.Context, req *pbGroup.GetJoinedGroupListReq) (*pbGroup.GetJoinedGroupListResp, error) {
	claims, err := token_verify.ParseToken(req.Token)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,parse token failed", err.Error())
		return &pbGroup.GetJoinedGroupListResp{ErrCode: constant.ErrParseToken.ErrCode, ErrorMsg: constant.ErrParseToken.ErrMsg}, nil
	}
	log.Info(claims.UID, req.OperationID, "recv req: ", req.String())

	joinedGroupList, err := imdb.GetJoinedGroupIdListByMemberId(claims.UID)
	if err != nil {
		log.Error(claims.UID, req.OperationID, "GetJoinedGroupIdListByMemberId failed, err: ", err.Error())
		return &pbGroup.GetJoinedGroupListResp{ErrCode: constant.ErrParam.ErrCode, ErrorMsg: constant.ErrParam.ErrMsg}, nil
	}

	var resp pbGroup.GetJoinedGroupListResp

	for _, v := range joinedGroupList {
		var groupNode pbGroup.GroupInfo
		num := imdb.GetGroupMemberNumByGroupId(v.GroupID)
		owner := imdb.GetGroupOwnerByGroupId(v.GroupID)
		group, err := imdb.FindGroupInfoByGroupId(v.GroupID)
		if num > 0 && owner != "" && err == nil {
			groupNode.GroupId = v.GroupID
			groupNode.FaceUrl = group.FaceUrl
			groupNode.CreateTime = uint64(group.CreateTime.Unix())
			groupNode.GroupName = group.GroupName
			groupNode.Introduction = group.Introduction
			groupNode.Notification = group.Notification
			groupNode.OwnerId = owner
			groupNode.MemberCount = uint32(int32(num))
			resp.GroupList = append(resp.GroupList, &groupNode)
		}
		log.Info(claims.UID, req.OperationID, "member num: ", num, "owner: ", owner)
	}
	resp.ErrCode = 0
	return &resp, nil
}

func (s *groupServer) InviteUserToGroup(ctx context.Context, req *pbGroup.InviteUserToGroupReq) (*pbGroup.InviteUserToGroupResp, error) {
	log.NewInfo(req.OperationID, "InviteUserToGroup args: ", req.String())
	if !imdb.IsExistGroupMember(req.GroupID, req.OpUserID) && !utils.IsContain(req.OpUserID, config.Config.Manager.AppManagerUid) {
		log.NewError(req.OperationID, "no permission InviteUserToGroup ", req.GroupID)
		return &pbGroup.InviteUserToGroupResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}, nil
	}

	groupInfoFromMysql, err := imdb.FindGroupInfoByGroupId(req.GroupID)
	if err != nil || groupInfoFromMysql == nil {
		log.NewError(req.OperationID, "FindGroupInfoByGroupId failed", req.GroupID)
		return &pbGroup.InviteUserToGroupResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}, nil
	}

	//
	//from User:  invite: applicant
	//to user:  invite: invited
	//to application
	var resp pbGroup.InviteUserToGroupResp
	opUser, err := imdb.FindUserByUID(req.OpUserID)
	if err != nil {
		log.NewError(req.OperationID, "FindUserByUID failed ", err.Error(), req.OpUserID)
	}
	var nicknameList string
	for _, v := range req.UidList {
		var resultNode pbGroup.Id2Result
		resultNode.UId = v
		resultNode.Result = 0
		toUserInfo, err := imdb.FindUserByUID(v)
		if err != nil {
			log.NewError(req.OperationID, "FindUserByUID failed ", err.Error())
			resultNode.Result = -1
			resp.Id2Result = append(resp.Id2Result, &resultNode)
			continue
		}

		if imdb.IsExistGroupMember(req.GroupID, v) {
			log.NewError(req.OperationID, "ExistGroupMember failed ", req.GroupID, v)
			resultNode.Result = -1
			resp.Id2Result = append(resp.Id2Result, &resultNode)
			continue
		}

		err = imdb.InsertGroupMember(req.GroupID, toUserInfo.UserID, toUserInfo.Nickname, toUserInfo.FaceUrl, 0)
		if err != nil {
			log.NewError(req.OperationID, "InsertGroupMember failed ", req.GroupID, toUserInfo.UserID, toUserInfo.Nickname, toUserInfo.FaceUrl)
			resultNode.Result = -1
			resp.Id2Result = append(resp.Id2Result, &resultNode)
			continue
		}
		member, err := imdb.GetMemberInfoById(req.GroupID, v)
		if groupInfoFromMysql != nil && opUser != nil && member != nil {
			chat.MemberInvitedNotification(req.OperationID, *groupInfoFromMysql, *opUser, *member)
		} else {
			log.NewError(req.OperationID, "args failed, nil ", groupInfoFromMysql, opUser, member)
		}

		err = db.DB.AddGroupMember(req.GroupID, toUserInfo.UserID)
		if err != nil {
			log.Error("", "", "add mongo group member failed, db.DB.AddGroupMember fail [err: %s]", err.Error())
		}
		nicknameList = nicknameList + toUserInfo.Nickname + " "
		resp.Id2Result = append(resp.Id2Result, &resultNode)
	}

	resp.ErrCode = 0
	return &resp, nil
}

type inviteUserToGroupReq struct {
	GroupID     string   `json:"groupID"`
	UidList     []string `json:"uidList"`
	Reason      string   `json:"reason"`
	OperationID string   `json:"operationID"`
}

func (c *inviteUserToGroupReq) ContentToString() string {
	data, _ := json.Marshal(c)
	dataString := string(data)
	return dataString
}

func (s *groupServer) GetGroupAllMember(ctx context.Context, req *pbGroup.GetGroupAllMemberReq) (*pbGroup.GetGroupAllMemberResp, error) {
	var resp pbGroup.GetGroupAllMemberResp
	resp.ErrCode = 0
	memberList, err := imdb.FindGroupMemberListByGroupId(req.GroupID)
	if err != nil {
		resp.ErrCode = constant.ErrDb.ErrCode
		resp.ErrMsg = err.Error()
		log.NewError(req.OperationID, "FindGroupMemberListByGroupId failed,", err.Error(), req.GroupID)
		return &resp, nil
	}

	for _, v := range memberList {
		var node pbGroup.GroupMemberFullInfo
		node.Role = v.AdministratorLevel
		node.NickName = v.NickName
		node.UserId = v.UserID
		node.FaceUrl = v.FaceUrl
		node.JoinTime = uint64(v.JoinTime.Unix())
		resp.MemberList = append(resp.MemberList, &node)
	}

	resp.ErrCode = 0
	return &resp, nil
}

func (s *groupServer) GetGroupMemberList(ctx context.Context, req *pbGroup.GetGroupMemberListReq) (*pbGroup.GetGroupMemberListResp, error) {
	claims, err := token_verify.ParseToken(req.Token)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,parse token failed", err.Error())
		return &pbGroup.GetGroupMemberListResp{ErrCode: constant.ErrParseToken.ErrCode, ErrMsg: constant.ErrParseToken.ErrMsg}, nil
	}
	//	log.Info(claims.UID, req.OperationID, "recv req: ", req.String())
	fmt.Println("req: ", req.GroupID)
	var resp pbGroup.GetGroupMemberListResp
	resp.ErrCode = 0
	memberList, err := imdb.GetGroupMemberByGroupId(req.GroupID, req.Filter, req.NextSeq, 30)
	if err != nil {
		resp.ErrCode = constant.ErrDb.ErrCode
		resp.ErrMsg = err.Error()
		log.Error(claims.UID, req.OperationID, "GetGroupMemberByGroupId failed, ", err.Error(), "params: ", req.GroupID, req.Filter, req.NextSeq)
		return &resp, nil
	}

	for _, v := range memberList {
		var node pbGroup.GroupMemberFullInfo
		node.Role = v.AdministratorLevel
		node.NickName = v.NickName
		node.UserId = v.UserID
		//	node.FaceUrl =
		node.JoinTime = uint64(v.JoinTime.Unix())
		resp.MemberList = append(resp.MemberList, &node)
	}
	//db operate  get db sorted by join time
	if int32(len(memberList)) < 30 {
		resp.NextSeq = 0
	} else {
		resp.NextSeq = req.NextSeq + int32(len(memberList))
	}

	resp.ErrCode = 0
	return &resp, nil
}

type groupMemberFullInfo struct {
	GroupId  string `json:"groupID"`
	UserId   string `json:"userId"`
	Role     int    `json:"role"`
	JoinTime uint64 `json:"joinTime"`
	NickName string `json:"nickName"`
	FaceUrl  string `json:"faceUrl"`
}

type kickGroupMemberApiReq struct {
	GroupID     string                `json:"groupID"`
	UidListInfo []groupMemberFullInfo `json:"uidListInfo"`
	Reason      string                `json:"reason"`
	OperationID string                `json:"operationID"`
}

func (c *kickGroupMemberApiReq) ContentToString() string {
	data, _ := json.Marshal(c)
	dataString := string(data)
	return dataString
}

func (s *groupServer) KickGroupMember(ctx context.Context, req *pbGroup.KickGroupMemberReq) (*pbGroup.KickGroupMemberResp, error) {
	log.NewInfo(req.OperationID, "KickGroupMember failed ", req.String())
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
			log.NewInfo(req.OperationID, "is group owner ", req.OpUserID, req.GroupID)
			break
		}
	}

	if flag != 1 {
		if utils.IsContain(req.OpUserID, config.Config.Manager.AppManagerUid) {
			flag = 1
			log.NewInfo(req.OperationID, "is app manager ", req.OpUserID, req.GroupID)
		}
	}

	if flag != 1 {
		log.NewError(req.OperationID, "failed, no access kick")
		return &pbGroup.KickGroupMemberResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}, nil
	}

	if len(req.UidListInfo) == 0 {
		log.NewError(req.OperationID, "failed, kick list 0")
		return &pbGroup.KickGroupMemberResp{ErrCode: constant.ErrParam.ErrCode, ErrMsg: constant.ErrParam.ErrMsg}, nil
	}

	//remove
	var resp pbGroup.KickGroupMemberResp
	for _, v := range req.UidListInfo {
		//owner cant kicked
		if v.UserId == req.OpUserID {
			log.NewError(req.OperationID, v.UserId, "failed, can't kick owner ", req.OpUserID)
			resp.Id2Result = append(resp.Id2Result, &pbGroup.Id2Result{UId: v.UserId, Result: -1})
			continue
		}
		err := imdb.RemoveGroupMember(req.GroupID, v.UserId)
		if err != nil {
			log.NewError(req.OperationID, "RemoveGroupMember failed ", err.Error(), req.GroupID, v.UserId)
			resp.Id2Result = append(resp.Id2Result, &pbGroup.Id2Result{UId: v.UserId, Result: -1})
		} else {
			resp.Id2Result = append(resp.Id2Result, &pbGroup.Id2Result{UId: v.UserId, Result: 0})
		}

		err = db.DB.DelGroupMember(req.GroupID, v.UserId)
		if err != nil {
			log.NewError(req.OperationID, "DelGroupMember failed ", err.Error(), req.GroupID, v.UserId)
		}
	}

	for _, v := range req.UidListInfo {
		chat.MemberKickedNotificationID(req.OperationID, req.GroupID, req.OpUserID, v.UserId, req.Reason)
	}

	resp.ErrCode = 0
	return &resp, nil
}

func (s *groupServer) GetGroupMembersInfo(ctx context.Context, req *pbGroup.GetGroupMembersInfoReq) (*pbGroup.GetGroupMembersInfoResp, error) {
	claims, err := token_verify.ParseToken(req.Token)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,parse token failed", err.Error())
		return &pbGroup.GetGroupMembersInfoResp{ErrCode: constant.ErrParseToken.ErrCode, ErrMsg: constant.ErrParseToken.ErrMsg}, nil
	}
	log.InfoByKv(claims.UID, req.OperationID, "param: ", req.MemberList)
	var resp pbGroup.GetGroupMembersInfoResp

	for _, v := range req.MemberList {
		var memberNode pbGroup.GroupMemberFullInfo
		memberInfo, err := imdb.GetMemberInfoById(req.GroupID, v)
		memberNode.UserId = v
		fmt.Println("id : ", memberNode.UserId)
		if err != nil {
			log.Error(claims.UID, req.OperationID, req.GroupID, v, "GetMemberInfoById failed, ", err.Error())
			//error occurs, only id is valid
			resp.MemberList = append(resp.MemberList, &memberNode)
			continue
		}
		user, err := imdb.FindUserByUID(v)
		if err == nil && user != nil {
			memberNode.FaceUrl = user.FaceUrl
			memberNode.JoinTime = uint64(memberInfo.JoinTime.Unix())
			memberNode.UserId = user.UserID
			memberNode.NickName = memberInfo.NickName
			memberNode.Role = memberInfo.AdministratorLevel
		}
		resp.MemberList = append(resp.MemberList, &memberNode)
	}
	resp.ErrCode = 0
	return &resp, nil
}

func (s *groupServer) GetGroupApplicationList(_ context.Context, pb *pbGroup.GetGroupApplicationListReq) (*pbGroup.GetGroupApplicationListResp, error) {
	log.Info("", "", "rpc GetGroupApplicationList call start..., [pb: %s]", pb.String())

	reply, err := im_mysql_model.GetGroupApplicationList(pb.OpUserID)
	if err != nil {
		return &pbGroup.GetGroupApplicationListResp{ErrCode: 701, ErrMsg: "GetGroupApplicationList failed"}, nil
	}
	log.Info("", "", "rpc GetGroupApplicationList call..., im_mysql_model.GetGroupApplicationList")

	return reply, nil
}

func (s *groupServer) GetGroupsInfo(ctx context.Context, req *pbGroup.GetGroupsInfoReq) (*pbGroup.GetGroupsInfoResp, error) {
	log.Info(req.Token, req.OperationID, "rpc get group info is server,args=%s", req.String())
	//Parse token, to find current user information
	claims, err := token_verify.ParseToken(req.Token)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,parse token failed", err.Error())
		return &pbGroup.GetGroupsInfoResp{ErrCode: constant.ErrParseToken.ErrCode, ErrMsg: constant.ErrParseToken.ErrMsg}, nil
	}
	log.Info("", req.OperationID, "args:", req.GroupIDList, claims.UID)
	groupsInfoList := make([]*pbGroup.GroupInfo, 0)
	for _, groupID := range req.GroupIDList {
		groupInfoFromMysql, err := im_mysql_model.FindGroupInfoByGroupId(groupID)
		if err != nil {
			log.Error(req.Token, req.OperationID, "find group info failed,err=%s", err.Error())
			continue
		}
		var groupInfo pbGroup.GroupInfo
		groupInfo.GroupId = groupID
		groupInfo.GroupName = groupInfoFromMysql.GroupName
		groupInfo.Introduction = groupInfoFromMysql.Introduction
		groupInfo.Notification = groupInfoFromMysql.Notification
		groupInfo.FaceUrl = groupInfoFromMysql.FaceUrl
		groupInfo.OwnerId = im_mysql_model.GetGroupOwnerByGroupId(groupID)
		groupInfo.MemberCount = uint32(im_mysql_model.GetGroupMemberNumByGroupId(groupID))
		groupInfo.CreateTime = uint64(groupInfoFromMysql.CreateTime.Unix())

		groupsInfoList = append(groupsInfoList, &groupInfo)
	}
	log.Info(req.Token, req.OperationID, "rpc get groupsInfo success return")
	return &pbGroup.GetGroupsInfoResp{Data: groupsInfoList}, nil
}

func (s *groupServer) GroupApplicationResponse(_ context.Context, pb *pbGroup.GroupApplicationResponseReq) (*pbGroup.GroupApplicationResponseResp, error) {
	log.NewInfo(pb.OperationID, "GroupApplicationResponse args: ", pb.String())
	reply, err := imdb.GroupApplicationResponse(pb)
	if err != nil {
		log.NewError(pb.OperationID, "GroupApplicationResponse failed ", err.Error(), pb.String())
		return &pbGroup.GroupApplicationResponseResp{ErrCode: 702, ErrMsg: err.Error()}, nil
	}

	if pb.HandleResult == 1 {
		if pb.ToUserID == "0" {
			err = db.DB.AddGroupMember(pb.GroupID, pb.FromUserID)
			if err != nil {
				log.Error("", "", "rpc GroupApplicationResponse call..., db.DB.AddGroupMember fail [pb: %s] [err: %s]", pb.String(), err.Error())
				return nil, err
			}
		} else {
			err = db.DB.AddGroupMember(pb.GroupID, pb.ToUserID)
			if err != nil {
				log.Error("", "", "rpc GroupApplicationResponse call..., db.DB.AddGroupMember fail [pb: %s] [err: %s]", pb.String(), err.Error())
				return nil, err
			}
		}
	}
	if pb.ToUserID == "0" {
		group, err := imdb.FindGroupInfoByGroupId(pb.GroupID)
		if err != nil {
			log.NewError(pb.OperationID, "FindGroupInfoByGroupId failed ", pb.GroupID)
			return reply, nil
		}
		member, err := imdb.FindGroupMemberInfoByGroupIdAndUserId(pb.GroupID, pb.OpUserID)
		if err != nil {
			log.NewError(pb.OperationID, "FindGroupMemberInfoByGroupIdAndUserId failed ", pb.GroupID, pb.OpUserID)
			return reply, nil
		}
		chat.ApplicationProcessedNotification(pb.OperationID, pb.FromUserID, *group, *member, pb.HandleResult, pb.HandledMsg)
		if pb.HandleResult == 1 {
			entrantUser, err := imdb.FindGroupMemberInfoByGroupIdAndUserId(pb.GroupID, pb.FromUserID)
			if err != nil {
				log.NewError(pb.OperationID, "FindGroupMemberInfoByGroupIdAndUserId failed ", err.Error(), pb.GroupID, pb.FromUserID)
				return reply, nil
			}
			chat.MemberEnterNotification(pb.OperationID, group, entrantUser)
		}
	} else {
		log.NewError(pb.OperationID, "args failed ", pb.String())
	}

	log.NewInfo(pb.OperationID, "rpc GroupApplicationResponse ok ", reply)
	return reply, nil
}

func (s *groupServer) JoinGroup(ctx context.Context, req *pbGroup.JoinGroupReq) (*pbGroup.CommonResp, error) {
	log.NewInfo(req.Token, req.OperationID, "JoinGroup args ", req.String())
	//Parse token, to find current user information
	claims, err := token_verify.ParseToken(req.Token)
	if err != nil {
		log.NewError(req.OperationID, "ParseToken failed", err.Error(), req.String())
		return &pbGroup.CommonResp{ErrCode: constant.ErrParseToken.ErrCode, ErrMsg: constant.ErrParseToken.ErrMsg}, nil
	}
	applicationUserInfo, err := im_mysql_model.FindUserByUID(claims.UID)
	if err != nil {
		log.NewError(req.OperationID, "FindUserByUID failed", err.Error(), claims.UID)
		return &pbGroup.CommonResp{ErrCode: constant.ErrSearchUserInfo.ErrCode, ErrMsg: constant.ErrSearchUserInfo.ErrMsg}, nil
	}

	_, err = im_mysql_model.FindGroupRequestUserInfoByGroupIDAndUid(req.GroupID, claims.UID)
	if err == nil {
		err = im_mysql_model.DelGroupRequest(req.GroupID, claims.UID, "0")
	}

	if err = im_mysql_model.InsertIntoGroupRequest(req.GroupID, claims.UID, "0", req.Message, applicationUserInfo.Nickname, applicationUserInfo.FaceUrl); err != nil {
		log.Error(req.Token, req.OperationID, "Insert into group request failed,er=%s", err.Error())
		return &pbGroup.CommonResp{ErrCode: constant.ErrJoinGroupApplication.ErrCode, ErrMsg: constant.ErrJoinGroupApplication.ErrMsg}, nil
	}

	memberList, err := im_mysql_model.FindGroupMemberListByGroupIdAndFilterInfo(req.GroupID, constant.GroupOwner)
	if len(memberList) == 0 {
		log.NewError(req.OperationID, "FindGroupMemberListByGroupIdAndFilterInfo failed ", req.GroupID, constant.GroupOwner, err)
		return &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}, nil
	}
	group, err := im_mysql_model.FindGroupInfoByGroupId(req.GroupID)
	if err != nil {
		log.NewError(req.OperationID, "FindGroupInfoByGroupId failed ", req.GroupID)
		return &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}, nil
	}
	chat.ReceiveJoinApplicationNotification(req.OperationID, memberList[0].UserID, applicationUserInfo, group)

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

	chat.MemberLeaveNotification(req.OperationID, req.GroupID, req.OpUserID)
	log.NewInfo(req.OperationID, "rpc quit group is success return")
	return &pbGroup.CommonResp{}, nil
}

func hasAccess(req *pbGroup.SetGroupInfoReq) bool {
	if utils.IsContain(req.OpUserID, config.Config.Manager.AppManagerUid) {
		return true
	}
	groupUserInfo, err := im_mysql_model.FindGroupMemberInfoByGroupIdAndUserId(req.GroupID, req.OpUserID)
	if err != nil {
		log.NewError(req.OperationID, "FindGroupMemberInfoByGroupIdAndUserId failed, ", err.Error(), req.GroupID, req.OpUserID)
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

	group, err := im_mysql_model.FindGroupInfoByGroupId(req.GroupID)
	if err != nil {
		log.NewError(req.OperationID, "FindGroupInfoByGroupId failed, ", err.Error(), req.GroupID)
		return &pbGroup.CommonResp{ErrCode: constant.ErrSetGroupInfo.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}, nil
	}

	////bitwise operators: 1:groupName; 10:Notification  100:Introduction; 1000:FaceUrl
	var changedType int32
	if group.GroupName != req.GroupName && req.GroupName != "" {
		changedType = 1
	}
	if group.Notification != req.Notification && req.Notification != "" {
		changedType = changedType | (1 << 1)
	}
	if group.Introduction != req.Introduction && req.Introduction != "" {
		changedType = changedType | (1 << 2)
	}
	if group.FaceUrl != req.FaceUrl && req.FaceUrl != "" {
		changedType = changedType | (1 << 3)
	}
	//only administrators can set group information
	if err = im_mysql_model.SetGroupInfo(req.GroupID, req.GroupName, req.Introduction, req.Notification, req.FaceUrl, ""); err != nil {
		return &pbGroup.CommonResp{ErrCode: constant.ErrSetGroupInfo.ErrCode, ErrMsg: constant.ErrSetGroupInfo.ErrMsg}, nil
	}

	if changedType != 0 {
		chat.GroupInfoChangedNotification(req.OperationID, changedType, req.GroupID, req.OpUserID)
	}

	return &pbGroup.CommonResp{}, nil
}

func (s *groupServer) TransferGroupOwner(_ context.Context, pb *pbGroup.TransferGroupOwnerReq) (*pbGroup.TransferGroupOwnerResp, error) {
	log.Info("", "", "rpc TransferGroupOwner call start..., [pb: %s]", pb.String())

	reply, err := im_mysql_model.TransferGroupOwner(pb)
	if err != nil {
		log.Error("", "", "rpc TransferGroupOwner call..., im_mysql_model.TransferGroupOwner fail [pb: %s] [err: %s]", pb.String(), err.Error())
		return nil, err
	}
	log.Info("", "", "rpc TransferGroupOwner call..., im_mysql_model.TransferGroupOwner")

	return reply, nil
}
