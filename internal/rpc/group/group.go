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
	"Open_IM/pkg/common/middleware"
	promePkg "Open_IM/pkg/common/prometheus"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/common/trace_log"
	cp "Open_IM/pkg/common/utils"
	"Open_IM/pkg/getcdv3"
	pbCache "Open_IM/pkg/proto/cache"
	pbConversation "Open_IM/pkg/proto/conversation"
	pbGroup "Open_IM/pkg/proto/group"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	pbUser "Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"
	"context"
	"errors"
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
		grpc.UnaryInterceptor(middleware.RpcServerInterceptor),
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

func OperationID(ctx context.Context) string {
	s, _ := ctx.Value("operationID").(string)
	return s
}

func OpUserID(ctx context.Context) string {
	s, _ := ctx.Value("opUserID").(string)
	return s
}

func (s *groupServer) CreateGroup(ctx context.Context, req *pbGroup.CreateGroupReq) (*pbGroup.CreateGroupResp, error) {

	resp := &pbGroup.CreateGroupResp{GroupInfo: &open_im_sdk.GroupInfo{}}
	if err := token_verify.CheckAccessV2(ctx, req.OpUserID, req.OwnerUserID); err != nil {
		return nil, err
	}
	var groupOwnerNum int
	var userIDs []string
	for _, info := range req.InitMemberList {
		if info.RoleLevel == constant.GroupOwner {
			groupOwnerNum++
		}
		userIDs = append(userIDs, info.UserID)
	}
	if req.OwnerUserID != "" {
		groupOwnerNum++
		userIDs = append(userIDs, req.OwnerUserID)
	}
	if groupOwnerNum != 1 {
		return nil, utils.Wrap(constant.ErrArgs, "")
	}
	if utils.IsRepeatStringSlice(userIDs) {
		return nil, utils.Wrap(constant.ErrArgs, "")
	}
	users, err := rocksCache.GetUserInfoFromCacheBatch(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	if len(users) != len(userIDs) {
		return nil, utils.Wrap(constant.ErrArgs, "")
	}
	userMap := make(map[string]*imdb.User)
	for i, user := range users {
		userMap[user.UserID] = users[i]
	}
	if err := s.DelGroupAndUserCache(ctx, "", userIDs); err != nil {
		return nil, err
	}
	if err := callbackBeforeCreateGroup(ctx, req); err != nil {
		return nil, err
	}
	groupId := req.GroupInfo.GroupID
	if groupId == "" {
		groupId = utils.Md5(req.OperationID + strconv.FormatInt(time.Now().UnixNano(), 10))
		bi := big.NewInt(0)
		bi.SetString(groupId[0:8], 16)
		groupId = bi.String()
	}
	groupInfo := imdb.Group{}
	utils.CopyStructFields(&groupInfo, req.GroupInfo)
	groupInfo.CreatorUserID = req.OpUserID
	groupInfo.GroupID = groupId
	groupInfo.CreateTime = time.Now()
	if groupInfo.NotificationUpdateTime.Unix() < 0 {
		groupInfo.NotificationUpdateTime = utils.UnixSecondToTime(0)
	}
	if req.GroupInfo.GroupType != constant.SuperGroup {
		var groupMembers []*imdb.GroupMember
		joinGroup := func(userID string, roleLevel int32) error {
			groupMember := &imdb.GroupMember{GroupID: groupId, RoleLevel: roleLevel, OperatorUserID: req.OpUserID, JoinSource: constant.JoinByInvitation, InviterUserID: req.OpUserID}
			user := userMap[userID]
			utils.CopyStructFields(&groupMember, user)
			if err := CallbackBeforeMemberJoinGroup(ctx, req.OperationID, groupMember, groupInfo.Ex); err != nil {
				return err
			}
			groupMembers = append(groupMembers, groupMember)
			return nil
		}
		if req.OwnerUserID == "" {
			if err := joinGroup(req.OwnerUserID, constant.GroupOwner); err != nil {
				return nil, err
			}
		}
		for _, info := range req.InitMemberList {
			if err := joinGroup(info.UserID, info.RoleLevel); err != nil {
				return nil, err
			}
		}
		if err := (*imdb.GroupMember)(nil).Create(ctx, groupMembers); err != nil {
			return nil, err
		}
	} else {
		if err := db.DB.CreateSuperGroup(groupId, userIDs, len(userIDs)); err != nil {
			return nil, err
		}
	}
	if err := (*imdb.Group)(nil).Create(ctx, []*imdb.Group{&groupInfo}); err != nil {
		return nil, err
	}
	utils.CopyStructFields(resp.GroupInfo, groupInfo)
	resp.GroupInfo.MemberCount = uint32(len(userIDs))
	if req.GroupInfo.GroupType != constant.SuperGroup {
		chat.GroupCreatedNotification(req.OperationID, req.OpUserID, groupId, userIDs)
	} else {
		for _, userID := range userIDs {
			if err := rocksCache.DelJoinedSuperGroupIDListFromCache(ctx, userID); err != nil {
				trace_log.SetCtxInfo(ctx, "DelJoinedSuperGroupIDListFromCache", err, "userID", userID)
			}
		}
		go func() {
			for _, v := range userIDs {
				chat.SuperGroupNotification(req.OperationID, v, v)
			}
		}()
	}
	return resp, nil
}

func (s *groupServer) GetJoinedGroupList(ctx context.Context, req *pbGroup.GetJoinedGroupListReq) (*pbGroup.GetJoinedGroupListResp, error) {
	resp := &pbGroup.GetJoinedGroupListResp{}

	if err := token_verify.CheckAccessV2(ctx, req.OpUserID, req.FromUserID); err != nil {
		return nil, err
	}
	joinedGroupList, err := rocksCache.GetJoinedGroupIDListFromCache(ctx, req.FromUserID)
	if err != nil {
		return nil, err
	}
	for _, groupID := range joinedGroupList {
		var groupNode open_im_sdk.GroupInfo
		num, err := rocksCache.GetGroupMemberNumFromCache(ctx, groupID)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), groupID)
			continue
		}
		owner, err := (*imdb.GroupMember)(nil).TakeOwnerInfo(ctx, groupID)
		//owner, err2 := imdb.GetGroupOwnerInfoByGroupID(groupID)
		if err != nil {
			continue
		}
		group, err := rocksCache.GetGroupInfoFromCache(ctx, groupID)
		if err != nil {
			continue
		}
		if group.GroupType == constant.SuperGroup {
			continue
		}
		if group.Status == constant.GroupStatusDismissed {
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
	}
	return resp, nil
}

func (s *groupServer) InviteUserToGroup(ctx context.Context, req *pbGroup.InviteUserToGroupReq) (*pbGroup.InviteUserToGroupResp, error) {
	resp := &pbGroup.InviteUserToGroupResp{}

	if !imdb.IsExistGroupMember(req.GroupID, req.OpUserID) && !token_verify.IsManagerUserID(req.OpUserID) {
		constant.SetErrorForResp(constant.ErrIdentity, resp.CommonResp)
		return nil, utils.Wrap(constant.ErrIdentity, "")
	}
	groupInfo, err := (*imdb.Group)(nil).Take(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if groupInfo.Status == constant.GroupStatusDismissed {
		return nil, utils.Wrap(constant.ErrDismissedAlready, "")
	}
	if groupInfo.NeedVerification == constant.AllNeedVerification &&
		!imdb.IsGroupOwnerAdmin(req.GroupID, req.OpUserID) && !token_verify.IsManagerUserID(req.OpUserID) {
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
		return resp, nil
	}
	if err := s.DelGroupAndUserCache(ctx, req.GroupID, req.InvitedUserIDList); err != nil {
		return nil, err
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
				trace_log.SetCtxInfo(ctx, "GetUserByUserID", err, "userID", v)
				resultNode.Result = -1
				resp.Id2ResultList = append(resp.Id2ResultList, &resultNode)
				continue
			}

			if imdb.IsExistGroupMember(req.GroupID, v) {
				trace_log.SetCtxInfo(ctx, "IsExistGroupMember", err, "groupID", req.GroupID, "userID", v)
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
			if err := CallbackBeforeMemberJoinGroup(ctx, req.OperationID, &toInsertInfo, groupInfo.Ex); err != nil {
				return nil, err
			}
			err = imdb.InsertIntoGroupMember(toInsertInfo)
			if err != nil {
				trace_log.SetCtxInfo(ctx, "InsertIntoGroupMember", err, "args", toInsertInfo)
				resultNode.Result = -1
				resp.Id2ResultList = append(resp.Id2ResultList, &resultNode)
				continue
			}
			okUserIDList = append(okUserIDList, v)
			err = db.DB.AddGroupMember(req.GroupID, toUserInfo.UserID)
			if err != nil {
				trace_log.SetCtxInfo(ctx, "AddGroupMember", err, "groupID", req.GroupID, "userID", toUserInfo.UserID)
			}
			resp.Id2ResultList = append(resp.Id2ResultList, &resultNode)
		}
	} else {
		okUserIDList = req.InvitedUserIDList
		if err := db.DB.AddUserToSuperGroup(req.GroupID, req.InvitedUserIDList); err != nil {
			return nil, err
		}
	}

	if groupInfo.GroupType != constant.SuperGroup {
		chat.MemberInvitedNotification(req.OperationID, req.GroupID, req.OpUserID, req.Reason, okUserIDList)
	} else {
		for _, userID := range req.InvitedUserIDList {
			if err := rocksCache.DelJoinedSuperGroupIDListFromCache(ctx, userID); err != nil {
				trace_log.SetCtxInfo(ctx, "DelJoinedSuperGroupIDListFromCache", err, "userID", userID)
			}
		}
		for _, v := range req.InvitedUserIDList {
			chat.SuperGroupNotification(req.OperationID, v, v)
		}
	}
	return resp, nil
}

func (s *groupServer) GetGroupAllMember(ctx context.Context, req *pbGroup.GetGroupAllMemberReq) (*pbGroup.GetGroupAllMemberResp, error) {
	resp := &pbGroup.GetGroupAllMemberResp{}

	groupInfo, err := rocksCache.GetGroupInfoFromCache(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if groupInfo.GroupType != constant.SuperGroup {
		memberList, err := rocksCache.GetGroupMembersInfoFromCache(ctx, req.Count, req.Offset, req.GroupID)
		if err != nil {
			return nil, err
		}
		for _, v := range memberList {
			var node open_im_sdk.GroupMemberFullInfo
			cp.GroupMemberDBCopyOpenIM(&node, v)
			resp.MemberList = append(resp.MemberList, &node)
		}
	}
	return resp, nil
}

func (s *groupServer) GetGroupMemberList(ctx context.Context, req *pbGroup.GetGroupMemberListReq) (*pbGroup.GetGroupMemberListResp, error) {
	resp := &pbGroup.GetGroupMemberListResp{}

	memberList, err := imdb.GetGroupMemberByGroupID(req.GroupID, req.Filter, req.NextSeq, 30)
	if err != nil {
		return nil, err
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
	return resp, nil
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
	resp := &pbGroup.KickGroupMemberResp{}

	groupInfo, err := rocksCache.GetGroupInfoFromCache(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	var okUserIDList []string
	if groupInfo.GroupType != constant.SuperGroup {
		opFlag := 0
		if !token_verify.IsManagerUserID(req.OpUserID) {
			opInfo, err := rocksCache.GetGroupMemberInfoFromCache(ctx, req.GroupID, req.OpUserID)
			if err != nil {
				constant.SetErrorForResp(err, resp.CommonResp)
				return nil, err
			}
			if opInfo.RoleLevel == constant.GroupOrdinaryUsers {
				return nil, utils.Wrap(constant.ErrNoPermission, "")
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
			//log.NewError(req.OperationID, "failed, kick list 0")
			//return &pbGroup.KickGroupMemberResp{ErrCode: constant.ErrArgs.ErrCode, ErrMsg: constant.ErrArgs.ErrMsg}, nil
			return nil, utils.Wrap(constant.ErrArgs, "")
		}
		if err := s.DelGroupAndUserCache(ctx, req.GroupID, req.KickedUserIDList); err != nil {
			return nil, err
		}
		//remove
		for _, v := range req.KickedUserIDList {
			kickedInfo, err := rocksCache.GetGroupMemberInfoFromCache(ctx, req.GroupID, v)
			if err != nil {
				resp.Id2ResultList = append(resp.Id2ResultList, &pbGroup.Id2Result{UserID: v, Result: -1})
				trace_log.SetCtxInfo(ctx, "GetGroupMemberInfoFromCache", err, "groupID", req.GroupID, "userID", v)
				continue
			}

			if kickedInfo.RoleLevel == constant.GroupAdmin && opFlag == 3 {
				resp.Id2ResultList = append(resp.Id2ResultList, &pbGroup.Id2Result{UserID: v, Result: -1})
				trace_log.SetCtxInfo(ctx, "", nil, "msg", "is constant.GroupAdmin, can't kicked", "groupID", req.GroupID, "userID", v)
				continue
			}
			if kickedInfo.RoleLevel == constant.GroupOwner && opFlag != 1 {
				resp.Id2ResultList = append(resp.Id2ResultList, &pbGroup.Id2Result{UserID: v, Result: -1})
				trace_log.SetCtxInfo(ctx, "", nil, "msg", "is constant.GroupOwner, can't kicked", "groupID", req.GroupID, "userID", v)
				continue
			}

			err = imdb.DeleteGroupMemberByGroupIDAndUserID(req.GroupID, v)
			trace_log.SetCtxInfo(ctx, "RemoveGroupMember", err, "groupID", req.GroupID, "userID", v)
			if err != nil {
				log.NewError(req.OperationID, "RemoveGroupMember failed ", err.Error(), req.GroupID, v)
				resp.Id2ResultList = append(resp.Id2ResultList, &pbGroup.Id2Result{UserID: v, Result: -1})
			} else {
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
			etcdConn, err := getcdv3.GetConn(ctx, config.Config.RpcRegisterName.OpenImUserName)
			if err != nil {
				return nil, err
			}
			client := pbUser.NewUserClient(etcdConn)
			respPb, err := client.SetConversation(context.Background(), &reqPb)
			trace_log.SetCtxInfo(ctx, "SetConversation", err, "req", &reqPb, "resp", respPb)
		}
	} else {
		okUserIDList = req.KickedUserIDList
		if err := db.DB.RemoverUserFromSuperGroup(req.GroupID, okUserIDList); err != nil {
			return nil, err
		}
	}

	if groupInfo.GroupType != constant.SuperGroup {
		for _, userID := range okUserIDList {
			if err := rocksCache.DelGroupMemberInfoFromCache(ctx, req.GroupID, userID); err != nil {
				trace_log.SetCtxInfo(ctx, "DelGroupMemberInfoFromCache", err, "groupID", req.GroupID, "userID", userID)
			}
		}
		chat.MemberKickedNotification(req, okUserIDList)
	} else {
		for _, userID := range okUserIDList {
			if err = rocksCache.DelJoinedSuperGroupIDListFromCache(ctx, userID); err != nil {
				trace_log.SetCtxInfo(ctx, "DelGroupMemberInfoFromCache", err, "userID", userID)
			}
		}
		go func() {
			for _, v := range req.KickedUserIDList {
				chat.SuperGroupNotification(req.OperationID, v, v)
			}
		}()

	}
	return resp, nil
}

func (s *groupServer) GetGroupMembersInfo(ctx context.Context, req *pbGroup.GetGroupMembersInfoReq) (*pbGroup.GetGroupMembersInfoResp, error) {
	resp := &pbGroup.GetGroupMembersInfoResp{}
	resp.MemberList = []*open_im_sdk.GroupMemberFullInfo{}
	for _, userID := range req.MemberList {
		groupMember, err := rocksCache.GetGroupMemberInfoFromCache(ctx, req.GroupID, userID)
		if err != nil {
			return nil, err
		}
		var memberNode open_im_sdk.GroupMemberFullInfo
		utils.CopyStructFields(&memberNode, groupMember)
		memberNode.JoinTime = int32(groupMember.JoinTime.Unix())
		resp.MemberList = append(resp.MemberList, &memberNode)
	}
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

func (s *groupServer) GetGroupApplicationList(ctx context.Context, req *pbGroup.GetGroupApplicationListReq) (*pbGroup.GetGroupApplicationListResp, error) {
	resp := &pbGroup.GetGroupApplicationListResp{}
	reply, err := imdb.GetRecvGroupApplicationList(req.FromUserID)
	if err != nil {
		return nil, err
	}
	var errResult error
	trace_log.SetCtxInfo(ctx, "GetRecvGroupApplicationList", nil, " FromUserID: ", req.FromUserID, "GroupApplicationList: ", reply)
	for _, v := range reply {
		node := open_im_sdk.GroupRequest{UserInfo: &open_im_sdk.PublicUserInfo{}, GroupInfo: &open_im_sdk.GroupInfo{}}
		err := FillGroupInfoByGroupID(req.OperationID, v.GroupID, node.GroupInfo)
		if err != nil {
			if !errors.Is(errors.Unwrap(err), constant.ErrDismissedAlready) {
				errResult = err
			}
			continue
		}
		trace_log.SetCtxInfo(ctx, "FillGroupInfoByGroupID ", nil, " groupID: ", v.GroupID, " groupInfo: ", node.GroupInfo)
		err = FillPublicUserInfoByUserID(req.OperationID, v.UserID, node.UserInfo)
		if err != nil {
			errResult = err
			continue
		}
		cp.GroupRequestDBCopyOpenIM(&node, &v)
		resp.GroupRequestList = append(resp.GroupRequestList, &node)
	}
	if errResult != nil && len(resp.GroupRequestList) == 0 {
		return nil, err
	}
	trace_log.SetRpcRespInfo(ctx, utils.GetSelfFuncName(), resp.String())
	return resp, nil
}

func (s *groupServer) GetGroupsInfo(ctx context.Context, req *pbGroup.GetGroupsInfoReq) (*pbGroup.GetGroupsInfoResp, error) {
	resp := &pbGroup.GetGroupsInfoResp{}
	groupsInfoList := make([]*open_im_sdk.GroupInfo, 0)
	for _, groupID := range req.GroupIDList {
		groupInfoFromRedis, err := rocksCache.GetGroupInfoFromCache(ctx, groupID)
		if err != nil {
			constant.SetErrorForResp(err, resp.CommonResp)
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
		trace_log.SetCtxInfo(ctx, utils.GetSelfFuncName(), err, "groupID", groupID, "userID", userID)
	}()
	if !token_verify.IsManagerUserID(userID) && !imdb.IsGroupOwnerAdmin(groupID, userID) {
		return utils.Wrap(constant.ErrNoPermission, utils.GetSelfFuncName())
	}
	return nil
}

func (s *groupServer) GroupApplicationResponse(ctx context.Context, req *pbGroup.GroupApplicationResponseReq) (*pbGroup.GroupApplicationResponseResp, error) {
	resp := &pbGroup.GroupApplicationResponseResp{}

	if err := CheckPermission(ctx, req.GroupID, req.OpUserID); err != nil {
		return nil, err
	}
	groupRequest := imdb.GroupRequest{}
	utils.CopyStructFields(&groupRequest, req)
	groupRequest.UserID = req.FromUserID
	groupRequest.HandleUserID = req.OpUserID
	groupRequest.HandledTime = time.Now()
	if err := (&imdb.GroupRequest{}).Update(ctx, []*imdb.GroupRequest{&groupRequest}); err != nil {
		return nil, err
	}
	groupInfo, err := rocksCache.GetGroupInfoFromCache(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if req.HandleResult == constant.GroupResponseAgree {
		user, err := imdb.GetUserByUserID(req.FromUserID)
		if err != nil {
			return nil, err
		}
		request, err := (&imdb.GroupRequest{}).Take(ctx, req.GroupID, req.FromUserID)
		if err != nil {
			return nil, err
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
			return nil, err
		}

		err = (&imdb.GroupMember{}).Create(ctx, []*imdb.GroupMember{&member})
		if err != nil {
			return nil, err
		}
		etcdCacheConn, err := fault_tolerant.GetDefaultConn(config.Config.RpcRegisterName.OpenImCacheName, req.OperationID)
		if err != nil {
			return nil, err
		}
		cacheClient := pbCache.NewCacheClient(etcdCacheConn)
		cacheResp, err := cacheClient.DelGroupMemberIDListFromCache(context.Background(), &pbCache.DelGroupMemberIDListFromCacheReq{OperationID: req.OperationID, GroupID: req.GroupID})
		if err != nil {
			return nil, err
		}
		if cacheResp.CommonResp.ErrCode != 0 {
			return nil, utils.Wrap(&constant.ErrInfo{
				ErrCode: cacheResp.CommonResp.ErrCode,
				ErrMsg:  cacheResp.CommonResp.ErrMsg,
			}, "")
		}
		_ = rocksCache.DelGroupMemberListHashFromCache(ctx, req.GroupID)
		_ = rocksCache.DelJoinedGroupIDListFromCache(ctx, req.FromUserID)
		_ = rocksCache.DelGroupMemberNumFromCache(ctx, req.GroupID)
		chat.GroupApplicationAcceptedNotification(req)
		chat.MemberEnterNotification(req)
	} else if req.HandleResult == constant.GroupResponseRefuse {
		chat.GroupApplicationRejectedNotification(req)
	} else {
		//return nil, utils.Wrap(constant.ErrArgs, "")
		return nil, constant.ErrArgs.Warp()
	}
	return resp, nil
}

func (s *groupServer) JoinGroup(ctx context.Context, req *pbGroup.JoinGroupReq) (*pbGroup.JoinGroupResp, error) {
	resp := &pbGroup.JoinGroupResp{}

	if _, err := imdb.GetUserByUserID(req.OpUserID); err != nil {
		return nil, err
	}
	groupInfo, err := rocksCache.GetGroupInfoFromCache(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if groupInfo.Status == constant.GroupStatusDismissed {
		constant.SetErrorForResp(constant.ErrDismissedAlready, resp.CommonResp)
		return nil, utils.Wrap(constant.ErrDismissedAlready, "")
	}

	if groupInfo.NeedVerification == constant.Directly {
		if groupInfo.GroupType != constant.SuperGroup {
			us, err := imdb.GetUserByUserID(req.OpUserID)
			if err != nil {
				return nil, err
			}
			//to group member
			groupMember := imdb.GroupMember{GroupID: req.GroupID, RoleLevel: constant.GroupOrdinaryUsers, OperatorUserID: req.OpUserID}
			utils.CopyStructFields(&groupMember, us)
			if err := CallbackBeforeMemberJoinGroup(ctx, req.OperationID, &groupMember, groupInfo.Ex); err != nil {
				return nil, err
			}
			if err := s.DelGroupAndUserCache(ctx, req.GroupID, []string{req.OpUserID}); err != nil {
				return nil, err
			}
			err = imdb.InsertIntoGroupMember(groupMember)
			if err != nil {
				return nil, err
			}

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
			etcdConn, err := getcdv3.GetConn(ctx, config.Config.RpcRegisterName.OpenImUserName)
			if err != nil {
				return nil, err
			}
			client := pbUser.NewUserClient(etcdConn)
			respPb, err := client.SetConversation(context.Background(), &reqPb)
			trace_log.SetCtxInfo(ctx, "SetConversation", err, "req", reqPb, "resp", respPb)
			chat.MemberEnterDirectlyNotification(req.GroupID, req.OpUserID, req.OperationID)
			return resp, nil
		} else {
			constant.SetErrorForResp(constant.ErrGroupTypeNotSupport, resp.CommonResp)
			return resp, nil
		}
	}
	var groupRequest imdb.GroupRequest
	groupRequest.UserID = req.OpUserID
	groupRequest.ReqMsg = req.ReqMessage
	groupRequest.GroupID = req.GroupID
	groupRequest.JoinSource = req.JoinSource
	err = imdb.InsertIntoGroupRequest(groupRequest)
	if err != nil {
		return nil, err
	}
	chat.JoinGroupApplicationNotification(req)
	return resp, nil
}

func (s *groupServer) QuitGroup(ctx context.Context, req *pbGroup.QuitGroupReq) (*pbGroup.QuitGroupResp, error) {
	resp := &pbGroup.QuitGroupResp{}

	groupInfo, err := imdb.GetGroupInfoByGroupID(req.GroupID)
	if err != nil {
		return nil, err
	}
	if groupInfo.GroupType != constant.SuperGroup {
		_, err = rocksCache.GetGroupMemberInfoFromCache(ctx, req.GroupID, req.OpUserID)
		if err != nil {
			return nil, err
		}
		if err := s.DelGroupAndUserCache(ctx, req.GroupID, []string{req.OpUserID}); err != nil {
			return nil, err
		}
		err = imdb.DeleteGroupMemberByGroupIDAndUserID(req.GroupID, req.OpUserID)
		if err != nil {
			return nil, err
		}
	} else {
		okUserIDList := []string{req.OpUserID}
		if err := db.DB.RemoverUserFromSuperGroup(req.GroupID, okUserIDList); err != nil {
			return nil, err
		}
	}

	if groupInfo.GroupType != constant.SuperGroup {
		_ = rocksCache.DelGroupMemberInfoFromCache(ctx, req.GroupID, req.OpUserID)
		chat.MemberQuitNotification(req)
	} else {
		_ = rocksCache.DelJoinedSuperGroupIDListFromCache(ctx, req.OpUserID)
		_ = rocksCache.DelGroupMemberListHashFromCache(ctx, req.GroupID)
		chat.SuperGroupNotification(req.OperationID, req.OpUserID, req.OpUserID)
	}
	return resp, nil
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
	resp := &pbGroup.SetGroupInfoResp{}

	if !hasAccess(req) {
		return nil, utils.Wrap(constant.ErrIdentity, "")
	}
	group, err := imdb.GetGroupInfoByGroupID(req.GroupInfoForSet.GroupID)
	if err != nil {
		return nil, err
	}
	if group.Status == constant.GroupStatusDismissed {
		return nil, utils.Wrap(constant.ErrDismissedAlready, "")
	}

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
			return nil, err
		}
	}
	if req.GroupInfoForSet.LookMemberInfo != nil {
		changedType = changedType | (1 << 5)
		m := make(map[string]interface{})
		m["look_member_info"] = req.GroupInfoForSet.LookMemberInfo.Value
		if err := imdb.UpdateGroupInfoDefaultZero(req.GroupInfoForSet.GroupID, m); err != nil {
			return nil, err
		}
	}
	if req.GroupInfoForSet.ApplyMemberFriend != nil {
		changedType = changedType | (1 << 6)
		m := make(map[string]interface{})
		m["apply_member_friend"] = req.GroupInfoForSet.ApplyMemberFriend.Value
		if err := imdb.UpdateGroupInfoDefaultZero(req.GroupInfoForSet.GroupID, m); err != nil {
			return nil, err
		}
	}
	//only administrators can set group information
	var groupInfo imdb.Group
	utils.CopyStructFields(&groupInfo, req.GroupInfoForSet)
	if req.GroupInfoForSet.Notification != "" {
		groupInfo.NotificationUserID = req.OpUserID
		groupInfo.NotificationUpdateTime = time.Now()
	}
	if err := rocksCache.DelGroupInfoFromCache(ctx, req.GroupInfoForSet.GroupID); err != nil {
		return nil, err
	}
	err = imdb.SetGroupInfo(groupInfo)
	if err != nil {
		return nil, err
	}
	if changedType != 0 {
		chat.GroupInfoSetNotification(req.OperationID, req.OpUserID, req.GroupInfoForSet.GroupID, groupName, notification, introduction, faceURL, req.GroupInfoForSet.NeedVerification)
	}
	if req.GroupInfoForSet.Notification != "" {
		//get group member user id
		getGroupMemberIDListFromCacheReq := &pbCache.GetGroupMemberIDListFromCacheReq{OperationID: req.OperationID, GroupID: req.GroupInfoForSet.GroupID}
		etcdConn, err := getcdv3.GetConn(ctx, config.Config.RpcRegisterName.OpenImCacheName)
		if err != nil {
			return nil, err
		}
		client := pbCache.NewCacheClient(etcdConn)
		cacheResp, err := client.GetGroupMemberIDListFromCache(ctx, getGroupMemberIDListFromCacheReq)
		if err != nil {
			return nil, err
		}
		if err = constant.CommonResp2Err(cacheResp.CommonResp); err != nil {
			return nil, err
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
		nEtcdConn, err := getcdv3.GetConn(ctx, config.Config.RpcRegisterName.OpenImConversationName)
		if err != nil {
			return nil, err
		}
		nClient := pbConversation.NewConversationClient(nEtcdConn)
		conversationReply, err := nClient.ModifyConversationField(context.Background(), &conversationReq)
		trace_log.SetCtxInfo(ctx, "ModifyConversationField", err, "req", &conversationReq, "resp", conversationReply)
	}
	return resp, nil
}

func (s *groupServer) TransferGroupOwner(ctx context.Context, req *pbGroup.TransferGroupOwnerReq) (*pbGroup.TransferGroupOwnerResp, error) {
	resp := &pbGroup.TransferGroupOwnerResp{}

	groupInfo, err := imdb.GetGroupInfoByGroupID(req.GroupID)
	if err != nil {
		return nil, err
	}
	if groupInfo.Status == constant.GroupStatusDismissed {
		return nil, utils.Wrap(constant.ErrDismissedAlready, "")
	}

	if req.OldOwnerUserID == req.NewOwnerUserID {
		return nil, err
	}
	err = rocksCache.DelGroupMemberInfoFromCache(ctx, req.GroupID, req.NewOwnerUserID)
	if err != nil {
		return nil, err
	}
	err = rocksCache.DelGroupMemberInfoFromCache(ctx, req.GroupID, req.OldOwnerUserID)
	if err != nil {
		return nil, err
	}

	groupMemberInfo := imdb.GroupMember{GroupID: req.GroupID, UserID: req.OldOwnerUserID, RoleLevel: constant.GroupOrdinaryUsers}
	err = imdb.UpdateGroupMemberInfo(groupMemberInfo)
	if err != nil {
		return nil, err
	}
	groupMemberInfo = imdb.GroupMember{GroupID: req.GroupID, UserID: req.NewOwnerUserID, RoleLevel: constant.GroupOwner}
	err = imdb.UpdateGroupMemberInfo(groupMemberInfo)
	if err != nil {
		return nil, err
	}
	chat.GroupOwnerTransferredNotification(req)
	return resp, nil
}

func (s *groupServer) GetGroups(ctx context.Context, req *pbGroup.GetGroupsReq) (*pbGroup.GetGroupsResp, error) {
	resp := &pbGroup.GetGroupsResp{
		Groups:     []*pbGroup.CMSGroup{},
		Pagination: &open_im_sdk.ResponsePagination{CurrentPage: req.Pagination.PageNumber, ShowNumber: req.Pagination.ShowNumber},
	}

	if req.GroupID != "" {
		groupInfoDB, err := imdb.GetGroupInfoByGroupID(req.GroupID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return resp, nil
			}
			return nil, err
		}
		resp.GroupNum = 1
		groupInfo := &open_im_sdk.GroupInfo{}
		utils.CopyStructFields(groupInfo, groupInfoDB)
		groupMember, err := imdb.GetGroupOwnerInfoByGroupID(req.GroupID)
		if err != nil {
			return nil, err
		}
		memberNum, err := imdb.GetGroupMembersCount(req.GroupID, "")
		if err != nil {
			return nil, err
		}
		groupInfo.MemberCount = uint32(memberNum)
		groupInfo.CreateTime = uint32(groupInfoDB.CreateTime.Unix())
		resp.Groups = append(resp.Groups, &pbGroup.CMSGroup{GroupInfo: groupInfo, GroupOwnerUserName: groupMember.Nickname, GroupOwnerUserID: groupMember.UserID})
	} else {
		groups, count, err := imdb.GetGroupsByName(req.GroupName, req.Pagination.PageNumber, req.Pagination.ShowNumber)
		if err != nil {
			trace_log.SetCtxInfo(ctx, "GetGroupsByName", err, "GroupName", req.GroupName, "PageNumber", req.Pagination.PageNumber, "ShowNumber", req.Pagination.ShowNumber)
		}
		for _, v := range groups {
			group := &pbGroup.CMSGroup{GroupInfo: &open_im_sdk.GroupInfo{}}
			utils.CopyStructFields(group.GroupInfo, v)
			groupMember, err := imdb.GetGroupOwnerInfoByGroupID(v.GroupID)
			if err != nil {
				trace_log.SetCtxInfo(ctx, "GetGroupOwnerInfoByGroupID", err, "GroupID", v.GroupID)
				continue
			}
			group.GroupInfo.CreateTime = uint32(v.CreateTime.Unix())
			group.GroupOwnerUserID = groupMember.UserID
			group.GroupOwnerUserName = groupMember.Nickname
			resp.Groups = append(resp.Groups, group)
		}
		resp.GroupNum = int32(count)
	}
	return resp, nil
}

func (s *groupServer) GetGroupMembersCMS(ctx context.Context, req *pbGroup.GetGroupMembersCMSReq) (*pbGroup.GetGroupMembersCMSResp, error) {
	resp := &pbGroup.GetGroupMembersCMSResp{}
	groupMembers, err := imdb.GetGroupMembersByGroupIdCMS(req.GroupID, req.UserName, req.Pagination.ShowNumber, req.Pagination.PageNumber)
	if err != nil {
		return nil, err
	}
	groupMembersCount, err := imdb.GetGroupMembersCount(req.GroupID, req.UserName)
	if err != nil {
		return nil, err
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
	return resp, nil
}

func (s *groupServer) GetUserReqApplicationList(ctx context.Context, req *pbGroup.GetUserReqApplicationListReq) (*pbGroup.GetUserReqApplicationListResp, error) {
	resp := &pbGroup.GetUserReqApplicationListResp{}
	groupRequests, err := imdb.GetUserReqGroupByUserID(req.UserID)
	if err != nil {
		return nil, err
	}
	for _, groupReq := range groupRequests {
		node := open_im_sdk.GroupRequest{UserInfo: &open_im_sdk.PublicUserInfo{}, GroupInfo: &open_im_sdk.GroupInfo{}}
		group, err := imdb.GetGroupInfoByGroupID(groupReq.GroupID)
		if err != nil {
			trace_log.SetCtxInfo(ctx, "GetGroupInfoByGroupID", err, "GroupID", groupReq.GroupID)
			continue
		}
		user, err := imdb.GetUserByUserID(groupReq.UserID)
		if err != nil {
			trace_log.SetCtxInfo(ctx, "GetUserByUserID", err, "UserID", groupReq.UserID)
			continue
		}
		cp.GroupRequestDBCopyOpenIM(&node, &groupReq)
		cp.UserDBCopyOpenIMPublicUser(node.UserInfo, user)
		cp.GroupDBCopyOpenIM(node.GroupInfo, group)
		resp.GroupRequestList = append(resp.GroupRequestList, &node)
	}
	return resp, nil
}

func (s *groupServer) DismissGroup(ctx context.Context, req *pbGroup.DismissGroupReq) (*pbGroup.DismissGroupResp, error) {
	resp := &pbGroup.DismissGroupResp{}

	if !token_verify.IsManagerUserID(req.OpUserID) && !imdb.IsGroupOwnerAdmin(req.GroupID, req.OpUserID) {
		return nil, utils.Wrap(constant.ErrIdentity, "")
	}

	if err := rocksCache.DelGroupInfoFromCache(ctx, req.GroupID); err != nil {
		return nil, err
	}
	if err := s.DelGroupAndUserCache(ctx, req.GroupID, nil); err != nil {
		return nil, err
	}

	err := imdb.OperateGroupStatus(req.GroupID, constant.GroupStatusDismissed)
	if err != nil {
		return nil, err
	}
	groupInfo, err := imdb.GetGroupInfoByGroupID(req.GroupID)
	if err != nil {
		return nil, err
	}
	if groupInfo.GroupType != constant.SuperGroup {
		memberList, err := imdb.GetGroupMemberListByGroupID(req.GroupID)
		if err != nil {
			trace_log.SetCtxInfo(ctx, "GetGroupMemberListByGroupID", err, "groupID", req.GroupID)
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
			etcdConn, err := getcdv3.GetConn(ctx, config.Config.RpcRegisterName.OpenImUserName)
			client := pbUser.NewUserClient(etcdConn)
			respPb, err := client.SetConversation(context.Background(), &reqPb)
			trace_log.SetCtxInfo(ctx, "SetConversation", err, "req", &reqPb, "resp", respPb)
		}
		err = imdb.DeleteGroupMemberByGroupID(req.GroupID)
		if err != nil {
			return nil, err
		}
		chat.GroupDismissedNotification(req)
	} else {
		err = db.DB.DeleteSuperGroup(req.GroupID)
		if err != nil {
			return nil, err
		}
	}
	return resp, nil
}

func (s *groupServer) MuteGroupMember(ctx context.Context, req *pbGroup.MuteGroupMemberReq) (*pbGroup.MuteGroupMemberResp, error) {
	resp := &pbGroup.MuteGroupMemberResp{}

	opFlag, err := s.getGroupUserLevel(req.GroupID, req.OpUserID)
	if err != nil {
		return nil, err
	}
	if opFlag == 0 {
		return nil, err
	}

	mutedInfo, err := rocksCache.GetGroupMemberInfoFromCache(ctx, req.GroupID, req.UserID)
	if err != nil {
		return nil, err
	}
	if mutedInfo.RoleLevel == constant.GroupOwner && opFlag != 1 {
		return nil, err
	}
	if mutedInfo.RoleLevel == constant.GroupAdmin && opFlag == 3 {
		return nil, err
	}

	if err := rocksCache.DelGroupMemberInfoFromCache(ctx, req.GroupID, req.UserID); err != nil {
		return nil, err
	}
	groupMemberInfo := imdb.GroupMember{GroupID: req.GroupID, UserID: req.UserID}
	groupMemberInfo.MuteEndTime = time.Unix(int64(time.Now().Second())+int64(req.MutedSeconds), time.Now().UnixNano())
	err = imdb.UpdateGroupMemberInfo(groupMemberInfo)
	if err != nil {
		return nil, err
	}
	chat.GroupMemberMutedNotification(req.OperationID, req.OpUserID, req.GroupID, req.UserID, req.MutedSeconds)
	return resp, nil
}

func (s *groupServer) CancelMuteGroupMember(ctx context.Context, req *pbGroup.CancelMuteGroupMemberReq) (*pbGroup.CancelMuteGroupMemberResp, error) {
	resp := &pbGroup.CancelMuteGroupMemberResp{}

	opFlag, err := s.getGroupUserLevel(req.GroupID, req.OpUserID)
	if err != nil {
		return nil, err
	}
	if opFlag == 0 {
		return nil, err
	}

	mutedInfo, err := imdb.GetGroupMemberInfoByGroupIDAndUserID(req.GroupID, req.UserID)
	if err != nil {
		return nil, err
	}
	if mutedInfo.RoleLevel == constant.GroupOwner && opFlag != 1 {
		return nil, err
	}
	if mutedInfo.RoleLevel == constant.GroupAdmin && opFlag == 3 {
		return nil, err
	}
	if err := rocksCache.DelGroupMemberInfoFromCache(ctx, req.GroupID, req.UserID); err != nil {
		return nil, err
	}

	groupMemberInfo := imdb.GroupMember{GroupID: req.GroupID, UserID: req.UserID}
	groupMemberInfo.MuteEndTime = time.Unix(0, 0)
	err = imdb.UpdateGroupMemberInfo(groupMemberInfo)
	if err != nil {
		return nil, err
	}
	chat.GroupMemberCancelMutedNotification(req.OperationID, req.OpUserID, req.GroupID, req.UserID)
	return resp, nil
}

func (s *groupServer) MuteGroup(ctx context.Context, req *pbGroup.MuteGroupReq) (*pbGroup.MuteGroupResp, error) {
	resp := &pbGroup.MuteGroupResp{}

	opFlag, err := s.getGroupUserLevel(req.GroupID, req.OpUserID)
	if err != nil {
		return nil, err
	}
	if opFlag == 0 {
		//errMsg := req.OperationID + "opFlag == 0  " + req.GroupID + req.OpUserID
		//log.Error(req.OperationID, errMsg)
		//return &pbGroup.MuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
		return nil, utils.Wrap(constant.ErrNoPermission, "")
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
	if err := rocksCache.DelGroupInfoFromCache(ctx, req.GroupID); err != nil {
		return nil, err
	}

	err = imdb.OperateGroupStatus(req.GroupID, constant.GroupStatusMuted)
	if err != nil {
		return nil, err
	}

	chat.GroupMutedNotification(req.OperationID, req.OpUserID, req.GroupID)
	return resp, nil
}

func (s *groupServer) CancelMuteGroup(ctx context.Context, req *pbGroup.CancelMuteGroupReq) (*pbGroup.CancelMuteGroupResp, error) {
	resp := &pbGroup.CancelMuteGroupResp{}

	opFlag, err := s.getGroupUserLevel(req.GroupID, req.OpUserID)
	if err != nil {
		return nil, err
	}
	if opFlag == 0 {
		return nil, err
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
	if err := rocksCache.DelGroupInfoFromCache(ctx, req.GroupID); err != nil {
		return nil, err
	}
	err = imdb.UpdateGroupInfoDefaultZero(req.GroupID, map[string]interface{}{"status": constant.GroupOk})
	if err != nil {
		return nil, err
	}
	chat.GroupCancelMutedNotification(req.OperationID, req.OpUserID, req.GroupID)
	return resp, nil
}

func (s *groupServer) SetGroupMemberNickname(ctx context.Context, req *pbGroup.SetGroupMemberNicknameReq) (*pbGroup.SetGroupMemberNicknameResp, error) {
	resp := &pbGroup.SetGroupMemberNicknameResp{}

	if req.OpUserID != req.UserID && !token_verify.IsManagerUserID(req.OpUserID) {
		return nil, utils.Wrap(constant.ErrIdentity, "")
	}
	cbReq := &pbGroup.SetGroupMemberInfoReq{
		GroupID:     req.GroupID,
		UserID:      req.UserID,
		OperationID: req.OperationID,
		OpUserID:    req.OpUserID,
		Nickname:    &wrapperspb.StringValue{Value: req.Nickname},
	}
	if err := CallbackBeforeSetGroupMemberInfo(ctx, cbReq); err != nil {
		return nil, err
	}
	nickName := cbReq.Nickname.Value
	groupMemberInfo := imdb.GroupMember{}
	groupMemberInfo.UserID = req.UserID
	groupMemberInfo.GroupID = req.GroupID
	if nickName == "" {
		userNickname, err := imdb.GetUserNameByUserID(groupMemberInfo.UserID)
		if err != nil {
			return nil, err
		}
		groupMemberInfo.Nickname = userNickname
	} else {
		groupMemberInfo.Nickname = nickName
	}

	if err := rocksCache.DelGroupMemberInfoFromCache(ctx, req.GroupID, req.UserID); err != nil {
		return nil, err
	}

	if err := imdb.UpdateGroupMemberInfo(groupMemberInfo); err != nil {
		return nil, err
	}
	chat.GroupMemberInfoSetNotification(req.OperationID, req.OpUserID, req.GroupID, req.UserID)
	return resp, nil
}

func (s *groupServer) SetGroupMemberInfo(ctx context.Context, req *pbGroup.SetGroupMemberInfoReq) (*pbGroup.SetGroupMemberInfoResp, error) {
	resp := &pbGroup.SetGroupMemberInfoResp{}

	if err := rocksCache.DelGroupMemberInfoFromCache(ctx, req.GroupID, req.UserID); err != nil {
		return nil, err
	}
	if err := CallbackBeforeSetGroupMemberInfo(ctx, req); err != nil {
		return nil, err
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
	if err := imdb.UpdateGroupMemberInfoByMap(groupMember, m); err != nil {
		return nil, err
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
	return resp, nil
}

func (s *groupServer) GetGroupAbstractInfo(ctx context.Context, req *pbGroup.GetGroupAbstractInfoReq) (*pbGroup.GetGroupAbstractInfoResp, error) {
	resp := &pbGroup.GetGroupAbstractInfoResp{}

	hashCode, err := rocksCache.GetGroupMemberListHashFromCache(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	resp.GroupMemberListHash = hashCode
	num, err := rocksCache.GetGroupMemberNumFromCache(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	resp.GroupMemberNumber = int32(num)
	return resp, nil
}

func (s *groupServer) DelGroupAndUserCache(ctx context.Context, groupID string, userIDList []string) error {
	operationID := trace_log.GetOperationID(ctx)
	if groupID != "" {
		etcdConn, err := getcdv3.GetConn(ctx, config.Config.RpcRegisterName.OpenImCacheName)
		if err != nil {
			return err
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
		err = constant.CommonResp2Err(cacheResp.CommonResp)
		err = rocksCache.DelGroupMemberListHashFromCache(ctx, groupID)
		if err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), groupID, err.Error())
			return utils.Wrap(err, "")
		}
		err = rocksCache.DelGroupMemberNumFromCache(ctx, groupID)
		if err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), groupID)
			return utils.Wrap(err, "")
		}
	}
	if userIDList != nil {
		for _, userID := range userIDList {
			err := rocksCache.DelJoinedGroupIDListFromCache(ctx, userID)
			if err != nil {
				log.NewError(operationID, utils.GetSelfFuncName(), err.Error())
				return utils.Wrap(err, "")
			}
		}
	}
	return nil
}
