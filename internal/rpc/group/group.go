package group

import (
	chat "Open_IM/internal/rpc/msg"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/cache"
	"Open_IM/pkg/common/db/controller"
	"Open_IM/pkg/common/db/relation"
	relation2 "Open_IM/pkg/common/db/table/relation"
	"Open_IM/pkg/common/db/unrelation"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/middleware"
	promePkg "Open_IM/pkg/common/prometheus"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/common/tracelog"
	"fmt"
	"github.com/OpenIMSDK/getcdv3"

	pbConversation "Open_IM/pkg/proto/conversation"
	pbGroup "Open_IM/pkg/proto/group"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	pbUser "Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"
	"context"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"net"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type groupServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
	controller.GroupInterface

	etcdConn *getcdv3.EtcdConn
	//userRpc         pbUser.UserClient
	//conversationRpc pbConversation.ConversationClient
}

func NewGroupServer(port int) *groupServer {
	log.NewPrivateLog(constant.LogFileName)
	g := groupServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImGroupName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
	ttl := 10
	etcdClient, err := getcdv3.NewEtcdConn(config.Config.Etcd.EtcdSchema, strings.Join(g.etcdAddr, ","), config.Config.RpcRegisterIP, config.Config.Etcd.UserName, config.Config.Etcd.Password, port, ttl)
	if err != nil {
		panic("NewEtcdConn failed" + err.Error())
	}
	err = etcdClient.RegisterEtcd("", g.rpcRegisterName)
	if err != nil {
		panic("NewEtcdConn failed" + err.Error())
	}
	etcdClient.SetDefaultEtcdConfig(config.Config.RpcRegisterName.OpenImUserName, config.Config.RpcPort.OpenImUserPort)
	//conn := etcdClient.GetConn("", config.Config.RpcRegisterName.OpenImUserName)
	//g.userRpc = pbUser.NewUserClient(conn)

	etcdClient.SetDefaultEtcdConfig(config.Config.RpcRegisterName.OpenImConversationName, config.Config.RpcPort.OpenImConversationPort)
	//conn = etcdClient.GetConn("", config.Config.RpcRegisterName.OpenImConversationName)
	//g.conversationRpc = pbConversation.NewConversationClient(conn)

	//mysql init
	var mysql relation.Mysql
	var mongo unrelation.Mongo
	var groupModel relation2.GroupModel
	var redis cache.RedisClient
	err = mysql.InitConn().AutoMigrateModel(&groupModel)
	if err != nil {
		panic("db init err:" + err.Error())
	}
	if mysql.GormConn() != nil {
		groupModel.DB = mysql.GormConn()
	} else {
		panic("db init err:" + "conn is nil")
	}
	mongo.InitMongo()
	redis.InitRedis()
	mongo.CreateSuperGroupIndex()
	g.GroupInterface = controller.NewGroupController(groupModel.DB, redis.GetClient(), mongo.GetClient())
	return &g
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

func (s *groupServer) CreateGroup(ctx context.Context, req *pbGroup.CreateGroupReq) (*pbGroup.CreateGroupResp, error) {
	resp := &pbGroup.CreateGroupResp{GroupInfo: &open_im_sdk.GroupInfo{}}
	if err := token_verify.CheckAccessV3(ctx, req.OwnerUserID); err != nil {
		return nil, err
	}
	if req.OwnerUserID == "" {
		return nil, constant.ErrArgs.Wrap("no group owner")
	}
	userIDs := append(append(req.InitMembers, req.AdminUserIDs...), req.OwnerUserID)
	if utils.Duplicate(userIDs) {
		return nil, constant.ErrArgs.Wrap("group member repeated")
	}
	userMap, err := GetUserInfoMap(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	if ids := utils.Single(userIDs, utils.Keys(userMap)); len(ids) > 0 {
		return nil, constant.ErrUserIDNotFound.Wrap(strings.Join(ids, ","))
	}
	if err := callbackBeforeCreateGroup(ctx, req); err != nil {
		return nil, err
	}
	var group relation2.GroupModel
	var groupMembers []*relation2.GroupMemberModel
	utils.CopyStructFields(&group, req.GroupInfo)
	group.GroupID = genGroupID(ctx, req.GroupInfo.GroupID)
	if req.GroupInfo.GroupType == constant.SuperGroup {
		if err := s.GroupInterface.CreateSuperGroup(ctx, group.GroupID, userIDs); err != nil {
			return nil, err
		}
	} else {
		joinGroup := func(userID string, roleLevel int32) error {
			user := userMap[userID]
			groupMember := &relation2.GroupMemberModel{GroupID: group.GroupID, RoleLevel: roleLevel, OperatorUserID: tracelog.GetOpUserID(ctx), JoinSource: constant.JoinByInvitation, InviterUserID: tracelog.GetOpUserID(ctx)}
			utils.CopyStructFields(&groupMember, user)
			if err := CallbackBeforeMemberJoinGroup(ctx, tracelog.GetOperationID(ctx), groupMember, group.Ex); err != nil {
				return err
			}
			groupMembers = append(groupMembers, groupMember)
			return nil
		}
		if err := joinGroup(req.OwnerUserID, constant.GroupOwner); err != nil {
			return nil, err
		}
		for _, userID := range req.AdminUserIDs {
			if err := joinGroup(userID, constant.GroupAdmin); err != nil {
				return nil, err
			}
		}
		for _, userID := range req.InitMembers {
			if err := joinGroup(userID, constant.GroupOrdinaryUsers); err != nil {
				return nil, err
			}
		}
	}
	if err := s.GroupInterface.CreateGroup(ctx, []*relation2.GroupModel{&group}, groupMembers); err != nil {
		return nil, err
	}
	utils.CopyStructFields(resp.GroupInfo, group)
	resp.GroupInfo.MemberCount = uint32(len(userIDs))
	if req.GroupInfo.GroupType == constant.SuperGroup {
		go func() {
			for _, userID := range userIDs {
				chat.SuperGroupNotification(tracelog.GetOperationID(ctx), userID, userID)
			}
		}()
	} else {
		chat.GroupCreatedNotification(tracelog.GetOperationID(ctx), tracelog.GetOpUserID(ctx), group.GroupID, userIDs)
	}
	return resp, nil
}

func (s *groupServer) GetJoinedGroupList(ctx context.Context, req *pbGroup.GetJoinedGroupListReq) (*pbGroup.GetJoinedGroupListResp, error) {
	resp := &pbGroup.GetJoinedGroupListResp{}
	if err := token_verify.CheckAccessV3(ctx, req.FromUserID); err != nil {
		return nil, err
	}
	groups, err := s.GroupInterface.GetJoinedGroupList(ctx, req.FromUserID)
	if err != nil {
		return nil, err
	}
	if len(groups) == 0 {
		return resp, nil
	}
	groupIDs := utils.Slice(groups, func(e *relation2.GroupModel) string {
		return e.GroupID
	})
	groupMemberNum, err := s.GroupInterface.GetGroupMemberNum(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	groupOwnerUserID, err := s.GroupInterface.GetGroupOwnerUserID(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	for _, group := range groups {
		if group.Status == constant.GroupStatusDismissed || group.GroupType == constant.SuperGroup {
			continue
		}
		var groupNode open_im_sdk.GroupInfo
		utils.CopyStructFields(&groupNode, group)
		groupNode.MemberCount = uint32(groupMemberNum[group.GroupID])
		groupNode.OwnerUserID = groupOwnerUserID[group.GroupID]
		groupNode.CreateTime = group.CreateTime.UnixMilli()
		groupNode.NotificationUpdateTime = group.NotificationUpdateTime.UnixMilli()
		resp.Groups = append(resp.Groups, &groupNode)
	}
	resp.Total = int32(len(resp.Groups))
	return resp, nil
}

func (s *groupServer) InviteUserToGroup(ctx context.Context, req *pbGroup.InviteUserToGroupReq) (*pbGroup.InviteUserToGroupResp, error) {
	resp := &pbGroup.InviteUserToGroupResp{}
	if len(req.InvitedUserIDs) == 0 {
		return nil, constant.ErrArgs.Wrap("user empty")
	}
	if utils.Duplicate(req.InvitedUserIDs) {
		return nil, constant.ErrArgs.Wrap("userID duplicate")
	}
	group, err := s.GroupInterface.TakeGroupByID(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if group.Status == constant.GroupStatusDismissed {
		return nil, constant.ErrDismissedAlready.Wrap()
	}
	members, err := s.GroupInterface.GetGroupMemberList(ctx, group.GroupID)
	if err != nil {
		return nil, err
	}
	memberMap := utils.SliceToMap(members, func(e *relation2.GroupMemberModel) string {
		return e.UserID
	})
	if ids := utils.Single(req.InvitedUserIDs, utils.Keys(memberMap)); len(ids) > 0 {
		return nil, constant.ErrArgs.Wrap("user in group " + strings.Join(ids, ","))
	}
	userMap, err := GetUserInfoMap(ctx, req.InvitedUserIDs)
	if err != nil {
		return nil, err
	}
	if ids := utils.Single(req.InvitedUserIDs, utils.Keys(userMap)); len(ids) > 0 {
		return nil, constant.ErrArgs.Wrap("user not found " + strings.Join(ids, ","))
	}
	if group.NeedVerification == constant.AllNeedVerification {
		if !token_verify.IsAppManagerUid(ctx) {
			opUserID := tracelog.GetOpUserID(ctx)
			member, ok := memberMap[opUserID]
			if !ok {
				return nil, constant.ErrNoPermission.Wrap("not in group")
			}
			if !(member.RoleLevel == constant.GroupOwner || member.RoleLevel == constant.GroupAdmin) {
				var requests []*relation2.GroupRequestModel
				for _, userID := range req.InvitedUserIDs {
					requests = append(requests, &relation2.GroupRequestModel{
						UserID:        userID,
						GroupID:       req.GroupID,
						JoinSource:    constant.JoinByInvitation,
						InviterUserID: opUserID,
					})
				}
				if err := s.GroupInterface.CreateGroupRequest(ctx, requests); err != nil {
					return nil, err
				}
				for _, request := range requests {
					chat.JoinGroupApplicationNotification(ctx, &pbGroup.JoinGroupReq{
						GroupID:       request.GroupID,
						ReqMessage:    request.ReqMsg,
						JoinSource:    request.JoinSource,
						InviterUserID: request.InviterUserID,
					})
				}
				return resp, nil
			}
		}
	}
	if group.GroupType == constant.SuperGroup {
		if err := s.GroupInterface.AddUserToSuperGroup(ctx, req.GroupID, req.InvitedUserIDs); err != nil {
			return nil, err
		}
		for _, userID := range req.InvitedUserIDs {
			chat.SuperGroupNotification(tracelog.GetOperationID(ctx), userID, userID)
		}
	} else {
		opUserID := tracelog.GetOpUserID(ctx)
		var groupMembers []*relation2.GroupMemberModel
		for _, userID := range req.InvitedUserIDs {
			user := userMap[userID]
			var member relation2.GroupMemberModel
			utils.CopyStructFields(&member, user)
			member.GroupID = req.GroupID
			member.RoleLevel = constant.GroupOrdinaryUsers
			member.OperatorUserID = opUserID
			member.InviterUserID = opUserID
			member.JoinSource = constant.JoinByInvitation
			if err := CallbackBeforeMemberJoinGroup(ctx, tracelog.GetOperationID(ctx), &member, group.Ex); err != nil {
				return nil, err
			}
			groupMembers = append(groupMembers, &member)
		}
		if err := s.GroupInterface.CreateGroupMember(ctx, groupMembers); err != nil {
			return nil, err
		}
		chat.MemberInvitedNotification(tracelog.GetOperationID(ctx), req.GroupID, tracelog.GetOpUserID(ctx), req.Reason, req.InvitedUserIDs)
	}
	return resp, nil
}

func (s *groupServer) GetGroupAllMember(ctx context.Context, req *pbGroup.GetGroupAllMemberReq) (*pbGroup.GetGroupAllMemberResp, error) {
	resp := &pbGroup.GetGroupAllMemberResp{}
	group, err := s.GroupInterface.TakeGroupByID(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if group.GroupType != constant.SuperGroup {
		members, err := s.GroupInterface.GetGroupMemberList(ctx, req.GroupID)
		if err != nil {
			return nil, err
		}
		for _, member := range members {
			var node open_im_sdk.GroupMemberFullInfo
			utils.CopyStructFields(&node, member)
			resp.Members = append(resp.Members, &node)
		}
	}
	return resp, nil
}

func (s *groupServer) GetGroupMemberList(ctx context.Context, req *pbGroup.GetGroupMemberListReq) (*pbGroup.GetGroupMemberListResp, error) {
	resp := &pbGroup.GetGroupMemberListResp{}
	members, err := s.GroupInterface.GetGroupMemberFilterList(ctx, req.GroupID, req.Filter, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}
	for _, member := range members {
		var info open_im_sdk.GroupMemberFullInfo
		utils.CopyStructFields(&info, &member)
		resp.Members = append(resp.Members, &info)
	}
	return resp, nil
}

//func (s *groupServer) getGroupUserLevel(groupID, userID string) (int, error) {
//	opFlag := 0
//	if !token_verify.IsManagerUserID(userID) {
//		opInfo, err := relation.GetGroupMemberInfoByGroupIDAndUserID(groupID, userID)
//		if err != nil {
//			return opFlag, utils.Wrap(err, "")
//		}
//		if opInfo.RoleLevel == constant.GroupOrdinaryUsers {
//			opFlag = 0
//		} else if opInfo.RoleLevel == constant.GroupOwner {
//			opFlag = 2 // owner
//		} else {
//			opFlag = 3 // admin
//		}
//	} else {
//		opFlag = 1 // app manager
//	}
//	return opFlag, nil
//}

func (s *groupServer) KickGroupMember(ctx context.Context, req *pbGroup.KickGroupMemberReq) (*pbGroup.KickGroupMemberResp, error) {
	resp := &pbGroup.KickGroupMemberResp{}
	group, err := s.GroupInterface.TakeGroupByID(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if len(req.KickedUserIDs) == 0 {
		return nil, constant.ErrArgs.Wrap("KickedUserIDs empty")
	}
	if utils.IsDuplicateStringSlice(req.KickedUserIDs) {
		return nil, constant.ErrArgs.Wrap("KickedUserIDs duplicate")
	}
	opUserID := tracelog.GetOpUserID(ctx)
	if utils.IsContain(opUserID, req.KickedUserIDs) {
		return nil, constant.ErrArgs.Wrap("opUserID in KickedUserIDs")
	}
	if group.GroupType == constant.SuperGroup {
		if err := s.GroupInterface.DelSuperGroupMember(ctx, req.GroupID, req.KickedUserIDs); err != nil {
			return nil, err
		}
		go func() {
			for _, userID := range req.KickedUserIDs {
				chat.SuperGroupNotification(tracelog.GetOperationID(ctx), userID, userID)
			}
		}()
	} else {
		members, err := s.GroupInterface.FindGroupMembersByID(ctx, req.GroupID, append(req.KickedUserIDs, opUserID))
		if err != nil {
			return nil, err
		}
		memberMap := make(map[string]*relation2.GroupMemberModel)
		for i, member := range members {
			memberMap[member.UserID] = members[i]
		}
		for _, userID := range req.KickedUserIDs {
			if _, ok := memberMap[userID]; !ok {
				return nil, constant.ErrUserIDNotFound.Wrap(userID)
			}
		}
		if !token_verify.IsAppManagerUid(ctx) {
			member := memberMap[opUserID]
			if member == nil {
				return nil, constant.ErrNoPermission.Wrap(fmt.Sprintf("opUserID %s no in group", opUserID))
			}
			switch member.RoleLevel {
			case constant.GroupOwner:
			case constant.GroupAdmin:
				for _, member := range members {
					if member.UserID == opUserID {
						continue
					}
					if member.RoleLevel == constant.GroupOwner || member.RoleLevel == constant.GroupAdmin {
						return nil, constant.ErrNoPermission.Wrap("userID:" + member.UserID)
					}
				}
			default:
				return nil, constant.ErrNoPermission.Wrap("opUserID is OrdinaryUser")
			}
		}
		if err := s.GroupInterface.DelGroupMember(ctx, group.GroupID, req.KickedUserIDs); err != nil {
			return nil, err
		}
		chat.MemberKickedNotification(req, req.KickedUserIDs)
	}
	return resp, nil
}

func (s *groupServer) GetGroupMembersInfo(ctx context.Context, req *pbGroup.GetGroupMembersInfoReq) (*pbGroup.GetGroupMembersInfoResp, error) {
	resp := &pbGroup.GetGroupMembersInfoResp{}
	members, err := s.GroupInterface.GetGroupMemberListByUserID(ctx, req.GroupID, req.Members)
	if err != nil {
		return nil, err
	}
	for _, member := range members {
		var memberNode open_im_sdk.GroupMemberFullInfo
		utils.CopyStructFields(&memberNode, member)
		memberNode.JoinTime = member.JoinTime.UnixMilli()
		resp.Members = append(resp.Members, &memberNode)
	}
	return resp, nil
}

//func FillGroupInfoByGroupID(operationID, groupID string, groupInfo *open_im_sdk.GroupInfo) error {
//	group, err := relation.TakeGroupInfoByGroupID(groupID)
//	if err != nil {
//		log.Error(operationID, "TakeGroupInfoByGroupID failed ", err.Error(), groupID)
//		return utils.Wrap(err, "")
//	}
//	if group.Status == constant.GroupStatusDismissed {
//		log.Debug(operationID, " group constant.GroupStatusDismissed ", group.GroupID)
//		return utils.Wrap(constant.ErrDismissedAlready, "")
//	}
//	return utils.Wrap(cp.GroupDBCopyOpenIM(groupInfo, group), "")
//}

//func FillPublicUserInfoByUserID(operationID, userID string, userInfo *open_im_sdk.PublicUserInfo) error {
//	user, err := relation.TakeUserByUserID(userID)
//	if err != nil {
//		log.Error(operationID, "TakeUserByUserID failed ", err.Error(), userID)
//		return utils.Wrap(err, "")
//	}
//	cp.UserDBCopyOpenIMPublicUser(userInfo, user)
//	return nil
//}

func (s *groupServer) GetGroupApplicationList(ctx context.Context, req *pbGroup.GetGroupApplicationListReq) (*pbGroup.GetGroupApplicationListResp, error) {
	resp := &pbGroup.GetGroupApplicationListResp{}
	groupRequests, err := s.GroupInterface.GetGroupRecvApplicationList(ctx, req.FromUserID)
	if err != nil {
		return nil, err
	}
	if len(groupRequests) == 0 {
		return resp, nil
	}
	var (
		userIDs  []string
		groupIDs []string
	)
	for _, gr := range groupRequests {
		userIDs = append(userIDs, gr.UserID)
		groupIDs = append(groupIDs, gr.GroupID)
	}
	userIDs = utils.Distinct(userIDs)
	groupIDs = utils.Distinct(groupIDs)
	userMap, err := GetPublicUserInfoMap(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	if ids := utils.Single(utils.Keys(userMap), userIDs); len(ids) > 0 {
		return nil, constant.ErrUserIDNotFound.Wrap(strings.Join(ids, ","))
	}
	groups, err := s.GroupInterface.FindGroupsByID(ctx, utils.Distinct(groupIDs))
	if err != nil {
		return nil, err
	}
	groupMap := utils.SliceToMap(groups, func(e *relation2.GroupModel) string {
		return e.GroupID
	})
	if ids := utils.Single(utils.Keys(groupMap), groupIDs); len(ids) > 0 {
		return nil, constant.ErrGroupIDNotFound.Wrap(strings.Join(ids, ","))
	}
	groupMemberNumMap, err := s.GroupInterface.GetGroupMemberNum(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	groupOwnerUserIDMap, err := s.GroupInterface.GetGroupOwnerUserID(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	for _, gr := range groupRequests {
		groupRequest := open_im_sdk.GroupRequest{UserInfo: &open_im_sdk.PublicUserInfo{}}
		utils.CopyStructFields(&groupRequest, gr)
		groupRequest.UserInfo = userMap[gr.UserID]
		groupRequest.GroupInfo = DbToPbGroupInfo(groupMap[gr.GroupID], groupOwnerUserIDMap[gr.GroupID], uint32(groupMemberNumMap[gr.GroupID]))
		resp.GroupRequests = append(resp.GroupRequests, &groupRequest)
	}
	return resp, nil
}

func (s *groupServer) GetGroupsInfo(ctx context.Context, req *pbGroup.GetGroupsInfoReq) (*pbGroup.GetGroupsInfoResp, error) {
	resp := &pbGroup.GetGroupsInfoResp{}
	if len(req.GroupIDs) == 0 {
		return nil, constant.ErrArgs.Wrap("groupID is empty")
	}
	groups, err := s.GroupInterface.FindGroupsByID(ctx, req.GroupIDs)
	if err != nil {
		return nil, err
	}
	groupMemberNumMap, err := s.GroupInterface.GetGroupMemberNum(ctx, req.GroupIDs)
	if err != nil {
		return nil, err
	}
	groupOwnerUserIDMap, err := s.GroupInterface.GetGroupOwnerUserID(ctx, req.GroupIDs)
	if err != nil {
		return nil, err
	}
	resp.GroupInfos = utils.Slice(groups, func(e *relation2.GroupModel) *open_im_sdk.GroupInfo {
		return DbToPbGroupInfo(e, groupOwnerUserIDMap[e.GroupID], uint32(groupMemberNumMap[e.GroupID]))
	})
	return resp, nil
}

func (s *groupServer) GroupApplicationResponse(ctx context.Context, req *pbGroup.GroupApplicationResponseReq) (*pbGroup.GroupApplicationResponseResp, error) {
	resp := &pbGroup.GroupApplicationResponseResp{}
	if !utils.Contain(req.HandleResult, constant.GroupResponseAgree, constant.GroupResponseRefuse) {
		return nil, constant.ErrArgs.Wrap("HandleResult unknown")
	}
	if !token_verify.IsAppManagerUid(ctx) {
		groupMember, err := s.GroupInterface.TakeGroupMemberByID(ctx, req.GroupID, req.FromUserID)
		if err != nil {
			return nil, err
		}
		if !(groupMember.RoleLevel == constant.GroupOwner || groupMember.RoleLevel == constant.GroupAdmin) {
			return nil, constant.ErrNoPermission.Wrap("no group owner or admin")
		}
	}
	group, err := s.GroupInterface.TakeGroupByID(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	groupRequest, err := s.GroupInterface.TakeGroupRequest(ctx, req.GroupID, req.FromUserID)
	if err != nil {
		return nil, err
	}
	if groupRequest.HandleResult != 0 {
		return nil, constant.ErrArgs.Wrap("group request already processed")
	}
	if _, err := s.GroupInterface.TakeGroupMemberByID(ctx, req.GroupID, req.FromUserID); err != nil {
		if !IsNotFound(err) {
			return nil, err
		}
	} else {
		return nil, constant.ErrArgs.Wrap("already in group")
	}
	user, err := GetPublicUserInfoOne(ctx, req.FromUserID)
	if err != nil {
		return nil, err
	}
	var member *relation2.GroupMemberModel
	if req.HandleResult == constant.GroupResponseAgree {
		member = &relation2.GroupMemberModel{
			GroupID:        req.GroupID,
			UserID:         user.UserID,
			Nickname:       user.Nickname,
			FaceURL:        user.FaceURL,
			RoleLevel:      constant.GroupOrdinaryUsers,
			JoinTime:       time.Now(),
			JoinSource:     groupRequest.JoinSource,
			InviterUserID:  groupRequest.InviterUserID,
			OperatorUserID: tracelog.GetOpUserID(ctx),
			Ex:             groupRequest.Ex,
		}
		if err = CallbackBeforeMemberJoinGroup(ctx, tracelog.GetOperationID(ctx), member, group.Ex); err != nil {
			return nil, err
		}
	}
	if err := s.GroupInterface.HandlerGroupRequest(ctx, req.GroupID, req.FromUserID, req.HandledMsg, req.HandleResult, member); err != nil {
		return nil, err
	}
	if req.HandleResult == constant.GroupResponseAgree {
		chat.GroupApplicationAcceptedNotification(req)
		chat.MemberEnterNotification(req)
	} else if req.HandleResult == constant.GroupResponseRefuse {
		chat.GroupApplicationRejectedNotification(req)
	}
	return resp, nil
}

func (s *groupServer) JoinGroup(ctx context.Context, req *pbGroup.JoinGroupReq) (*pbGroup.JoinGroupResp, error) {
	resp := &pbGroup.JoinGroupResp{}
	if _, err := GetPublicUserInfoOne(ctx, tracelog.GetOpUserID(ctx)); err != nil {
		return nil, err
	}
	group, err := s.GroupInterface.TakeGroupByID(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if group.Status == constant.GroupStatusDismissed {
		return nil, constant.ErrDismissedAlready.Wrap()
	}
	if group.NeedVerification == constant.Directly {
		if group.GroupType == constant.SuperGroup {
			return nil, constant.ErrGroupTypeNotSupport.Wrap()
		}
		us, err := relation.GetUserByUserID(tracelog.GetOpUserID(ctx))
		if err != nil {
			return nil, err
		}
		groupMember := relation2.GroupMemberModel{GroupID: req.GroupID, RoleLevel: constant.GroupOrdinaryUsers, OperatorUserID: tracelog.GetOpUserID(ctx)}
		utils.CopyStructFields(&groupMember, us)
		if err := CallbackBeforeMemberJoinGroup(ctx, tracelog.GetOperationID(ctx), &groupMember, group.Ex); err != nil {
			return nil, err
		}
		if err := s.GroupInterface.CreateGroupMember(ctx, []*relation2.GroupMemberModel{&groupMember}); err != nil {
			return nil, err
		}
		chat.MemberEnterDirectlyNotification(req.GroupID, tracelog.GetOpUserID(ctx), tracelog.GetOperationID(ctx))
		return resp, nil
	}
	groupRequest := relation2.GroupRequestModel{
		UserID:     tracelog.GetOpUserID(ctx),
		ReqMsg:     req.ReqMessage,
		GroupID:    req.GroupID,
		JoinSource: req.JoinSource,
		ReqTime:    time.Now(),
	}
	if err := s.GroupInterface.CreateGroupRequest(ctx, []*relation2.GroupRequestModel{&groupRequest}); err != nil {
		return nil, err
	}
	chat.JoinGroupApplicationNotification(ctx, req)
	return resp, nil
}

func (s *groupServer) QuitGroup(ctx context.Context, req *pbGroup.QuitGroupReq) (*pbGroup.QuitGroupResp, error) {
	resp := &pbGroup.QuitGroupResp{}
	group, err := s.GroupInterface.TakeGroupByID(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if group.GroupType == constant.SuperGroup {
		if err := s.GroupInterface.DelSuperGroupMember(ctx, req.GroupID, []string{tracelog.GetOpUserID(ctx)}); err != nil {
			return nil, err
		}
		chat.SuperGroupNotification(tracelog.GetOperationID(ctx), tracelog.GetOpUserID(ctx), tracelog.GetOpUserID(ctx))
	} else {
		_, err := s.GroupInterface.TakeGroupMemberByID(ctx, req.GroupID, tracelog.GetOpUserID(ctx))
		if err != nil {
			return nil, err
		}
		chat.MemberQuitNotification(req)
	}
	return resp, nil
}

func (s *groupServer) SetGroupInfo(ctx context.Context, req *pbGroup.SetGroupInfoReq) (*pbGroup.SetGroupInfoResp, error) {
	resp := &pbGroup.SetGroupInfoResp{}
	if !token_verify.IsAppManagerUid(ctx) {
		groupMember, err := s.GroupInterface.TakeGroupMemberByID(ctx, req.GroupInfoForSet.GroupID, tracelog.GetOpUserID(ctx))
		if err != nil {
			return nil, err
		}
		if !(groupMember.RoleLevel == constant.GroupOwner || groupMember.RoleLevel == constant.GroupAdmin) {
			return nil, constant.ErrNoPermission.Wrap("no group owner or admin")
		}
	}
	group, err := s.TakeGroupByID(ctx, req.GroupInfoForSet.GroupID)
	if err != nil {
		return nil, err
	}
	if group.Status == constant.GroupStatusDismissed {
		return nil, utils.Wrap(constant.ErrDismissedAlready, "")
	}
	data := PbToDbMapGroupInfoForSet(req.GroupInfoForSet)
	if len(data) > 0 {
		return resp, nil
	}
	if err := s.GroupInterface.UpdateGroup(ctx, group.GroupID, data); err != nil {
		return nil, err
	}
	group, err = s.TakeGroupByID(ctx, req.GroupInfoForSet.GroupID)
	if err != nil {
		return nil, err
	}
	chat.GroupInfoSetNotification(tracelog.GetOperationID(ctx), tracelog.GetOpUserID(ctx), req.GroupInfoForSet.GroupID, group.GroupName, group.Notification, group.Introduction, group.FaceURL, req.GroupInfoForSet.NeedVerification)
	if req.GroupInfoForSet.Notification != "" {
		GroupNotification(ctx, group.GroupID)
	}
	return resp, nil
}

func (s *groupServer) TransferGroupOwner(ctx context.Context, req *pbGroup.TransferGroupOwnerReq) (*pbGroup.TransferGroupOwnerResp, error) {
	resp := &pbGroup.TransferGroupOwnerResp{}
	group, err := s.GroupInterface.TakeGroupByID(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if group.Status == constant.GroupStatusDismissed {
		return nil, utils.Wrap(constant.ErrDismissedAlready, "")
	}
	if req.OldOwnerUserID == req.NewOwnerUserID {
		return nil, constant.ErrArgs.Wrap("OldOwnerUserID == NewOwnerUserID")
	}
	members, err := s.GroupInterface.GetGroupMemberListByUserID(ctx, req.GroupID, []string{req.OldOwnerUserID, req.NewOwnerUserID})
	if err != nil {
		return nil, err
	}
	memberMap := utils.SliceToMap(members, func(e *relation2.GroupMemberModel) string { return e.UserID })
	if ids := utils.Single([]string{req.OldOwnerUserID, req.NewOwnerUserID}, utils.Keys(memberMap)); len(ids) > 0 {
		return nil, constant.ErrArgs.Wrap("user not in group " + strings.Join(ids, ","))
	}
	newOwner := memberMap[req.NewOwnerUserID]
	if newOwner == nil {
		return nil, constant.ErrArgs.Wrap("NewOwnerUser not in group " + req.NewOwnerUserID)
	}
	oldOwner := memberMap[req.OldOwnerUserID]
	if token_verify.IsAppManagerUid(ctx) {
		if oldOwner == nil {
			oldOwner, err = s.GroupInterface.GetGroupOwnerUser(ctx, req.OldOwnerUserID)
			if err != nil {
				return nil, err
			}
		}
	} else {
		if oldOwner == nil {
			return nil, constant.ErrArgs.Wrap("OldOwnerUser not in group " + req.NewOwnerUserID)
		}
		if oldOwner.GroupID != tracelog.GetOpUserID(ctx) {
			return nil, constant.ErrNoPermission.Wrap(fmt.Sprintf("user %s no permission transfer group owner", tracelog.GetOpUserID(ctx)))
		}
	}
	if err := s.GroupInterface.TransferGroupOwner(ctx, req.GroupID, req.OldOwnerUserID, req.NewOwnerUserID); err != nil {
		return nil, err
	}
	chat.GroupOwnerTransferredNotification(req)
	return resp, nil
}

func (s *groupServer) GetGroups(ctx context.Context, req *pbGroup.GetGroupsReq) (*pbGroup.GetGroupsResp, error) {
	resp := &pbGroup.GetGroupsResp{}
	var (
		groups []*relation2.GroupModel
		err    error
	)
	if req.GroupID != "" {
		groups, err = s.GroupInterface.FindGroupsByID(ctx, []string{req.GroupID})
		resp.GroupNum = int32(len(groups))
	} else {
		resp.GroupNum, groups, err = s.GroupInterface.FindSearchGroup(ctx, req.GroupName, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	}
	if err != nil {
		return nil, err
	}
	groupIDs := utils.Slice(groups, func(e *relation2.GroupModel) string {
		return e.GroupID
	})
	ownerMembers, err := s.GroupInterface.FindGroupOwnerUser(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	ownerMemberMap := utils.SliceToMap(ownerMembers, func(e *relation2.GroupMemberModel) string {
		return e.GroupID
	})
	if ids := utils.Single(groupIDs, utils.Keys(ownerMemberMap)); len(ids) > 0 {
		return nil, constant.ErrDB.Wrap("group not owner " + strings.Join(ids, ","))
	}
	groupMemberNumMap, err := s.GroupInterface.GetGroupMemberNum(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	for _, group := range groups {
		member := ownerMemberMap[group.GroupID]
		resp.Groups = append(resp.Groups, DbToBpCMSGroup(group, member.UserID, member.Nickname, uint32(groupMemberNumMap[group.GroupID])))
	}
	return resp, nil
}

func (s *groupServer) GetGroupMembersCMS(ctx context.Context, req *pbGroup.GetGroupMembersCMSReq) (*pbGroup.GetGroupMembersCMSResp, error) {
	resp := &pbGroup.GetGroupMembersCMSResp{}
	groupMembers, err := relation.GetGroupMembersByGroupIdCMS(req.GroupID, req.UserName, req.Pagination.ShowNumber, req.Pagination.PageNumber)
	if err != nil {
		return nil, err
	}
	groupMembersCount, err := relation.GetGroupMembersCount(req.GroupID, req.UserName)
	if err != nil {
		return nil, err
	}
	log.NewInfo(tracelog.GetOperationID(ctx), groupMembersCount)
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
	groupRequests, err := relation.GetUserReqGroupByUserID(req.UserID)
	if err != nil {
		return nil, err
	}
	for _, groupReq := range groupRequests {
		node := open_im_sdk.GroupRequest{UserInfo: &open_im_sdk.PublicUserInfo{}, GroupInfo: &open_im_sdk.GroupInfo{}}
		group, err := relation.GetGroupInfoByGroupID(groupReq.GroupID)
		if err != nil {
			tracelog.SetCtxInfo(ctx, "GetGroupInfoByGroupID", err, "GroupID", groupReq.GroupID)
			continue
		}
		user, err := relation.GetUserByUserID(groupReq.UserID)
		if err != nil {
			tracelog.SetCtxInfo(ctx, "GetUserByUserID", err, "UserID", groupReq.UserID)
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

	if !token_verify.IsManagerUserID(tracelog.GetOpUserID(ctx)) && !relation.IsGroupOwnerAdmin(req.GroupID, tracelog.GetOpUserID(ctx)) {
		return nil, utils.Wrap(constant.ErrIdentity, "")
	}

	if err := rocksCache.DelGroupInfoFromCache(ctx, req.GroupID); err != nil {
		return nil, err
	}
	if err := s.DelGroupAndUserCache(ctx, req.GroupID, nil); err != nil {
		return nil, err
	}

	err := relation.OperateGroupStatus(req.GroupID, constant.GroupStatusDismissed)
	if err != nil {
		return nil, err
	}
	groupInfo, err := relation.GetGroupInfoByGroupID(req.GroupID)
	if err != nil {
		return nil, err
	}
	if groupInfo.GroupType != constant.SuperGroup {
		memberList, err := relation.GetGroupMemberListByGroupID(req.GroupID)
		if err != nil {
			tracelog.SetCtxInfo(ctx, "GetGroupMemberListByGroupID", err, "groupID", req.GroupID)
		}
		//modify quitter conversation info
		var reqPb pbUser.SetConversationReq
		var c pbConversation.Conversation
		for _, v := range memberList {
			reqPb.OperationID = tracelog.GetOperationID(ctx)
			c.OwnerUserID = v.UserID
			c.ConversationID = utils.GetConversationIDBySessionType(req.GroupID, constant.GroupChatType)
			c.ConversationType = constant.GroupChatType
			c.GroupID = req.GroupID
			c.IsNotInGroup = true
			reqPb.Conversation = &c
			etcdConn, err := getcdv3.GetConn(ctx, config.Config.RpcRegisterName.OpenImUserName)
			client := pbUser.NewUserClient(etcdConn)
			respPb, err := client.SetConversation(context.Background(), &reqPb)
			tracelog.SetCtxInfo(ctx, "SetConversation", err, "req", &reqPb, "resp", respPb)
		}
		err = relation.DeleteGroupMemberByGroupID(req.GroupID)
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

	opFlag, err := s.getGroupUserLevel(req.GroupID, tracelog.GetOpUserID(ctx))
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
	groupMemberInfo := relation2.GroupMemberModel{GroupID: req.GroupID, UserID: req.UserID}
	groupMemberInfo.MuteEndTime = time.Unix(int64(time.Now().Second())+int64(req.MutedSeconds), time.Now().UnixNano())
	err = relation.UpdateGroupMemberInfo(groupMemberInfo)
	if err != nil {
		return nil, err
	}
	chat.GroupMemberMutedNotification(tracelog.GetOperationID(ctx), tracelog.GetOpUserID(ctx), req.GroupID, req.UserID, req.MutedSeconds)
	return resp, nil
}

func (s *groupServer) CancelMuteGroupMember(ctx context.Context, req *pbGroup.CancelMuteGroupMemberReq) (*pbGroup.CancelMuteGroupMemberResp, error) {
	resp := &pbGroup.CancelMuteGroupMemberResp{}

	opFlag, err := s.getGroupUserLevel(req.GroupID, tracelog.GetOpUserID(ctx))
	if err != nil {
		return nil, err
	}
	if opFlag == 0 {
		return nil, err
	}

	mutedInfo, err := relation.GetGroupMemberInfoByGroupIDAndUserID(req.GroupID, req.UserID)
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

	groupMemberInfo := relation2.GroupMemberModel{GroupID: req.GroupID, UserID: req.UserID}
	groupMemberInfo.MuteEndTime = time.Unix(0, 0)
	err = relation.UpdateGroupMemberInfo(groupMemberInfo)
	if err != nil {
		return nil, err
	}
	chat.GroupMemberCancelMutedNotification(tracelog.GetOperationID(ctx), tracelog.GetOpUserID(ctx), req.GroupID, req.UserID)
	return resp, nil
}

func (s *groupServer) MuteGroup(ctx context.Context, req *pbGroup.MuteGroupReq) (*pbGroup.MuteGroupResp, error) {
	resp := &pbGroup.MuteGroupResp{}

	opFlag, err := s.getGroupUserLevel(req.GroupID, tracelog.GetOpUserID(ctx))
	if err != nil {
		return nil, err
	}
	if opFlag == 0 {
		//errMsg := req.OperationID + "opFlag == 0  " + req.GroupID + req.OpUserID
		//log.Error(req.OperationID, errMsg)
		//return &pbGroup.MuteGroupResp{CommonResp: &pbGroup.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
		return nil, utils.Wrap(constant.ErrNoPermission, "")
	}

	//mutedInfo, err := relation.GetGroupMemberInfoByGroupIDAndUserID(req.GroupID, req.UserID)
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

	err = relation.OperateGroupStatus(req.GroupID, constant.GroupStatusMuted)
	if err != nil {
		return nil, err
	}

	chat.GroupMutedNotification(tracelog.GetOperationID(ctx), tracelog.GetOpUserID(ctx), req.GroupID)
	return resp, nil
}

func (s *groupServer) CancelMuteGroup(ctx context.Context, req *pbGroup.CancelMuteGroupReq) (*pbGroup.CancelMuteGroupResp, error) {
	resp := &pbGroup.CancelMuteGroupResp{}

	opFlag, err := s.getGroupUserLevel(req.GroupID, tracelog.GetOpUserID(ctx))
	if err != nil {
		return nil, err
	}
	if opFlag == 0 {
		return nil, err
	}
	//mutedInfo, err := relation.GetGroupMemberInfoByGroupIDAndUserID(req.GroupID, req.)
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
	log.Debug(tracelog.GetOperationID(ctx), "UpdateGroupInfoDefaultZero ", req.GroupID, map[string]interface{}{"status": constant.GroupOk})
	if err := rocksCache.DelGroupInfoFromCache(ctx, req.GroupID); err != nil {
		return nil, err
	}
	err = relation.UpdateGroupInfoDefaultZero(req.GroupID, map[string]interface{}{"status": constant.GroupOk})
	if err != nil {
		return nil, err
	}
	chat.GroupCancelMutedNotification(tracelog.GetOperationID(ctx), tracelog.GetOpUserID(ctx), req.GroupID)
	return resp, nil
}

func (s *groupServer) SetGroupMemberNickname(ctx context.Context, req *pbGroup.SetGroupMemberNicknameReq) (*pbGroup.SetGroupMemberNicknameResp, error) {
	resp := &pbGroup.SetGroupMemberNicknameResp{}
	if tracelog.GetOpUserID(ctx) != req.UserID && !token_verify.IsManagerUserID(tracelog.GetOpUserID(ctx)) {
		return nil, utils.Wrap(constant.ErrIdentity, "")
	}
	cbReq := &pbGroup.SetGroupMemberInfoReq{
		GroupID:  req.GroupID,
		UserID:   req.UserID,
		Nickname: &wrapperspb.StringValue{Value: req.Nickname},
	}
	if err := CallbackBeforeSetGroupMemberInfo(ctx, cbReq); err != nil {
		return nil, err
	}
	nickName := cbReq.Nickname.Value
	groupMemberInfo := relation2.GroupMemberModel{}
	groupMemberInfo.UserID = req.UserID
	groupMemberInfo.GroupID = req.GroupID
	if nickName == "" {
		userNickname, err := relation.GetUserNameByUserID(groupMemberInfo.UserID)
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

	if err := relation.UpdateGroupMemberInfo(groupMemberInfo); err != nil {
		return nil, err
	}
	chat.GroupMemberInfoSetNotification(tracelog.GetOperationID(ctx), tracelog.GetOpUserID(ctx), req.GroupID, req.UserID)
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
	groupMember := relation2.GroupMemberModel{
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
	if err := relation.UpdateGroupMemberInfoByMap(groupMember, m); err != nil {
		return nil, err
	}
	if req.RoleLevel != nil {
		switch req.RoleLevel.Value {
		case constant.GroupOrdinaryUsers:
			//msg.GroupMemberRoleLevelChangeNotification(req.OperationID, req.OpUserID, req.GroupID, req.UserID, constant.GroupMemberSetToOrdinaryUserNotification)
			chat.GroupMemberInfoSetNotification(tracelog.GetOperationID(ctx), tracelog.GetOpUserID(ctx), req.GroupID, req.UserID)
		case constant.GroupAdmin, constant.GroupOwner:
			//msg.GroupMemberRoleLevelChangeNotification(req.OperationID, req.OpUserID, req.GroupID, req.UserID, constant.GroupMemberSetToAdminNotification)
			chat.GroupMemberInfoSetNotification(tracelog.GetOperationID(ctx), tracelog.GetOpUserID(ctx), req.GroupID, req.UserID)
		}
	} else {
		chat.GroupMemberInfoSetNotification(tracelog.GetOperationID(ctx), tracelog.GetOpUserID(ctx), req.GroupID, req.UserID)
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
	operationID := tracelog.GetOperationID(ctx)
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
