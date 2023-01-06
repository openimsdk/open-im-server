package group

import (
	"Open_IM/internal/rpc/fault_tolerant"
	chat "Open_IM/internal/rpc/msg"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	rocksCache "Open_IM/pkg/common/db/rocks_cache"
	"Open_IM/pkg/common/log"
	promePkg "Open_IM/pkg/common/prometheus"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/common/trace_log"
	cp "Open_IM/pkg/common/utils"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"github.com/OpenIMSDK/getcdv3"

	pbCache "Open_IM/pkg/proto/cache"
	pbConversation "Open_IM/pkg/proto/conversation"
	pbGroup "Open_IM/pkg/proto/group"
	pbUser "Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"fmt"
	"math/big"
	"net"
	"strconv"
	"strings"
	"time"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"
)

type groupServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func NewGroupServer(port int) *groupServer {
	log.NewPrivateLog(constant.LogFileName)
	return &groupServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImGroupName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (s *groupServer) Run() {
	log.NewInfo("", "group rpc start ")
	listenIP := ""
	if config.Config.ListenIP == "" {
		listenIP = "0.0.0.0"
	} else {
		listenIP = config.Config.ListenIP
	}
	address := listenIP + ":" + strconv.Itoa(s.rpcPort)
	//listener network
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic("listening err:" + err.Error() + s.rpcRegisterName)
	}
	log.NewInfo("", "listen network success, ", address, listener)
	defer listener.Close()
	//grpc server
	recvSize := 1024 * 1024 * constant.GroupRPCRecvSize
	sendSize := 1024 * 1024 * constant.GroupRPCSendSize
	var grpcOpts = []grpc.ServerOption{
		grpc.MaxRecvMsgSize(recvSize),
		grpc.MaxSendMsgSize(sendSize),
	}
	if config.Config.Prometheus.Enable {
		promePkg.NewGrpcRequestCounter()
		promePkg.NewGrpcRequestFailedCounter()
		promePkg.NewGrpcRequestSuccessCounter()
		grpcOpts = append(grpcOpts, []grpc.ServerOption{
			// grpc.UnaryInterceptor(promePkg.UnaryServerInterceptorProme),
			grpc.StreamInterceptor(grpcPrometheus.StreamServerInterceptor),
			grpc.UnaryInterceptor(grpcPrometheus.UnaryServerInterceptor),
		}...)
	}
	srv := grpc.NewServer(grpcOpts...)
	defer srv.GracefulStop()
	//Service registers with etcd
	pbGroup.RegisterGroupServer(srv, s)

	rpcRegisterIP := config.Config.RpcRegisterIP
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error("", "GetLocalIP failed ", err.Error())
		}
	}
	log.NewInfo("", "rpcRegisterIP", rpcRegisterIP)
	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName, 10, "")
	if err != nil {
		log.NewError("", "RegisterEtcd failed ", err.Error())
		panic(utils.Wrap(err, "register group module  rpc to etcd err"))

	}
	log.Info("", "RegisterEtcd ", s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName)
	err = srv.Serve(listener)
	if err != nil {
		log.NewError("", "Serve failed ", err.Error())
		return
	}
	log.NewInfo("", "group rpc success")
}

func (s *groupServer) CreateGroup(ctx context.Context, req *pbGroup.CreateGroupReq) (resp *pbGroup.CreateGroupResp, _ error) {
	resp = &pbGroup.CreateGroupResp{CommonResp: &open_im_sdk.CommonResp{}, GroupInfo: &open_im_sdk.GroupInfo{}}
	ctx = trace_log.NewRpcCtx(ctx, utils.GetSelfFuncName(), req.OperationID)
	defer func() {
		trace_log.SetContextInfo(ctx, utils.GetFuncName(1), nil, "req", req.String(), "resp", resp.String())
		trace_log.ShowLog(ctx)
	}()

	if err := token_verify.CheckAccessV2(ctx, req.OpUserID, req.OwnerUserID); err != nil {
		SetErrorForResp(err, &resp.CommonResp.ErrCode, &resp.CommonResp.ErrMsg)
		return
	}
	if err := callbackBeforeCreateGroup(ctx, req); err != nil {
		SetErrorForResp(err, &resp.CommonResp.ErrCode, &resp.CommonResp.ErrMsg)
		return
	}
	groupId := req.GroupInfo.GroupID
	if groupId == "" {
		groupId = utils.Md5(req.OperationID + strconv.FormatInt(time.Now().UnixNano(), 10))
		bi := big.NewInt(0)
		bi.SetString(groupId[0:8], 16)
		groupId = bi.String()
	}
	//to group
	groupInfo := imdb.Group{}
	utils.CopyStructFields(&groupInfo, req.GroupInfo)
	groupInfo.CreatorUserID = req.OpUserID
	groupInfo.GroupID = groupId
	groupInfo.CreateTime = time.Now()
	if groupInfo.NotificationUpdateTime.Unix() < 0 {
		groupInfo.NotificationUpdateTime = utils.UnixSecondToTime(0)
	}
	if err := (*imdb.Group)(nil).Create(ctx, []*imdb.Group{&groupInfo}); err != nil {
		SetErrorForResp(err, &resp.CommonResp.ErrCode, &resp.CommonResp.ErrMsg)
		return
	}
	groupMember := imdb.GroupMember{}
	us := &imdb.User{}
	if req.OwnerUserID != "" {
		var userIDList []string
		for _, v := range req.InitMemberList {
			userIDList = append(userIDList, v.UserID)
		}
		userIDList = append(userIDList, req.OwnerUserID)
		if err := s.DelGroupAndUserCache(req.OperationID, "", userIDList); err != nil {
			SetErr(ctx, "DelGroupAndUserCache", err, &resp.CommonResp.ErrCode, &resp.CommonResp.ErrMsg)
			return
		}
		var err error
		us, err = imdb.GetUserByUserID(req.OwnerUserID)
		if err != nil {
			SetErr(ctx, "GetUserByUserID", err, &resp.CommonResp.ErrCode, &resp.CommonResp.ErrMsg, "userID", req.OwnerUserID)
			return
		}
		//to group member
		groupMember = imdb.GroupMember{GroupID: groupId, RoleLevel: constant.GroupOwner, OperatorUserID: req.OpUserID, JoinSource: constant.JoinByInvitation, InviterUserID: req.OpUserID}
		utils.CopyStructFields(&groupMember, us)
		if err := CallbackBeforeMemberJoinGroup(ctx, req.OperationID, &groupMember, groupInfo.Ex); err != nil {
			SetErrorForResp(err, &resp.CommonResp.ErrCode, &resp.CommonResp.ErrMsg)
			return
		}
		if err := (*imdb.GroupMember)(nil).Create(ctx, []*imdb.GroupMember{&groupMember}); err != nil {
			SetErrorForResp(err, &resp.CommonResp.ErrCode, &resp.CommonResp.ErrMsg)
			return
		}
	}
	var okUserIDList []string
	if req.GroupInfo.GroupType != constant.SuperGroup {
		//to group member
		var groupMembers []*imdb.GroupMember
		for _, user := range req.InitMemberList {
			us, err := rocksCache.GetUserInfoFromCache(user.UserID)
			if err != nil {
				trace_log.SetContextInfo(ctx, "GetUserInfoFromCache", err, "userID", user.UserID)
				continue
			}
			if user.RoleLevel == constant.GroupOwner {
				trace_log.SetContextInfo(ctx, "GetUserInfoFromCache", nil, "userID", user.UserID, "msg", "only one owner, failed ")
				continue
			}
			groupMember.RoleLevel = user.RoleLevel
			groupMember.JoinSource = constant.JoinByInvitation
			groupMember.InviterUserID = req.OpUserID
			utils.CopyStructFields(&groupMember, us)
			if err := CallbackBeforeMemberJoinGroup(ctx, req.OperationID, &groupMember, groupInfo.Ex); err != nil {
				SetErrorForResp(err, &resp.CommonResp.ErrCode, &resp.CommonResp.ErrMsg)
				return
			}
			groupMembers = append(groupMembers, &groupMember)
			okUserIDList = append(okUserIDList, user.UserID)
		}
		if err := (*imdb.GroupMember)(nil).Create(ctx, groupMembers); err != nil {
			SetErrorForResp(err, &resp.CommonResp.ErrCode, &resp.CommonResp.ErrMsg)
			return
		}
		group, err := rocksCache.GetGroupInfoFromCache(ctx, groupId)
		if err != nil {
			SetErr(ctx, "GetGroupInfoFromCache", err, &resp.CommonResp.ErrCode, &resp.CommonResp.ErrMsg, "groupID", groupId)
			return
		}
		utils.CopyStructFields(resp.GroupInfo, group)
		memberCount, err := rocksCache.GetGroupMemberNumFromCache(groupId)
		if err != nil {
			SetErr(ctx, "GetGroupMemberNumFromCache", err, &resp.CommonResp.ErrCode, &resp.CommonResp.ErrMsg, "groupID", groupId)
			return
		}
		resp.GroupInfo.MemberCount = uint32(memberCount)
		if req.OwnerUserID != "" {
			resp.GroupInfo.OwnerUserID = req.OwnerUserID
			okUserIDList = append(okUserIDList, req.OwnerUserID)
		}
		// superGroup stored in mongodb
	} else {
		for _, v := range req.InitMemberList {
			okUserIDList = append(okUserIDList, v.UserID)
		}
		if err := db.DB.CreateSuperGroup(groupId, okUserIDList, len(okUserIDList)); err != nil {
			SetErr(ctx, "CreateSuperGroup", err, &resp.CommonResp.ErrCode, &resp.CommonResp.ErrMsg, "groupID", groupId, "userIDList", okUserIDList)
			return
		}
	}

	if len(okUserIDList) != 0 {
		if req.GroupInfo.GroupType != constant.SuperGroup {
			chat.GroupCreatedNotification(req.OperationID, req.OpUserID, groupId, okUserIDList)
		} else {
			for _, userID := range okUserIDList {
				if err := rocksCache.DelJoinedSuperGroupIDListFromCache(userID); err != nil {
					trace_log.SetContextInfo(ctx, "DelJoinedSuperGroupIDListFromCache", err, "userID", userID)
					//log.NewWarn(req.OperationID, utils.GetSelfFuncName(), userID, err.Error())
				}
			}
			go func() {
				for _, v := range okUserIDList {
					chat.SuperGroupNotification(req.OperationID, v, v)
				}
			}()
		}
		return
	} else {
		//log.NewInfo(req.OperationID, "rpc CreateGroup return ", resp.String())
		return
	}
}

func (s *groupServer) GetJoinedGroupList(ctx context.Context, req *pbGroup.GetJoinedGroupListReq) (resp *pbGroup.GetJoinedGroupListResp,_ error) {
	resp = &pbGroup.GetJoinedGroupListResp{CommonResp: &open_im_sdk.CommonResp{}}
	ctx = trace_log.NewRpcCtx(ctx, utils.GetSelfFuncName(), req.OperationID)
	trace_log.SetContextInfo(ctx, utils.GetSelfFuncName(), nil, "req", req, "resp", resp)
	defer trace_log.ShowLog(ctx)
	if err := token_verify.CheckAccessV2(ctx, req.OpUserID, req.FromUserID);err != nil {
		SetErrorForResp(err, &resp.CommonResp.ErrCode, &resp.CommonResp.ErrMsg)
		return
	}
	joinedGroupList, err := rocksCache.GetJoinedGroupIDListFromCache(req.FromUserID)
	if err != nil {
		SetErr(ctx, "GetJoinedGroupIDListFromCache", err, &resp.CommonResp.ErrCode, &resp.CommonResp.ErrMsg, "userID", req.FromUserID)
		return
	}
	log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "joinedGroupList: ", joinedGroupList)
	for _, v := range joinedGroupList {
		var groupNode open_im_sdk.GroupInfo
		num, err := rocksCache.GetGroupMemberNumFromCache(v)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), v)
			continue
		}
		owner, err2 := imdb.GetGroupOwnerInfoByGroupID(v)
		if err2 != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err2.Error(), v)
			continue
		}
		group, err := rocksCache.GetGroupInfoFromCache(v)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), v)
			continue
		}
		if group.GroupType == constant.SuperGroup {
			continue
		}
		if group.Status == constant.GroupStatusDismissed {
			log.NewError(req.OperationID, "constant.GroupStatusDismissed ", group)
			continue
		}
		utils.CopyStructFields(&groupNode, group)
		groupNode.CreateTime = uint32(group.CreateTime.Unix())
		groupNode.NotificationUpdateTime = uint32(group.NotificationUpdateTime.Unix())
		if group.NotificationUpdateTime.Unix() < 0 {
			groupNode.NotificationUpdateTime = 0
		}

		groupNode.MemberCount = uint32(num)
		groupNode.OwnerUserID = owner.UserID
		resp.GroupList = append(resp.GroupList, &groupNode)
		log.NewDebug(req.OperationID, "joinedGroup ", groupNode)
	}
	log.NewInfo(req.OperationID, "GetJoinedGroupList rpc return ", resp.String())
	return &resp, nil
}

func (s *groupServer) InviteUserToGroup(ctx context.Context, req *pbGroup.InviteUserToGroupReq) (*pbGroup.InviteUserToGroupResp, error) {
	log.NewInfo(req.OperationID, "InviteUserToGroup args ", req.String())
	if !imdb.IsExistGroupMember(req.GroupID, req.OpUserID) && !token_verify.IsManagerUserID(req.OpUserID) {
		log.NewError(req.OperationID, "no permission InviteUserToGroup ", req.GroupID, req.OpUserID)
		return &pbGroup.InviteUserToGroupResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}, nil
	}

	groupInfo, err := imdb.GetGroupInfoByGroupID(req.GroupID)
	if err != nil {
		log.NewError(req.OperationID, "GetGroupInfoByGroupID failed ", req.GroupID, err)
		return &pbGroup.InviteUserToGroupResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}, nil
	}
	if groupInfo.Status == constant.GroupStatusDismissed {
		errMsg := " group status is dismissed "
		return &pbGroup.InviteUserToGroupResp{ErrCode: constant.ErrStatus.ErrCode, ErrMsg: errMsg}, nil
	}
	var resp pbGroup.InviteUserToGroupResp
	if groupInfo.NeedVerification == constant.AllNeedVerification &&
		!imdb.IsGroupOwnerAdmin(req.GroupID, req.OpUserID) && !token_verify.IsManagerUserID(req.OpUserID) {
		var resp pbGroup.InviteUserToGroupResp
		joinReq := pbGroup.JoinGroupReq{}
		for _, v := range req.InvitedUserIDList {
			var groupRequest imdb.GroupRequest
			groupRequest.UserID = v
			groupRequest.GroupID = req.GroupID
			groupRequest.JoinSource = constant.JoinByInvitation
			groupRequest.InviterUserID = req.OpUserID
			err = imdb.InsertIntoGroupRequest(groupRequest)
			if err != nil {
				var resultNode pbGroup.Id2Result
				resultNode.Result = -1
				resultNode.UserID = v
				resp.Id2ResultList = append(resp.Id2ResultList, &resultNode)

				continue
				// log.NewError(req.OperationID, "InsertIntoGroupRequest failed ", err.Error(), groupRequest)
				//	return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
			} else {
				var resultNode pbGroup.Id2Result
				resultNode.Result = 0
				resultNode.UserID = v
				resp.Id2ResultList = append(resp.Id2ResultList, &resultNode)
				joinReq.GroupID = req.GroupID
				joinReq.OperationID = req.OperationID
				joinReq.OpUserID = v
				resp.Id2ResultList = append(resp.Id2ResultList, &resultNode)
				chat.JoinGroupApplicationNotification(&joinReq)
			}
		}
		log.NewInfo(req.OperationID, "InviteUserToGroup rpc return ", resp)
		return &resp, nil
	}
	if err := s.DelGroupAndUserCache(req.OperationID, req.GroupID, req.InvitedUserIDList); err != nil {
		log.NewError(req.OperationID, "DelGroupAndUserCache failed", err.Error())
		return &pbGroup.InviteUserToGroupResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}, nil
	}
	//from User:  invite: applicant
	//to user:  invite: invited
	var okUserIDList []string
	if groupInfo.GroupType != constant.SuperGroup {
		for _, v := range req.InvitedUserIDList {
			var resultNode pbGroup.Id2Result
			resultNode.UserID = v
			resultNode.Result = 0
			toUserInfo, err := imdb.GetUserByUserID(v)
			if err != nil {
				log.NewError(req.OperationID, "GetUserByUserID failed ", err.Error(), v)
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
			var toInsertInfo imdb.GroupMember
			utils.CopyStructFields(&toInsertInfo, toUserInfo)
			toInsertInfo.GroupID = req.GroupID
			toInsertInfo.RoleLevel = constant.GroupOrdinaryUsers
			toInsertInfo.OperatorUserID = req.OpUserID
			toInsertInfo.InviterUserID = req.OpUserID
			toInsertInfo.JoinSource = constant.JoinByInvitation
			callbackResp := CallbackBeforeMemberJoinGroup(req.OperationID, &toInsertInfo, groupInfo.Ex)
			if callbackResp.ErrCode != 0 {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "callbackBeforeSendSingleMsg resp: ", callbackResp)
			}
			if callbackResp.ActionCode != constant.ActionAllow {
				log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "callbackBeforeSendSingleMsg result", "end rpc and return", callbackResp)
				continue
			}

			err = imdb.InsertIntoGroupMember(toInsertInfo)
			if err != nil {
				log.NewError(req.OperationID, "InsertIntoGroupMember failed ", req.GroupID, toUserInfo.UserID, toUserInfo.Nickname, toUserInfo.FaceURL)
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
	} else {
		okUserIDList = req.InvitedUserIDList
		if err := db.DB.AddUserToSuperGroup(req.GroupID, req.InvitedUserIDList); err != nil {
			log.NewError(req.OperationID, "AddUserToSuperGroup failed ", req.GroupID, err)
			return &pbGroup.InviteUserToGroupResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}, nil
		}
	}

	// set conversations
	var haveConUserID []string
	var sessionType int
	if groupInfo.GroupType == constant.NormalGroup {
		sessionType = constant.GroupChatType
	} else {
		sessionType = constant.SuperGroupChatType
	}
	conversations, err := imdb.GetConversationsByConversationIDMultipleOwner(okUserIDList, utils.GetConversationIDBySessionType(req.GroupID, sessionType))
	if err != nil {
		log.NewError(req.OperationID, "GetConversationsByConversationIDMultipleOwner failed ", err.Error(), req.GroupID, sessionType)
	}
	for _, v := range conversations {
		haveConUserID = append(haveConUserID, v.OwnerUserID)
	}
	var reqPb pbUser.SetConversationReq
	var c pbConversation.Conversation
	for _, v := range conversations {
		reqPb.OperationID = req.OperationID
		c.OwnerUserID = v.OwnerUserID
		c.ConversationID = utils.GetConversationIDBySessionType(req.GroupID, sessionType)
		c.RecvMsgOpt = v.RecvMsgOpt
		c.ConversationType = int32(sessionType)
		c.GroupID = req.GroupID
		c.IsPinned = v.IsPinned
		c.AttachedInfo = v.AttachedInfo
		c.IsPrivateChat = v.IsPrivateChat
		c.GroupAtType = v.GroupAtType
		c.IsNotInGroup = false
		c.Ex = v.Ex
		reqPb.Conversation = &c
		etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
		if etcdConn == nil {
			errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
			log.NewError(req.OperationID, errMsg)
			return &pbGroup.InviteUserToGroupResp{ErrCode: constant.ErrInternal.ErrCode, ErrMsg: errMsg}, nil
		}
		client := pbUser.NewUserClient(etcdConn)
		respPb, err := client.SetConversation(context.Background(), &reqPb)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "SetConversation rpc failed, ", reqPb.String(), err.Error(), v.OwnerUserID)
		} else {
			log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "SetConversation success", respPb.String(), v.OwnerUserID)
		}
	}
	for _, v := range utils.DifferenceString(haveConUserID, okUserIDList) {
		reqPb.OperationID = req.OperationID
		c.OwnerUserID = v
		c.ConversationID = utils.GetConversationIDBySessionType(req.GroupID, sessionType)
		c.ConversationType = int32(sessionType)
		c.GroupID = req.GroupID
		c.IsNotInGroup = false
		c.UpdateUnreadCountTime = utils.GetCurrentTimestampByMill()
		reqPb.Conversation = &c
		etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
		if etcdConn == nil {
			errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
			log.NewError(req.OperationID, errMsg)
			return &pbGroup.InviteUserToGroupResp{ErrCode: constant.ErrInternal.ErrCode, ErrMsg: errMsg}, nil
		}
		client := pbUser.NewUserClient(etcdConn)
		respPb, err := client.SetConversation(context.Background(), &reqPb)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "SetConversation rpc failed, ", reqPb.String(), err.Error(), v)
		} else {
			log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "SetConversation success", respPb.String(), v)
		}
	}

	if groupInfo.GroupType != constant.SuperGroup {
		chat.MemberInvitedNotification(req.OperationID, req.GroupID, req.OpUserID, req.Reason, okUserIDList)
	} else {
		for _, v := range req.InvitedUserIDList {
			if err := rocksCache.DelJoinedSuperGroupIDListFromCache(v); err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error())
			}
		}
		for _, v := range req.InvitedUserIDList {
			chat.SuperGroupNotification(req.OperationID, v, v)
		}
	}

	log.NewInfo(req.OperationID, "InviteUserToGroup rpc return ", resp)
	return &resp, nil
}

func (s *groupServer) GetGroupAllMember(ctx context.Context, req *pbGroup.GetGroupAllMemberReq) (*pbGroup.GetGroupAllMemberResp, error) {
	log.NewInfo(req.OperationID, "GetGroupAllMember, args ", req.String())
	var resp pbGroup.GetGroupAllMemberResp
	groupInfo, err := rocksCache.GetGroupInfoFromCache(req.GroupID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.GroupID)
		resp.ErrCode = constant.ErrDB.ErrCode
		resp.ErrMsg = constant.ErrDB.ErrMsg
		return &resp, nil
	}
	if groupInfo.GroupType != constant.SuperGroup {
		memberList, err := rocksCache.GetGroupMembersInfoFromCache(req.Count, req.Offset, req.GroupID)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.GroupID)
			resp.ErrCode = constant.ErrDB.ErrCode
			resp.ErrMsg = constant.ErrDB.ErrMsg
			return &resp, nil
		}
		for _, v := range memberList {
			var node open_im_sdk.GroupMemberFullInfo
			cp.GroupMemberDBCopyOpenIM(&node, v)
			resp.MemberList = append(resp.MemberList, &node)
		}
	}
	log.NewInfo(req.OperationID, "GetGroupAllMember rpc return ", len(resp.MemberList))
	return &resp, nil
}

func (s *groupServer) GetGroupMemberList(ctx context.Context, req *pbGroup.GetGroupMemberListReq) (*pbGroup.GetGroupMemberListResp, error) {
	log.NewInfo(req.OperationID, "GetGroupMemberList args ", req.String())
	var resp pbGroup.GetGroupMemberListResp
	memberList, err := imdb.GetGroupMemberByGroupID(req.GroupID, req.Filter, req.NextSeq, 30)
	if err != nil {
		resp.ErrCode = constant.ErrDB.ErrCode
		resp.ErrMsg = constant.ErrDB.ErrMsg
		log.NewError(req.OperationID, "GetGroupMemberByGroupId failed,", req.GroupID, req.Filter, req.NextSeq, 30)
		return &resp, nil
	}

	for _, v := range memberList {
		var node open_im_sdk.GroupMemberFullInfo
		utils.CopyStructFields(&node, &v)
		resp.MemberList = append(resp.MemberList, &node)
	}
	//db operate  get db sorted by join time
	if int32(len(memberList)) < 30 {
		resp.NextSeq = 0
	} else {
		resp.NextSeq = req.NextSeq + int32(len(memberList))
	}
	resp.ErrCode = 0
	log.NewInfo(req.OperationID, "GetGroupMemberList rpc return ", resp.String())
	return &resp, nil
}

func (s *groupServer) getGroupUserLevel(groupID, userID string) (int, error) {
	opFlag := 0
	if !token_verify.IsManagerUserID(userID) {
		opInfo, err := imdb.GetGroupMemberInfoByGroupIDAndUserID(groupID, userID)
		if err != nil {
			return opFlag, utils.Wrap(err, "")
		}
		if opInfo.RoleLevel == constant.GroupOrdinaryUsers {
			opFlag = 0
		} else if opInfo.RoleLevel == constant.GroupOwner {
			opFlag = 2 //owner
		} else {
			opFlag = 3 //admin
		}
	} else {
		opFlag = 1 //app manager
	}
	return opFlag, nil
}

func (s *groupServer) KickGroupMember(ctx context.Context, req *pbGroup.KickGroupMemberReq) (*pbGroup.KickGroupMemberResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc args ", req.String())
	groupInfo, err := rocksCache.GetGroupInfoFromCache(req.GroupID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetGroupInfoByGroupID", req.GroupID, err.Error())
		return &pbGroup.KickGroupMemberResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}, nil
	}
	var okUserIDList []string
	var resp pbGroup.KickGroupMemberResp
	if groupInfo.GroupType != constant.SuperGroup {
		opFlag := 0
		if !token_verify.IsManagerUserID(req.OpUserID) {
			opInfo, err := rocksCache.GetGroupMemberInfoFromCache(req.GroupID, req.OpUserID)
			if err != nil {
				errMsg := req.OperationID + " GetGroupMemberInfoByGroupIDAndUserID  failed " + err.Error() + req.GroupID + req.OpUserID
				log.Error(req.OperationID, errMsg)
				return &pbGroup.KickGroupMemberResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
			}
			if opInfo.RoleLevel == constant.GroupOrdinaryUsers {
				errMsg := req.OperationID + " opInfo.RoleLevel == constant.GroupOrdinaryUsers " + opInfo.UserID + opInfo.GroupID
				log.Error(req.OperationID, errMsg)
				return &pbGroup.KickGroupMemberResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}, nil
			} else if opInfo.RoleLevel == constant.GroupOwner {
				opFlag = 2 //owner
			} else {
				opFlag = 3 //admin
			}
		} else {
			opFlag = 1 //app manager
		}

		//op is group owner?
		if len(req.KickedUserIDList) == 0 {
			log.NewError(req.OperationID, "failed, kick list 0")
			return &pbGroup.KickGroupMemberResp{ErrCode: constant.ErrArgs.ErrCode, ErrMsg: constant.ErrArgs.ErrMsg}, nil
		}
		if err := s.DelGroupAndUserCache(req.OperationID, req.GroupID, req.KickedUserIDList); err != nil {
			log.NewError(req.OperationID, "DelGroupAndUserCache failed", err.Error())
			return &pbGroup.KickGroupMemberResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}, nil
		}
		//remove
		for _, v := range req.KickedUserIDList {
			kickedInfo, err := rocksCache.GetGroupMemberInfoFromCache(req.GroupID, v)
			if err != nil {
				log.NewError(req.OperationID, " GetGroupMemberInfoByGroupIDAndUserID failed ", req.GroupID, v, err.Error())
				resp.Id2ResultList = append(resp.Id2ResultList, &pbGroup.Id2Result{UserID: v, Result: -1})
				continue
			}

			if kickedInfo.RoleLevel == constant.GroupAdmin && opFlag == 3 {
				log.Error(req.OperationID, "is constant.GroupAdmin, can't kicked ", v)
				resp.Id2ResultList = append(resp.Id2ResultList, &pbGroup.Id2Result{UserID: v, Result: -1})
				continue
			}
			if kickedInfo.RoleLevel == constant.GroupOwner && opFlag != 1 {
				log.NewDebug(req.OperationID, "is constant.GroupOwner, can't kicked ", v)
				resp.Id2ResultList = append(resp.Id2ResultList, &pbGroup.Id2Result{UserID: v, Result: -1})
				continue
			}

			err = imdb.DeleteGroupMemberByGroupIDAndUserID(req.GroupID, v)
			if err != nil {
				log.NewError(req.OperationID, "RemoveGroupMember failed ", err.Error(), req.GroupID, v)
				resp.Id2ResultList = append(resp.Id2ResultList, &pbGroup.Id2Result{UserID: v, Result: -1})
			} else {
				log.NewDebug(req.OperationID, "kicked ", v)
				resp.Id2ResultList = append(resp.Id2ResultList, &pbGroup.Id2Result{UserID: v, Result: 0})
				okUserIDList = append(okUserIDList, v)
			}
		}
		var reqPb pbUser.SetConversationReq
		var c pbConversation.Conversation
		for _, v := range okUserIDList {
			reqPb.OperationID = req.OperationID
			c.OwnerUserID = v
			c.ConversationID = utils.GetConversationIDBySessionType(req.GroupID, constant.GroupChatType)
			c.ConversationType = constant.GroupChatType
			c.GroupID = req.GroupID
			c.IsNotInGroup = true
			reqPb.Conversation = &c
			etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
			if etcdConn == nil {
				errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
				log.NewError(req.OperationID, errMsg)
				resp.ErrCode = constant.ErrInternal.ErrCode
				resp.ErrMsg = errMsg
				return &resp, nil
			}
			client := pbUser.NewUserClient(etcdConn)
			respPb, err := client.SetConversation(context.Background(), &reqPb)
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "SetConversation rpc failed, ", reqPb.String(), err.Error(), v)
			} else {
				log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "SetConversation success", respPb.String(), v)
			}
		}
	} else {
		okUserIDList = req.KickedUserIDList
		if err := db.DB.RemoverUserFromSuperGroup(req.GroupID, okUserIDList); err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), req.GroupID, req.KickedUserIDList, err.Error())
			resp.ErrCode = constant.ErrDB.ErrCode
			resp.ErrMsg = constant.ErrDB.ErrMsg
			return &resp, nil
		}
	}

	if groupInfo.GroupType != constant.SuperGroup {
		for _, userID := range okUserIDList {
			if err := rocksCache.DelGroupMemberInfoFromCache(req.GroupID, userID); err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.GroupID)
			}
		}
		chat.MemberKickedNotification(req, okUserIDList)
	} else {
		for _, userID := range okUserIDList {
			err = rocksCache.DelJoinedSuperGroupIDListFromCache(userID)
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), userID)
			}
		}
		go func() {
			for _, v := range req.KickedUserIDList {
				chat.SuperGroupNotification(req.OperationID, v, v)
			}
		}()

	}

	log.NewInfo(req.OperationID, "GetGroupMemberList rpc return ", resp.String())
	return &resp, nil
}

func (s *groupServer) GetGroupMembersInfo(ctx context.Context, req *pbGroup.GetGroupMembersInfoReq) (resp *pbGroup.GetGroupMembersInfoResp, err error) {
	defer func() {
		trace_log.SetContextInfo(ctx, utils.GetSelfFuncName(), err, "rpc req", req.String(), "rpc resp", resp.String())
		trace_log.ShowLog(ctx)
	}()
	ctx = trace_log.NewRpcCtx(ctx, utils.GetSelfFuncName(), req.OperationID)

	resp = &pbGroup.GetGroupMembersInfoResp{CommonResp: &open_im_sdk.CommonResp{}}
	resp.MemberList = []*open_im_sdk.GroupMemberFullInfo{}
	for _, userID := range req.MemberList {
		groupMember, err := rocksCache.GetGroupMemberInfoFromCache(req.GroupID, userID)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), req.GroupID, userID, err.Error())
			var errResult error

			continue
		}
		var memberNode open_im_sdk.GroupMemberFullInfo
		utils.CopyStructFields(&memberNode, groupMember)
		memberNode.JoinTime = int32(groupMember.JoinTime.Unix())
		resp.MemberList = append(resp.MemberList, &memberNode)
	}
	if
	SetErr(ctx, "GetGroupMemberInfoFromCache", err, &resp.CommonResp.ErrCode, &resp.CommonResp.ErrMsg, "groupID", req.GroupID)
	return resp, nil
}

func FillGroupInfoByGroupID(operationID, groupID string, groupInfo *open_im_sdk.GroupInfo) error {
	group, err := imdb.TakeGroupInfoByGroupID(groupID)
	if err != nil {
		log.Error(operationID, "TakeGroupInfoByGroupID failed ", err.Error(), groupID)
		return utils.Wrap(err, "")
	}
	if group.Status == constant.GroupStatusDismissed {
		log.Debug(operationID, " group constant.GroupStatusDismissed ", group.GroupID)
		return utils.Wrap(constant.ErrDismissedAlready, "")
	}
	return utils.Wrap(cp.GroupDBCopyOpenIM(groupInfo, group), "")
}

func FillPublicUserInfoByUserID(operationID, userID string, userInfo *open_im_sdk.PublicUserInfo) error {
	user, err := imdb.TakeUserByUserID(userID)
	if err != nil {
		log.Error(operationID, "TakeUserByUserID failed ", err.Error(), userID)
		return utils.Wrap(err, "")
	}
	cp.UserDBCopyOpenIMPublicUser(userInfo, user)
	return nil
}

func SetErr(ctx context.Context, funcName string, err error, errCode *int32, errMsg *string, args ...interface{}) {
	errInfo := constant.ToAPIErrWithErr(err)
	*errCode = errInfo.ErrCode
	*errMsg = errInfo.WrapErrMsg
	trace_log.SetContextInfo(ctx, funcName, err, args)
}

func SetErrCodeMsg(ctx context.Context, funcName string, errCodeFrom int32, errMsgFrom string, errCode *int32, errMsg *string, args ...interface{}) {
	*errCode = errCodeFrom
	*errMsg = errMsgFrom
	trace_log.SetContextInfo(ctx, funcName, constant.ToAPIErrWithErrCode(errCodeFrom), args)
}

func SetErrorForResp(err error, commonResp *open_im_sdk.CommonResp) {
	errInfo := constant.ToAPIErrWithErr(err)
	commonResp.ErrCode = errInfo.ErrCode
	commonResp.ErrMsg = errInfo.ErrMsg
	commonResp.DetailErrMsg = err.Error()
}

func (s *groupServer) GetGroupApplicationList(ctx context.Context, req *pbGroup.GetGroupApplicationListReq) (resp *pbGroup.GetGroupApplicationListResp, err error) {
	ctx = trace_log.NewRpcCtx(ctx, utils.GetSelfFuncName(), req.OperationID)
	trace_log.SetContextInfo(ctx, utils.GetSelfFuncName(), err, "rpc req ", req.String(), "rpc resp ", resp.String())
	trace_log.ShowLog(ctx)

	resp = &pbGroup.GetGroupApplicationListResp{CommonResp: &open_im_sdk.CommonResp{}}
	reply, err := imdb.GetRecvGroupApplicationList(req.FromUserID)
	if err != nil {
		SetErrorForResp(err, &resp.CommonResp.ErrCode, &resp.CommonResp.ErrMsg)
		return resp, nil
	}
	var errResult error
	trace_log.SetContextInfo(ctx, "GetRecvGroupApplicationList", nil, " FromUserID: ", req.FromUserID, "GroupApplicationList: ", reply)
	for _, v := range reply {
		node := open_im_sdk.GroupRequest{UserInfo: &open_im_sdk.PublicUserInfo{}, GroupInfo: &open_im_sdk.GroupInfo{}}
		err := FillGroupInfoByGroupID(req.OperationID, v.GroupID, node.GroupInfo)
		if err != nil {
			if !errors.Is(errors.Unwrap(err), constant.ErrDismissedAlready) {
				errResult = err
			}
			continue
		}
		trace_log.SetContextInfo(ctx, "FillGroupInfoByGroupID ", nil, " groupID: ", v.GroupID, " groupInfo: ", node.GroupInfo)
		err = FillPublicUserInfoByUserID(req.OperationID, v.UserID, node.UserInfo)
		if err != nil {
			errResult = err
			continue
		}
		cp.GroupRequestDBCopyOpenIM(&node, &v)
		resp.GroupRequestList = append(resp.GroupRequestList, &node)
	}
	if errResult != nil && len(resp.GroupRequestList) == 0 {
		SetErr(ctx, "", errResult, &resp.CommonResp.ErrCode, &resp.CommonResp.ErrMsg)
		return resp, nil
	}
	trace_log.SetRpcRespInfo(ctx, utils.GetSelfFuncName(), resp.String())
	return resp, nil
}

func (s *groupServer) GetGroupsInfo(ctx context.Context, req *pbGroup.GetGroupsInfoReq) (resp *pbGroup.GetGroupsInfoResp, err error) {
	defer func() {
		trace_log.SetContextInfo(ctx, utils.GetSelfFuncName(), err, "rpc req ", req.String(), "rpc resp ", resp.String())
		trace_log.ShowLog(ctx)
	}()
	ctx = trace_log.NewRpcCtx(ctx, utils.GetSelfFuncName(), req.OperationID)

	resp = &pbGroup.GetGroupsInfoResp{CommonResp: &open_im_sdk.CommonResp{}}
	groupsInfoList := make([]*open_im_sdk.GroupInfo, 0)
	for _, groupID := range req.GroupIDList {
		groupInfoFromRedis, err := rocksCache.GetGroupInfoFromCache(ctx, groupID)
		if err != nil {
			SetErrorForResp(err, &resp.CommonResp.ErrCode, &resp.CommonResp.ErrMsg)
			continue
		}
		var groupInfo open_im_sdk.GroupInfo
		cp.GroupDBCopyOpenIM(&groupInfo, groupInfoFromRedis)
		groupInfo.NeedVerification = groupInfoFromRedis.NeedVerification
		groupsInfoList = append(groupsInfoList, &groupInfo)
	}
	resp.GroupInfoList = groupsInfoList
	return resp, nil
}

func CheckPermission(ctx context.Context, groupID string, userID string) (err error) {
	defer func() {
		trace_log.SetContextInfo(ctx, utils.GetSelfFuncName(), err, "groupID", groupID, "userID", userID)
	}()
	if !token_verify.IsManagerUserID(userID) && !imdb.IsGroupOwnerAdmin(groupID, userID) {
		return utils.Wrap(constant.ErrNoPermission, utils.GetSelfFuncName())
	}
	return nil
}

func (s *groupServer) GroupApplicationResponse(ctx context.Context, req *pbGroup.GroupApplicationResponseReq) (*pbGroup.GroupApplicationResponseResp, error) {
	trace_log.SetRpcReqInfo(trace_log.NewRpcCtx(ctx, utils.GetSelfFuncName(), req.OperationID), utils.GetSelfFuncName(), req.String())
	defer trace_log.ShowLog(ctx)
	resp := pbGroup.GroupApplicationResponseResp{CommonResp: &open_im_sdk.CommonResp{}}
	if err := CheckPermission(ctx, req.GroupID, req.OpUserID); err != nil {
		SetErrorForResp(err, &resp.CommonResp.ErrCode, &resp.CommonResp.ErrMsg)
		return &resp, nil
	}
	groupRequest := imdb.GroupRequest{}
	utils.CopyStructFields(&groupRequest, req)
	groupRequest.UserID = req.FromUserID
	groupRequest.HandleUserID = req.OpUserID
	groupRequest.HandledTime = time.Now()
	if err := (&imdb.GroupRequest{}).Update(ctx, []*imdb.GroupRequest{&groupRequest}); err != nil {
		SetErrorForResp(err, &resp.CommonResp.ErrCode, &resp.CommonResp.ErrMsg)
		return &resp, nil
	}
	groupInfo, err := rocksCache.GetGroupInfoFromCache(ctx, req.GroupID)
	if err != nil {
		SetErrorForResp(err, &resp.CommonResp.ErrCode, &resp.CommonResp.ErrMsg)
		return &resp, nil
	}
	if req.HandleResult == constant.GroupResponseAgree {
		user, err := imdb.GetUserByUserID(req.FromUserID)
		if err != nil {
			SetErrorForResp(err, &resp.CommonResp.ErrCode, &resp.CommonResp.ErrMsg)
			return &resp, nil
		}
		request, err := (&imdb.GroupRequest{}).Take(ctx, req.GroupID, req.FromUserID)
		if err != nil {
			SetErrorForResp(err, &resp.CommonResp.ErrCode, &resp.CommonResp.ErrMsg)
			return &resp, nil
		}
		member := imdb.GroupMember{}
		member.GroupID = req.GroupID
		member.UserID = req.FromUserID
		member.RoleLevel = constant.GroupOrdinaryUsers
		member.OperatorUserID = req.OpUserID
		member.FaceURL = user.FaceURL
		member.Nickname = user.Nickname
		member.JoinSource = request.JoinSource
		member.InviterUserID = request.InviterUserID
		member.MuteEndTime = time.Unix(int64(time.Now().Second()), 0)
		err = CallbackBeforeMemberJoinGroup(ctx, req.OperationID, &member, groupInfo.Ex)
		if err != nil {
			SetErrorForResp(err, &resp.CommonResp.ErrCode, &resp.CommonResp.ErrMsg)
			return &resp, nil
		}

		err = (&imdb.GroupMember{}).Create(ctx, []*imdb.GroupMember{&member})
		if err != nil {
			SetErrorForResp(err, &resp.CommonResp.ErrCode, &resp.CommonResp.ErrMsg)
			return &resp, nil
		}

		etcdCacheConn, err := fault_tolerant.GetDefaultConn(config.Config.RpcRegisterName.OpenImCacheName, req.OperationID)
		if err != nil {
			SetErrorForResp(err, &resp.CommonResp.ErrCode, &resp.CommonResp.ErrMsg)
			return &resp, nil
		}
		cacheClient := pbCache.NewCacheClient(etcdCacheConn)
		cacheResp, err := cacheClient.DelGroupMemberIDListFromCache(context.Background(), &pbCache.DelGroupMemberIDListFromCacheReq{OperationID: req.OperationID, GroupID: req.GroupID})
		if err != nil {
			SetErr(ctx, "DelGroupMemberIDListFromCache", err, &resp.CommonResp.ErrCode, &resp.CommonResp.ErrMsg, "groupID", req.GroupID)
			return &resp, nil
		}
		if cacheResp.CommonResp.ErrCode != 0 {
			SetErrCodeMsg(ctx, "DelGroupMemberIDListFromCache", cacheResp.CommonResp.ErrCode, cacheResp.CommonResp.ErrMsg, &resp.CommonResp.ErrCode, &resp.CommonResp.ErrMsg)
			return &resp, nil
		}
		if err := rocksCache.DelGroupMemberListHashFromCache(req.GroupID); err != nil {
			trace_log.SetContextInfo(ctx, "DelGroupMemberListHashFromCache", err, "groupID ", req.GroupID)
		}
		if err := rocksCache.DelJoinedGroupIDListFromCache(req.FromUserID); err != nil {
			trace_log.SetContextInfo(ctx, "DelJoinedGroupIDListFromCache", err, "userID ", req.FromUserID)
		}
		if err := rocksCache.DelGroupMemberNumFromCache(req.GroupID); err != nil {
			trace_log.SetContextInfo(ctx, "DelJoinedGroupIDListFromCache", err, "groupID ", req.GroupID)
		}
		chat.GroupApplicationAcceptedNotification(req)
		chat.MemberEnterNotification(req)
	} else if req.HandleResult == constant.GroupResponseRefuse {
		chat.GroupApplicationRejectedNotification(req)
	} else {
		SetErr(ctx, "", constant.ErrArgs, &resp.CommonResp.ErrCode, &resp.CommonResp.ErrMsg)
		return &resp, nil
	}

	return &resp, nil
}

func (s *groupServer) JoinGroup(ctx context.Context, req *pbGroup.JoinGroupReq) (*pbGroup.JoinGroupResp, error) {
	log.NewInfo(req.OperationID, "JoinGroup args ", req.String())
	_, err := imdb.GetUserByUserID(req.OpUserID)
	if err != nil {
		log.NewError(req.OperationID, "GetUserByUserID failed ", err.Error(), req.OpUserID)
		return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}

	groupInfo, err := rocksCache.GetGroupInfoFromCache(req.GroupID)
	if err != nil {
		log.NewError(req.OperationID, "GetGroupInfoByGroupID failed ", req.GroupID, err)
		return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	if groupInfo.Status == constant.GroupStatusDismissed {
		errMsg := " group status is dismissed "
		return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrStatus.ErrCode, ErrMsg: errMsg}}, nil
	}

	if groupInfo.NeedVerification == constant.Directly {
		if groupInfo.GroupType != constant.SuperGroup {
			us, err := imdb.GetUserByUserID(req.OpUserID)
			if err != nil {
				log.NewError(req.OperationID, "GetUserByUserID failed ", err.Error(), req.OpUserID)
				return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
			}
			//to group member
			groupMember := imdb.GroupMember{GroupID: req.GroupID, RoleLevel: constant.GroupOrdinaryUsers, OperatorUserID: req.OpUserID}
			utils.CopyStructFields(&groupMember, us)
			callbackResp := CallbackBeforeMemberJoinGroup(req.OperationID, &groupMember, groupInfo.Ex)
			if callbackResp.ErrCode != 0 {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "callbackBeforeSendSingleMsg resp: ", callbackResp)
			}
			if callbackResp.ActionCode != constant.ActionAllow {
				if callbackResp.ErrCode == 0 {
					callbackResp.ErrCode = 201
				}
				log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "callbackBeforeSendSingleMsg result", "end rpc and return", callbackResp)
				return &pbGroup.JoinGroupResp{
					CommonResp: &pbGroup.CommonResp{
						ErrCode: int32(callbackResp.ErrCode),
						ErrMsg:  callbackResp.ErrMsg,
					},
				}, nil
			}

			if err := s.DelGroupAndUserCache(req.OperationID, req.GroupID, []string{req.OpUserID}); err != nil {
				log.NewError(req.OperationID, "DelGroupAndUserCache failed, ", err.Error(), req.OpUserID)
				return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
			}

			err = imdb.InsertIntoGroupMember(groupMember)
			if err != nil {
				log.NewError(req.OperationID, "InsertIntoGroupMember failed ", err.Error(), groupMember)
				return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
			}
			//}

			var sessionType int
			if groupInfo.GroupType == constant.NormalGroup {
				sessionType = constant.GroupChatType
			} else {
				sessionType = constant.SuperGroupChatType
			}
			var reqPb pbUser.SetConversationReq
			var c pbConversation.Conversation
			reqPb.OperationID = req.OperationID
			c.OwnerUserID = req.OpUserID
			c.ConversationID = utils.GetConversationIDBySessionType(req.GroupID, sessionType)
			c.ConversationType = int32(sessionType)
			c.GroupID = req.GroupID
			c.IsNotInGroup = false
			c.UpdateUnreadCountTime = utils.GetCurrentTimestampByMill()
			reqPb.Conversation = &c
			etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
			if etcdConn == nil {
				errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
				log.NewError(req.OperationID, errMsg)
				return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: 0, ErrMsg: errMsg}}, nil
			}
			client := pbUser.NewUserClient(etcdConn)
			respPb, err := client.SetConversation(context.Background(), &reqPb)
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "SetConversation rpc failed, ", reqPb.String(), err.Error())
			} else {
				log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "SetConversation success", respPb.String())
			}

			chat.MemberEnterDirectlyNotification(req.GroupID, req.OpUserID, req.OperationID)
			log.NewInfo(req.OperationID, "JoinGroup rpc return ")
			return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}}, nil
		} else {
			log.Error(req.OperationID, "JoinGroup rpc failed, group type:  ", groupInfo.GroupType, "not support directly")
			return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrArgs.ErrCode, ErrMsg: constant.ErrArgs.ErrMsg}}, nil
		}
	}
	var groupRequest imdb.GroupRequest
	groupRequest.UserID = req.OpUserID
	groupRequest.ReqMsg = req.ReqMessage
	groupRequest.GroupID = req.GroupID
	groupRequest.JoinSource = req.JoinSource
	err = imdb.InsertIntoGroupRequest(groupRequest)
	if err != nil {
		log.NewError(req.OperationID, "InsertIntoGroupRequest failed ", err.Error(), groupRequest)
		return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	//_, err = imdb.GetGroupMemberListByGroupIDAndRoleLevel(req.GroupID, constant.GroupOwner)
	//if err != nil {
	//	log.NewError(req.OperationID, "GetGroupMemberListByGroupIDAndRoleLevel failed ", err.Error(), req.GroupID, constant.GroupOwner)
	//	return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}}, nil

	chat.JoinGroupApplicationNotification(req)
	log.NewInfo(req.OperationID, "JoinGroup rpc return ")
	return &pbGroup.JoinGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}}, nil
}

func (s *groupServer) QuitGroup(ctx context.Context, req *pbGroup.QuitGroupReq) (*pbGroup.QuitGroupResp, error) {
	log.NewInfo(req.OperationID, "QuitGroup args ", req.String())
	groupInfo, err := imdb.GetGroupInfoByGroupID(req.GroupID)
	if err != nil {
		log.NewError(req.OperationID, "ReduceGroupMemberFromCache rpc call failed ", err.Error())
		return &pbGroup.QuitGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	if groupInfo.GroupType != constant.SuperGroup {
		_, err = imdb.GetGroupMemberInfoByGroupIDAndUserID(req.GroupID, req.OpUserID)
		if err != nil {
			log.NewError(req.OperationID, "GetGroupMemberInfoByGroupIDAndUserID failed ", err.Error(), req.GroupID, req.OpUserID)
			return &pbGroup.QuitGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
		}

		if err := s.DelGroupAndUserCache(req.OperationID, req.GroupID, []string{req.OpUserID}); err != nil {
			log.NewError(req.OperationID, "DelGroupAndUserCache failed, ", err.Error(), req.OpUserID)
			return &pbGroup.QuitGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
		}

		err = imdb.DeleteGroupMemberByGroupIDAndUserID(req.GroupID, req.OpUserID)
		if err != nil {
			log.NewError(req.OperationID, "DeleteGroupMemberByGroupIdAndUserId failed ", err.Error(), req.GroupID, req.OpUserID)
			return &pbGroup.QuitGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
		}

		err = db.DB.DelGroupMember(req.GroupID, req.OpUserID)
		if err != nil {
			log.NewError(req.OperationID, "DelGroupMember failed ", req.GroupID, req.OpUserID)
			//	return &pbGroup.CommonResp{ErrorCode: constant.ErrQuitGroup.ErrCode, ErrorMsg: constant.ErrQuitGroup.ErrMsg}, nil
		}
		//modify quitter conversation info
		var reqPb pbUser.SetConversationReq
		var c pbConversation.Conversation
		reqPb.OperationID = req.OperationID
		c.OwnerUserID = req.OpUserID
		c.ConversationID = utils.GetConversationIDBySessionType(req.GroupID, constant.GroupChatType)
		c.ConversationType = constant.GroupChatType
		c.GroupID = req.GroupID
		c.IsNotInGroup = true
		reqPb.Conversation = &c
		etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
		if etcdConn == nil {
			errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
			log.NewError(req.OperationID, errMsg)
			return &pbGroup.QuitGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrInternal.ErrCode, ErrMsg: errMsg}}, nil
		}
		client := pbUser.NewUserClient(etcdConn)
		respPb, err := client.SetConversation(context.Background(), &reqPb)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "SetConversation rpc failed, ", reqPb.String(), err.Error())
		} else {
			log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "SetConversation success", respPb.String())
		}
	} else {
		okUserIDList := []string{req.OpUserID}
		if err := db.DB.RemoverUserFromSuperGroup(req.GroupID, okUserIDList); err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), req.GroupID, okUserIDList, err.Error())
			return &pbGroup.QuitGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
		}
	}

	if groupInfo.GroupType != constant.SuperGroup {
		if err := rocksCache.DelGroupMemberInfoFromCache(req.GroupID, req.OpUserID); err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.GroupID)
		}
		chat.MemberQuitNotification(req)
	} else {
		if err := rocksCache.DelJoinedSuperGroupIDListFromCache(req.OpUserID); err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.OpUserID)
		}
		if err := rocksCache.DelGroupMemberListHashFromCache(req.GroupID); err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), req.GroupID, err.Error())
		}
		chat.SuperGroupNotification(req.OperationID, req.OpUserID, req.OpUserID)
	}
	log.NewInfo(req.OperationID, "rpc QuitGroup return ", pbGroup.QuitGroupResp{CommonResp: &pbGroup.CommonResp{}})
	return &pbGroup.QuitGroupResp{CommonResp: &pbGroup.CommonResp{}}, nil
}

func hasAccess(req *pbGroup.SetGroupInfoReq) bool {
	if utils.IsContain(req.OpUserID, config.Config.Manager.AppManagerUid) {
		return true
	}
	groupUserInfo, err := imdb.GetGroupMemberInfoByGroupIDAndUserID(req.GroupInfoForSet.GroupID, req.OpUserID)
	if err != nil {
		log.NewError(req.OperationID, "GetGroupMemberInfoByGroupIDAndUserID failed, ", err.Error(), req.GroupInfoForSet.GroupID, req.OpUserID)
		return false

	}
	if groupUserInfo.RoleLevel == constant.GroupOwner || groupUserInfo.RoleLevel == constant.GroupAdmin {
		return true
	}
	return false
}

func (s *groupServer) SetGroupInfo(ctx context.Context, req *pbGroup.SetGroupInfoReq) (*pbGroup.SetGroupInfoResp, error) {
	log.NewInfo(req.OperationID, "SetGroupInfo args ", req.String())
	if !hasAccess(req) {
		log.NewError(req.OperationID, "no access ", req)
		return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}}, nil
	}

	group, err := imdb.GetGroupInfoByGroupID(req.GroupInfoForSet.GroupID)
	if err != nil {
		log.NewError(req.OperationID, "GetGroupInfoByGroupID failed ", err.Error(), req.GroupInfoForSet.GroupID)
		return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}}, nil
	}

	if group.Status == constant.GroupStatusDismissed {
		errMsg := " group status is dismissed "
		return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrStatus.ErrCode, ErrMsg: errMsg}}, nil
	}

	////bitwise operators: 0001:groupName; 0010:Notification  0100:Introduction; 1000:FaceUrl; 10000:owner
	var changedType int32
	groupName := ""
	notification := ""
	introduction := ""
	faceURL := ""
	if group.GroupName != req.GroupInfoForSet.GroupName && req.GroupInfoForSet.GroupName != "" {
		changedType = 1
		groupName = req.GroupInfoForSet.GroupName
	}
	if group.Notification != req.GroupInfoForSet.Notification && req.GroupInfoForSet.Notification != "" {
		changedType = changedType | (1 << 1)
		notification = req.GroupInfoForSet.Notification
	}
	if group.Introduction != req.GroupInfoForSet.Introduction && req.GroupInfoForSet.Introduction != "" {
		changedType = changedType | (1 << 2)
		introduction = req.GroupInfoForSet.Introduction
	}
	if group.FaceURL != req.GroupInfoForSet.FaceURL && req.GroupInfoForSet.FaceURL != "" {
		changedType = changedType | (1 << 3)
		faceURL = req.GroupInfoForSet.FaceURL
	}

	if req.GroupInfoForSet.NeedVerification != nil {
		changedType = changedType | (1 << 4)
		m := make(map[string]interface{})
		m["need_verification"] = req.GroupInfoForSet.NeedVerification.Value
		if err := imdb.UpdateGroupInfoDefaultZero(req.GroupInfoForSet.GroupID, m); err != nil {
			log.NewError(req.OperationID, "UpdateGroupInfoDefaultZero failed ", err.Error(), m)
			return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
		}
	}
	if req.GroupInfoForSet.LookMemberInfo != nil {
		changedType = changedType | (1 << 5)
		m := make(map[string]interface{})
		m["look_member_info"] = req.GroupInfoForSet.LookMemberInfo.Value
		if err := imdb.UpdateGroupInfoDefaultZero(req.GroupInfoForSet.GroupID, m); err != nil {
			log.NewError(req.OperationID, "UpdateGroupInfoDefaultZero failed ", err.Error(), m)
			return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
		}
	}
	if req.GroupInfoForSet.ApplyMemberFriend != nil {
		changedType = changedType | (1 << 6)
		m := make(map[string]interface{})
		m["apply_member_friend"] = req.GroupInfoForSet.ApplyMemberFriend.Value
		if err := imdb.UpdateGroupInfoDefaultZero(req.GroupInfoForSet.GroupID, m); err != nil {
			log.NewError(req.OperationID, "UpdateGroupInfoDefaultZero failed ", err.Error(), m)
			return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
		}
	}
	//
	//if req.RoleLevel != nil {
	//
	//}
	//only administrators can set group information
	var groupInfo imdb.Group
	utils.CopyStructFields(&groupInfo, req.GroupInfoForSet)
	if req.GroupInfoForSet.Notification != "" {
		groupInfo.NotificationUserID = req.OpUserID
		groupInfo.NotificationUpdateTime = time.Now()
	}
	if err := rocksCache.DelGroupInfoFromCache(req.GroupInfoForSet.GroupID); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "DelGroupInfoFromCache failed ", err.Error(), req.GroupInfoForSet.GroupID)
		return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	err = imdb.SetGroupInfo(groupInfo)
	if err != nil {
		log.NewError(req.OperationID, "SetGroupInfo failed ", err.Error(), groupInfo)
		return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	log.NewInfo(req.OperationID, "SetGroupInfo rpc return ", pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{}})
	if changedType != 0 {
		chat.GroupInfoSetNotification(req.OperationID, req.OpUserID, req.GroupInfoForSet.GroupID, groupName, notification, introduction, faceURL, req.GroupInfoForSet.NeedVerification)
	}
	if req.GroupInfoForSet.Notification != "" {
		//get group member user id
		getGroupMemberIDListFromCacheReq := &pbCache.GetGroupMemberIDListFromCacheReq{OperationID: req.OperationID, GroupID: req.GroupInfoForSet.GroupID}
		etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCacheName, req.OperationID)
		if etcdConn == nil {
			errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
			log.NewError(req.OperationID, errMsg)
			return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrInternal.ErrCode, ErrMsg: errMsg}}, nil
		}
		client := pbCache.NewCacheClient(etcdConn)
		cacheResp, err := client.GetGroupMemberIDListFromCache(context.Background(), getGroupMemberIDListFromCacheReq)
		if err != nil {
			log.NewError(req.OperationID, "GetGroupMemberIDListFromCache rpc call failed ", err.Error())
			return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{}}, nil
		}
		if cacheResp.CommonResp.ErrCode != 0 {
			log.NewError(req.OperationID, "GetGroupMemberIDListFromCache rpc logic call failed ", cacheResp.String())
			return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{}}, nil
		}
		var conversationReq pbConversation.ModifyConversationFieldReq

		conversation := pbConversation.Conversation{
			OwnerUserID:      req.OpUserID,
			ConversationID:   utils.GetConversationIDBySessionType(req.GroupInfoForSet.GroupID, constant.GroupChatType),
			ConversationType: constant.GroupChatType,
			GroupID:          req.GroupInfoForSet.GroupID,
		}
		conversationReq.Conversation = &conversation
		conversationReq.OperationID = req.OperationID
		conversationReq.FieldType = constant.FieldGroupAtType
		conversation.GroupAtType = constant.GroupNotification
		conversationReq.UserIDList = cacheResp.UserIDList
		nEtcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImConversationName, req.OperationID)
		if etcdConn == nil {
			errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
			log.NewError(req.OperationID, errMsg)
			return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrInternal.ErrCode, ErrMsg: errMsg}}, nil
		}
		nClient := pbConversation.NewConversationClient(nEtcdConn)
		conversationReply, err := nClient.ModifyConversationField(context.Background(), &conversationReq)
		if err != nil {
			log.NewError(conversationReq.OperationID, "ModifyConversationField rpc failed, ", conversationReq.String(), err.Error())
		} else if conversationReply.CommonResp.ErrCode != 0 {
			log.NewError(conversationReq.OperationID, "ModifyConversationField rpc failed, ", conversationReq.String(), conversationReply.String())
		}
	}
	return &pbGroup.SetGroupInfoResp{CommonResp: &pbGroup.CommonResp{}}, nil
}

func (s *groupServer) TransferGroupOwner(_ context.Context, req *pbGroup.TransferGroupOwnerReq) (*pbGroup.TransferGroupOwnerResp, error) {
	log.NewInfo(req.OperationID, "TransferGroupOwner ", req.String())

	groupInfo, err := imdb.GetGroupInfoByGroupID(req.GroupID)
	if err != nil {
		log.NewError(req.OperationID, "GetGroupInfoByGroupID failed ", req.GroupID, err)
		return &pbGroup.TransferGroupOwnerResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	if groupInfo.Status == constant.GroupStatusDismissed {
		errMsg := " group status is dismissed "
		return &pbGroup.TransferGroupOwnerResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrStatus.ErrCode, ErrMsg: errMsg}}, nil
	}

	if req.OldOwnerUserID == req.NewOwnerUserID {
		log.NewError(req.OperationID, "same owner ", req.OldOwnerUserID, req.NewOwnerUserID)
		return &pbGroup.TransferGroupOwnerResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrArgs.ErrCode, ErrMsg: constant.ErrArgs.ErrMsg}}, nil
	}
	err = rocksCache.DelGroupMemberInfoFromCache(req.GroupID, req.NewOwnerUserID)
	if err != nil {
		log.NewError(req.OperationID, "DelGroupMemberInfoFromCache failed ", req.GroupID, req.NewOwnerUserID)
		return &pbGroup.TransferGroupOwnerResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil

	}
	err = rocksCache.DelGroupMemberInfoFromCache(req.GroupID, req.OldOwnerUserID)
	if err != nil {
		log.NewError(req.OperationID, "DelGroupMemberInfoFromCache failed ", req.GroupID, req.OldOwnerUserID)
		return &pbGroup.TransferGroupOwnerResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}

	groupMemberInfo := imdb.GroupMember{GroupID: req.GroupID, UserID: req.OldOwnerUserID, RoleLevel: constant.GroupOrdinaryUsers}
	err = imdb.UpdateGroupMemberInfo(groupMemberInfo)
	if err != nil {
		log.NewError(req.OperationID, "UpdateGroupMemberInfo failed ", groupMemberInfo)
		return &pbGroup.TransferGroupOwnerResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	groupMemberInfo = imdb.GroupMember{GroupID: req.GroupID, UserID: req.NewOwnerUserID, RoleLevel: constant.GroupOwner}
	err = imdb.UpdateGroupMemberInfo(groupMemberInfo)
	if err != nil {
		log.NewError(req.OperationID, "UpdateGroupMemberInfo failed ", groupMemberInfo)
		return &pbGroup.TransferGroupOwnerResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}

	chat.GroupOwnerTransferredNotification(req)
	return &pbGroup.TransferGroupOwnerResp{CommonResp: &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}}, nil

}

func (s *groupServer) GetGroups(_ context.Context, req *pbGroup.GetGroupsReq) (*pbGroup.GetGroupsResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "GetGroups ", req.String())
	resp := &pbGroup.GetGroupsResp{
		CommonResp: &pbGroup.CommonResp{},
		CMSGroups:  []*pbGroup.CMSGroup{},
		Pagination: &open_im_sdk.ResponsePagination{CurrentPage: req.Pagination.PageNumber, ShowNumber: req.Pagination.ShowNumber},
	}
	if req.GroupID != "" {
		groupInfoDB, err := imdb.GetGroupInfoByGroupID(req.GroupID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return resp, nil
			}
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.GroupID)
			resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resp.CommonResp.ErrMsg = err.Error()
			return resp, nil
		}
		resp.GroupNum = 1
		groupInfo := &open_im_sdk.GroupInfo{}
		utils.CopyStructFields(groupInfo, groupInfoDB)
		groupMember, err := imdb.GetGroupOwnerInfoByGroupID(req.GroupID)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.GroupID)
			resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resp.CommonResp.ErrMsg = err.Error()
			return resp, nil
		}
		memberNum, err := imdb.GetGroupMembersCount(req.GroupID, "")
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.GroupID)
			resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resp.CommonResp.ErrMsg = err.Error()
			return resp, nil
		}
		groupInfo.MemberCount = uint32(memberNum)
		groupInfo.CreateTime = uint32(groupInfoDB.CreateTime.Unix())
		resp.CMSGroups = append(resp.CMSGroups, &pbGroup.CMSGroup{GroupInfo: groupInfo, GroupOwnerUserName: groupMember.Nickname, GroupOwnerUserID: groupMember.UserID})
	} else {
		groups, count, err := imdb.GetGroupsByName(req.GroupName, req.Pagination.PageNumber, req.Pagination.ShowNumber)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetGroupsByName error", req.String(), req.GroupName, req.Pagination.PageNumber, req.Pagination.ShowNumber)
		}
		for _, v := range groups {
			group := &pbGroup.CMSGroup{GroupInfo: &open_im_sdk.GroupInfo{}}
			utils.CopyStructFields(group.GroupInfo, v)
			groupMember, err := imdb.GetGroupOwnerInfoByGroupID(v.GroupID)
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetGroupOwnerInfoByGroupID failed", err.Error(), v)
				continue
			}
			group.GroupInfo.CreateTime = uint32(v.CreateTime.Unix())
			group.GroupOwnerUserID = groupMember.UserID
			group.GroupOwnerUserName = groupMember.Nickname
			resp.CMSGroups = append(resp.CMSGroups, group)
		}
		resp.GroupNum = int32(count)
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "GetGroups resp", resp.String())
	return resp, nil
}

func (s *groupServer) GetGroupMembersCMS(_ context.Context, req *pbGroup.GetGroupMembersCMSReq) (*pbGroup.GetGroupMembersCMSResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "args:", req.String())
	resp := &pbGroup.GetGroupMembersCMSResp{CommonResp: &pbGroup.CommonResp{}}
	groupMembers, err := imdb.GetGroupMembersByGroupIdCMS(req.GroupID, req.UserName, req.Pagination.ShowNumber, req.Pagination.PageNumber)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetGroupMembersByGroupIdCMS Error", err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	groupMembersCount, err := imdb.GetGroupMembersCount(req.GroupID, req.UserName)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetGroupMembersCMS Error", err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	log.NewInfo(req.OperationID, groupMembersCount)
	resp.MemberNums = int32(groupMembersCount)
	for _, groupMember := range groupMembers {
		member := open_im_sdk.GroupMemberFullInfo{}
		utils.CopyStructFields(&member, groupMember)
		member.JoinTime = int32(groupMember.JoinTime.Unix())
		member.MuteEndTime = uint32(groupMember.MuteEndTime.Unix())
		resp.Members = append(resp.Members, &member)
	}
	resp.Pagination = &open_im_sdk.ResponsePagination{
		CurrentPage: req.Pagination.PageNumber,
		ShowNumber:  req.Pagination.ShowNumber,
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp:", resp.String())
	return resp, nil
}

func (s *groupServer) GetUserReqApplicationList(_ context.Context, req *pbGroup.GetUserReqApplicationListReq) (*pbGroup.GetUserReqApplicationListResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &pbGroup.GetUserReqApplicationListResp{}
	groupRequests, err := imdb.GetUserReqGroupByUserID(req.UserID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUserReqGroupByUserID failed ", err.Error())
		resp.CommonResp = &pbGroup.CommonResp{
			ErrCode: constant.ErrDB.ErrCode,
			ErrMsg:  constant.ErrDB.ErrMsg,
		}
		return resp, nil
	}
	for _, groupReq := range groupRequests {
		node := open_im_sdk.GroupRequest{UserInfo: &open_im_sdk.PublicUserInfo{}, GroupInfo: &open_im_sdk.GroupInfo{}}
		group, err := imdb.GetGroupInfoByGroupID(groupReq.GroupID)
		if err != nil {
			log.Error(req.OperationID, "GetGroupInfoByGroupID failed ", err.Error(), groupReq.GroupID)
			continue
		}
		user, err := imdb.GetUserByUserID(groupReq.UserID)
		if err != nil {
			log.Error(req.OperationID, "GetUserByUserID failed ", err.Error(), groupReq.UserID)
			continue
		}
		cp.GroupRequestDBCopyOpenIM(&node, &groupReq)
		cp.UserDBCopyOpenIMPublicUser(node.UserInfo, user)
		cp.GroupDBCopyOpenIM(node.GroupInfo, group)
		resp.GroupRequestList = append(resp.GroupRequestList, &node)
	}
	resp.CommonResp = &pbGroup.CommonResp{
		ErrCode: 0,
		ErrMsg:  "",
	}
	return resp, nil
}

func (s *groupServer) DismissGroup(ctx context.Context, req *pbGroup.DismissGroupReq) (*pbGroup.DismissGroupResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc args ", req.String())
	if !token_verify.IsManagerUserID(req.OpUserID) && !imdb.IsGroupOwnerAdmin(req.GroupID, req.OpUserID) {
		log.NewError(req.OperationID, "verify failed ", req.OpUserID, req.GroupID)
		return &pbGroup.DismissGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}}, nil
	}

	if err := rocksCache.DelGroupInfoFromCache(req.GroupID); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.GroupID)
		return &pbGroup.DismissGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}
	if err := s.DelGroupAndUserCache(req.OperationID, req.GroupID, nil); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.GroupID)
		return &pbGroup.DismissGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}

	err := imdb.OperateGroupStatus(req.GroupID, constant.GroupStatusDismissed)
	if err != nil {
		log.NewError(req.OperationID, "OperateGroupStatus failed ", req.GroupID, constant.GroupStatusDismissed)
		return &pbGroup.DismissGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	groupInfo, err := imdb.GetGroupInfoByGroupID(req.GroupID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.GroupID)
		return &pbGroup.DismissGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	if groupInfo.GroupType != constant.SuperGroup {
		memberList, err := imdb.GetGroupMemberListByGroupID(req.GroupID)
		if err != nil {
			log.NewError(req.OperationID, "GetGroupMemberListByGroupID failed,", err.Error(), req.GroupID)
		}
		//modify quitter conversation info
		var reqPb pbUser.SetConversationReq
		var c pbConversation.Conversation
		for _, v := range memberList {
			reqPb.OperationID = req.OperationID
			c.OwnerUserID = v.UserID
			c.ConversationID = utils.GetConversationIDBySessionType(req.GroupID, constant.GroupChatType)
			c.ConversationType = constant.GroupChatType
			c.GroupID = req.GroupID
			c.IsNotInGroup = true
			reqPb.Conversation = &c
			etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName, req.OperationID)
			if etcdConn == nil {
				errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
				log.NewError(req.OperationID, errMsg)
				return &pbGroup.DismissGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrInternal.ErrCode, ErrMsg: errMsg}}, nil
			}
			client := pbUser.NewUserClient(etcdConn)
			respPb, err := client.SetConversation(context.Background(), &reqPb)
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "SetConversation rpc failed, ", reqPb.String(), err.Error(), v.UserID)
			} else {
				log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "SetConversation success", respPb.String(), v.UserID)
			}
		}
		err = imdb.DeleteGroupMemberByGroupID(req.GroupID)
		if err != nil {
			log.NewError(req.OperationID, "DeleteGroupMemberByGroupID failed ", req.GroupID)
			return &pbGroup.DismissGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
		}
		chat.GroupDismissedNotification(req)
	} else {
		err = db.DB.DeleteSuperGroup(req.GroupID)
		if err != nil {
			log.NewError(req.OperationID, "DeleteGroupMemberByGroupID failed ", req.GroupID)
			return &pbGroup.DismissGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
		}
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc return ", pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""})
	return &pbGroup.DismissGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}}, nil
}

//  rpc MuteGroupMember(MuteGroupMemberReq) returns(MuteGroupMemberResp);
//  rpc CancelMuteGroupMember(CancelMuteGroupMemberReq) returns(CancelMuteGroupMemberResp);
//  rpc MuteGroup(MuteGroupReq) returns(MuteGroupResp);
//  rpc CancelMuteGroup(CancelMuteGroupReq) returns(CancelMuteGroupResp);

func (s *groupServer) MuteGroupMember(ctx context.Context, req *pbGroup.MuteGroupMemberReq) (*pbGroup.MuteGroupMemberResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc args ", req.String())
	opFlag, err := s.getGroupUserLevel(req.GroupID, req.OpUserID)
	if err != nil {
		errMsg := req.OperationID + " getGroupUserLevel failed " + req.GroupID + req.OpUserID + err.Error()
		log.Error(req.OperationID, errMsg)
		return &pbGroup.MuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}
	if opFlag == 0 {
		errMsg := req.OperationID + "opFlag == 0  " + req.GroupID + req.OpUserID
		log.Error(req.OperationID, errMsg)
		return &pbGroup.MuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}

	mutedInfo, err := rocksCache.GetGroupMemberInfoFromCache(req.GroupID, req.UserID)
	if err != nil {
		errMsg := req.OperationID + " GetGroupMemberInfoByGroupIDAndUserID failed " + req.GroupID + req.UserID + err.Error()
		return &pbGroup.MuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}
	if mutedInfo.RoleLevel == constant.GroupOwner && opFlag != 1 {
		errMsg := req.OperationID + " mutedInfo.RoleLevel == constant.GroupOwner " + req.GroupID + req.UserID
		return &pbGroup.MuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}
	if mutedInfo.RoleLevel == constant.GroupAdmin && opFlag == 3 {
		errMsg := req.OperationID + " mutedInfo.RoleLevel == constant.GroupAdmin " + req.GroupID + req.UserID
		return &pbGroup.MuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}

	if err := rocksCache.DelGroupMemberInfoFromCache(req.GroupID, req.UserID); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.GroupID, req.UserID)
		return &pbGroup.MuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}
	groupMemberInfo := imdb.GroupMember{GroupID: req.GroupID, UserID: req.UserID}
	groupMemberInfo.MuteEndTime = time.Unix(int64(time.Now().Second())+int64(req.MutedSeconds), time.Now().UnixNano())
	err = imdb.UpdateGroupMemberInfo(groupMemberInfo)
	if err != nil {
		log.Error(req.OperationID, "UpdateGroupMemberInfo failed ", err.Error(), groupMemberInfo)
		return &pbGroup.MuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	chat.GroupMemberMutedNotification(req.OperationID, req.OpUserID, req.GroupID, req.UserID, req.MutedSeconds)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc return ", pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""})
	return &pbGroup.MuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}}, nil
}

func (s *groupServer) CancelMuteGroupMember(ctx context.Context, req *pbGroup.CancelMuteGroupMemberReq) (*pbGroup.CancelMuteGroupMemberResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc args ", req.String())

	opFlag, err := s.getGroupUserLevel(req.GroupID, req.OpUserID)
	if err != nil {
		errMsg := req.OperationID + " getGroupUserLevel failed " + req.GroupID + req.OpUserID + err.Error()
		log.Error(req.OperationID, errMsg)
		return &pbGroup.CancelMuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}
	if opFlag == 0 {
		errMsg := req.OperationID + "opFlag == 0  " + req.GroupID + req.OpUserID
		log.Error(req.OperationID, errMsg)
		return &pbGroup.CancelMuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}

	mutedInfo, err := imdb.GetGroupMemberInfoByGroupIDAndUserID(req.GroupID, req.UserID)
	if err != nil {
		errMsg := req.OperationID + " GetGroupMemberInfoByGroupIDAndUserID failed " + req.GroupID + req.UserID + err.Error()
		return &pbGroup.CancelMuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}
	if mutedInfo.RoleLevel == constant.GroupOwner && opFlag != 1 {
		errMsg := req.OperationID + " mutedInfo.RoleLevel == constant.GroupOwner " + req.GroupID + req.UserID
		return &pbGroup.CancelMuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}
	if mutedInfo.RoleLevel == constant.GroupAdmin && opFlag == 3 {
		errMsg := req.OperationID + " mutedInfo.RoleLevel == constant.GroupAdmin " + req.GroupID + req.UserID
		return &pbGroup.CancelMuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}
	if err := rocksCache.DelGroupMemberInfoFromCache(req.GroupID, req.UserID); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.GroupID)
		return &pbGroup.CancelMuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}

	groupMemberInfo := imdb.GroupMember{GroupID: req.GroupID, UserID: req.UserID}
	groupMemberInfo.MuteEndTime = time.Unix(0, 0)
	err = imdb.UpdateGroupMemberInfo(groupMemberInfo)
	if err != nil {
		log.Error(req.OperationID, "UpdateGroupMemberInfo failed ", err.Error(), groupMemberInfo)
		return &pbGroup.CancelMuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}

	chat.GroupMemberCancelMutedNotification(req.OperationID, req.OpUserID, req.GroupID, req.UserID)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc return ", pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""})
	return &pbGroup.CancelMuteGroupMemberResp{CommonResp: &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}}, nil
}

func (s *groupServer) MuteGroup(ctx context.Context, req *pbGroup.MuteGroupReq) (*pbGroup.MuteGroupResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc args ", req.String())

	opFlag, err := s.getGroupUserLevel(req.GroupID, req.OpUserID)
	if err != nil {
		errMsg := req.OperationID + " getGroupUserLevel failed " + req.GroupID + req.OpUserID + err.Error()
		log.Error(req.OperationID, errMsg)
		return &pbGroup.MuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}
	if opFlag == 0 {
		errMsg := req.OperationID + "opFlag == 0  " + req.GroupID + req.OpUserID
		log.Error(req.OperationID, errMsg)
		return &pbGroup.MuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}

	//mutedInfo, err := imdb.GetGroupMemberInfoByGroupIDAndUserID(req.GroupID, req.UserID)
	//if err != nil {
	//	errMsg := req.OperationID + " GetGroupMemberInfoByGroupIDAndUserID failed " + req.GroupID + req.OpUserID + err.Error()
	//	return &pbGroup.MuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	//}
	//if mutedInfo.RoleLevel == constant.GroupOwner && opFlag != 1 {
	//	errMsg := req.OperationID + " mutedInfo.RoleLevel == constant.GroupOwner " + req.GroupID + req.OpUserID + err.Error()
	//	return &pbGroup.MuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	//}
	//if mutedInfo.RoleLevel == constant.GroupAdmin && opFlag == 3 {
	//	errMsg := req.OperationID + " mutedInfo.RoleLevel == constant.GroupAdmin " + req.GroupID + req.OpUserID + err.Error()
	//	return &pbGroup.MuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	//}
	if err := rocksCache.DelGroupInfoFromCache(req.GroupID); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.GroupID)
		return &pbGroup.MuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}

	err = imdb.OperateGroupStatus(req.GroupID, constant.GroupStatusMuted)
	if err != nil {
		log.Error(req.OperationID, "OperateGroupStatus failed ", err.Error(), req.GroupID, constant.GroupStatusMuted)
		return &pbGroup.MuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}

	chat.GroupMutedNotification(req.OperationID, req.OpUserID, req.GroupID)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc return ", pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""})
	return &pbGroup.MuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}}, nil
}

func (s *groupServer) CancelMuteGroup(ctx context.Context, req *pbGroup.CancelMuteGroupReq) (*pbGroup.CancelMuteGroupResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc args ", req.String())
	opFlag, err := s.getGroupUserLevel(req.GroupID, req.OpUserID)
	if err != nil {
		errMsg := req.OperationID + " getGroupUserLevel failed " + req.GroupID + req.OpUserID + err.Error()
		log.Error(req.OperationID, errMsg)
		return &pbGroup.CancelMuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}
	if opFlag == 0 {
		errMsg := req.OperationID + "opFlag == 0  " + req.GroupID + req.OpUserID
		log.Error(req.OperationID, errMsg)
		return &pbGroup.CancelMuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}
	//mutedInfo, err := imdb.GetGroupMemberInfoByGroupIDAndUserID(req.GroupID, req.)
	//if err != nil {
	//	errMsg := req.OperationID + " GetGroupMemberInfoByGroupIDAndUserID failed " + req.GroupID + req.OpUserID + err.Error()
	//	return &pbGroup.CancelMuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	//}
	//if mutedInfo.RoleLevel == constant.GroupOwner && opFlag != 1 {
	//	errMsg := req.OperationID + " mutedInfo.RoleLevel == constant.GroupOwner " + req.GroupID + req.OpUserID + err.Error()
	//	return &pbGroup.CancelMuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	//}
	//if mutedInfo.RoleLevel == constant.GroupAdmin && opFlag == 3 {
	//	errMsg := req.OperationID + " mutedInfo.RoleLevel == constant.GroupAdmin " + req.GroupID + req.OpUserID + err.Error()
	//	return &pbGroup.CancelMuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	//}
	log.Debug(req.OperationID, "UpdateGroupInfoDefaultZero ", req.GroupID, map[string]interface{}{"status": constant.GroupOk})
	if err := rocksCache.DelGroupInfoFromCache(req.GroupID); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.GroupID)
		return &pbGroup.CancelMuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}
	err = imdb.UpdateGroupInfoDefaultZero(req.GroupID, map[string]interface{}{"status": constant.GroupOk})
	if err != nil {
		log.Error(req.OperationID, "UpdateGroupInfoDefaultZero failed ", err.Error(), req.GroupID)
		return &pbGroup.CancelMuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}

	chat.GroupCancelMutedNotification(req.OperationID, req.OpUserID, req.GroupID)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc return ", pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""})
	return &pbGroup.CancelMuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}}, nil
}

func (s *groupServer) SetGroupMemberNickname(ctx context.Context, req *pbGroup.SetGroupMemberNicknameReq) (*pbGroup.SetGroupMemberNicknameResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc args ", req.String())
	if req.OpUserID != req.UserID && !token_verify.IsManagerUserID(req.OpUserID) {
		errMsg := req.OperationID + " verify failed " + req.OpUserID + req.GroupID
		log.Error(req.OperationID, errMsg)
		return &pbGroup.SetGroupMemberNicknameResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}}, nil
	}
	cbReq := &pbGroup.SetGroupMemberInfoReq{
		GroupID:     req.GroupID,
		UserID:      req.UserID,
		OperationID: req.OperationID,
		OpUserID:    req.OpUserID,
		Nickname:    &wrapperspb.StringValue{Value: req.Nickname},
	}
	callbackResp := CallbackBeforeSetGroupMemberInfo(cbReq)
	if callbackResp.ErrCode != 0 {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "CallbackBeforeMemberJoinGroup resp: ", callbackResp)
	}
	if callbackResp.ActionCode != constant.ActionAllow {
		if callbackResp.ErrCode == 0 {
			callbackResp.ErrCode = 201
		}
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CallbackBeforeMemberJoinGroup result", "end rpc and return", callbackResp)
		return &pbGroup.SetGroupMemberNicknameResp{
			CommonResp: &pbGroup.CommonResp{
				ErrCode: int32(callbackResp.ErrCode),
				ErrMsg:  callbackResp.ErrMsg,
			},
		}, nil
	}

	nickName := cbReq.Nickname.Value
	groupMemberInfo := imdb.GroupMember{}
	groupMemberInfo.UserID = req.UserID
	groupMemberInfo.GroupID = req.GroupID
	if nickName == "" {
		userNickname, err := imdb.GetUserNameByUserID(groupMemberInfo.UserID)
		if err != nil {
			errMsg := req.OperationID + " GetUserNameByUserID failed " + err.Error()
			log.Error(req.OperationID, errMsg)
			return &pbGroup.SetGroupMemberNicknameResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
		}
		groupMemberInfo.Nickname = userNickname
	} else {
		groupMemberInfo.Nickname = nickName
	}

	if err := rocksCache.DelGroupMemberInfoFromCache(req.GroupID, req.UserID); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.GroupID)
		return &pbGroup.SetGroupMemberNicknameResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}
	err := imdb.UpdateGroupMemberInfo(groupMemberInfo)
	if err != nil {
		errMsg := req.OperationID + " UpdateGroupMemberInfo failed " + err.Error()
		log.Error(req.OperationID, errMsg)
		return &pbGroup.SetGroupMemberNicknameResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	chat.GroupMemberInfoSetNotification(req.OperationID, req.OpUserID, req.GroupID, req.UserID)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc return ", pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""})
	return &pbGroup.SetGroupMemberNicknameResp{CommonResp: &pbGroup.CommonResp{ErrCode: 0, ErrMsg: ""}}, nil
}

func (s *groupServer) SetGroupMemberInfo(ctx context.Context, req *pbGroup.SetGroupMemberInfoReq) (resp *pbGroup.SetGroupMemberInfoResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp = &pbGroup.SetGroupMemberInfoResp{CommonResp: &pbGroup.CommonResp{}}
	if err := rocksCache.DelGroupMemberInfoFromCache(req.GroupID, req.UserID); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.GroupID)
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	callbackResp := CallbackBeforeSetGroupMemberInfo(req)
	if callbackResp.ErrCode != 0 {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "CallbackBeforeMemberJoinGroup resp: ", callbackResp)
	}
	if callbackResp.ActionCode != constant.ActionAllow {
		if callbackResp.ErrCode == 0 {
			callbackResp.ErrCode = 201
		}
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CallbackBeforeMemberJoinGroup result", "end rpc and return", callbackResp)
		return &pbGroup.SetGroupMemberInfoResp{
			CommonResp: &pbGroup.CommonResp{
				ErrCode: int32(callbackResp.ErrCode),
				ErrMsg:  callbackResp.ErrMsg,
			},
		}, nil
	}

	groupMember := imdb.GroupMember{
		GroupID: req.GroupID,
		UserID:  req.UserID,
	}
	m := make(map[string]interface{})
	if req.RoleLevel != nil {
		m["role_level"] = req.RoleLevel.Value
	}
	if req.FaceURL != nil {
		m["user_group_face_url"] = req.FaceURL.Value
	}
	if req.Nickname != nil {
		m["nickname"] = req.Nickname.Value
	}
	if req.Ex != nil {
		m["ex"] = req.Ex.Value
	} else {
		m["ex"] = nil
	}
	err = imdb.UpdateGroupMemberInfoByMap(groupMember, m)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "SetGroupMemberInfo failed", err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg + ":" + err.Error()
		return resp, nil
	}
	if req.RoleLevel != nil {
		switch req.RoleLevel.Value {
		case constant.GroupOrdinaryUsers:
			//msg.GroupMemberRoleLevelChangeNotification(req.OperationID, req.OpUserID, req.GroupID, req.UserID, constant.GroupMemberSetToOrdinaryUserNotification)
			chat.GroupMemberInfoSetNotification(req.OperationID, req.OpUserID, req.GroupID, req.UserID)
		case constant.GroupAdmin, constant.GroupOwner:
			//msg.GroupMemberRoleLevelChangeNotification(req.OperationID, req.OpUserID, req.GroupID, req.UserID, constant.GroupMemberSetToAdminNotification)
			chat.GroupMemberInfoSetNotification(req.OperationID, req.OpUserID, req.GroupID, req.UserID)
		}
	} else {
		chat.GroupMemberInfoSetNotification(req.OperationID, req.OpUserID, req.GroupID, req.UserID)
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

func (s *groupServer) GetGroupAbstractInfo(c context.Context, req *pbGroup.GetGroupAbstractInfoReq) (*pbGroup.GetGroupAbstractInfoResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &pbGroup.GetGroupAbstractInfoResp{CommonResp: &pbGroup.CommonResp{}}
	hashCode, err := rocksCache.GetGroupMemberListHashFromCache(req.GroupID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetGroupMemberListHashFromCache failed", req.GroupID, err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	resp.GroupMemberListHash = hashCode
	num, err := rocksCache.GetGroupMemberNumFromCache(req.GroupID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetGroupMemberNumByGroupID failed", req.GroupID, err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	resp.GroupMemberNumber = int32(num)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", resp.String())
	return resp, nil
}

func (s *groupServer) DelGroupAndUserCache(operationID, groupID string, userIDList []string) error {
	if groupID != "" {
		etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCacheName, operationID)
		if etcdConn == nil {
			errMsg := operationID + "getcdv3.GetDefaultConn == nil"
			log.NewError(operationID, errMsg)
			return errors.New("etcdConn is nil")
		}
		cacheClient := pbCache.NewCacheClient(etcdConn)
		cacheResp, err := cacheClient.DelGroupMemberIDListFromCache(context.Background(), &pbCache.DelGroupMemberIDListFromCacheReq{
			GroupID:     groupID,
			OperationID: operationID,
		})
		if err != nil {
			log.NewError(operationID, "DelGroupMemberIDListFromCache rpc call failed ", err.Error())
			return utils.Wrap(err, "")
		}
		if cacheResp.CommonResp.ErrCode != 0 {
			log.NewError(operationID, "DelGroupMemberIDListFromCache rpc logic call failed ", cacheResp.String())
			return errors.New(fmt.Sprintf("errMsg is %s, errCode is %d", cacheResp.CommonResp.ErrMsg, cacheResp.CommonResp.ErrCode))
		}
		err = rocksCache.DelGroupMemberListHashFromCache(groupID)
		if err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), groupID, err.Error())
			return utils.Wrap(err, "")
		}
		err = rocksCache.DelGroupMemberNumFromCache(groupID)
		if err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), groupID)
			return utils.Wrap(err, "")
		}
	}
	if userIDList != nil {
		for _, userID := range userIDList {
			err := rocksCache.DelJoinedGroupIDListFromCache(userID)
			if err != nil {
				log.NewError(operationID, utils.GetSelfFuncName(), err.Error())
				return utils.Wrap(err, "")
			}
		}
	}
	return nil
}
