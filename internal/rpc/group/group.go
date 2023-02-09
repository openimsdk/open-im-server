package group

import (
	chat "Open_IM/internal/rpc/msg"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/cache"
	"Open_IM/pkg/common/db/controller"
	relation_conn "Open_IM/pkg/common/db/relation"
	"Open_IM/pkg/common/db/table/relation"
	"Open_IM/pkg/common/db/unrelation"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/middleware"
	promePkg "Open_IM/pkg/common/prometheus"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/common/tracelog"
	discoveryRegistry "Open_IM/pkg/discovery_registry"
	pbGroup "Open_IM/pkg/proto/group"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"fmt"
	"github.com/OpenIMSDK/getcdv3"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"math/big"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/OpenIMSDK/openKeeper"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type groupServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
	GroupInterface  controller.GroupInterface

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
	var mysql relation_conn.Mysql
	var mongo unrelation.Mongo
	var groupModel relation.GroupModel
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
	zkClient, err := openKeeper.NewClient([]string{"43.154.157.177:2181"}, config.Config.Etcd.EtcdSchema, 10, "", "")
	if err != nil {
		panic(err.Error())
	}
	registerIP, err := utils.GetRpcIP()
	err = zkClient.Register(config.Config.RpcRegisterName.OpenImGroupName, registerIP, port)
	if err != nil {
		panic(err.Error())
	}
	var registerCenter discoveryRegistry.SvcDiscoveryRegistry = zkClient
	conns, err := registerCenter.GetConns(config.Config.RpcRegisterName.OpenImConversationName)
	g.GroupInterface = controller.NewGroupInterface(controller.NewGroupDatabase(mysql.GormConn(), redis.GetClient(), mongo.GetClient()))
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

func (s *groupServer) CheckGroupAdmin(ctx context.Context, groupID string) error {
	if !token_verify.IsAppManagerUid(ctx) {
		groupMember, err := s.GroupInterface.TakeGroupMember(ctx, groupID, tracelog.GetOpUserID(ctx))
		if err != nil {
			return err
		}
		if !(groupMember.RoleLevel == constant.GroupOwner || groupMember.RoleLevel == constant.GroupAdmin) {
			return constant.ErrNoPermission.Wrap("no group owner or admin")
		}
	}
	return nil
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
	var groupMembers []*relation.GroupMemberModel
	group := PbToDBGroupInfo(req.GroupInfo)
	group.GroupID = genGroupID(ctx, req.GroupInfo.GroupID)
	joinGroup := func(userID string, roleLevel int32) error {
		groupMember := PbToDbGroupMember(userMap[userID])
		groupMember.GroupID = group.GroupID
		groupMember.RoleLevel = roleLevel
		groupMember.OperatorUserID = tracelog.GetOpUserID(ctx)
		groupMember.JoinSource = constant.JoinByInvitation
		groupMember.InviterUserID = tracelog.GetOpUserID(ctx)
		if err := CallbackBeforeMemberJoinGroup(ctx, tracelog.GetOperationID(ctx), groupMember, group.Ex); err != nil {
			return err
		}
		groupMembers = append(groupMembers, groupMember)
		return nil
	}
	if err := joinGroup(req.OwnerUserID, constant.GroupOwner); err != nil {
		return nil, err
	}
	if req.GroupInfo.GroupType == constant.SuperGroup {
		if err := s.GroupInterface.CreateSuperGroup(ctx, group.GroupID, userIDs); err != nil {
			return nil, err
		}
	} else {
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
	if err := s.GroupInterface.CreateGroup(ctx, []*relation.GroupModel{group}, groupMembers); err != nil {
		return nil, err
	}
	resp.GroupInfo = DbToPbGroupInfo(group, req.OwnerUserID, uint32(len(userIDs)))
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
	total, members, err := s.GroupInterface.PageGroupMember(ctx, nil, []string{req.FromUserID}, nil, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}
	resp.Total = total
	if len(members) == 0 {
		return resp, nil
	}
	groupIDs := utils.Slice(members, func(e *relation.GroupMemberModel) string {
		return e.GroupID
	})
	groups, err := s.GroupInterface.FindGroup(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	groupMemberNum, err := s.GroupInterface.MapGroupMemberNum(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	owners, err := s.GroupInterface.FindGroupMember(ctx, groupIDs, nil, []int32{constant.GroupOwner})
	if err != nil {
		return nil, err
	}
	ownerMap := utils.SliceToMap(owners, func(e *relation.GroupMemberModel) string {
		return e.GroupID
	})
	resp.Groups = utils.Slice(utils.Order(groupIDs, groups, func(group *relation.GroupModel) string {
		return group.GroupID
	}), func(group *relation.GroupModel) *open_im_sdk.GroupInfo {
		return DbToPbGroupInfo(group, ownerMap[group.GroupID].UserID, uint32(groupMemberNum[group.GroupID]))
	})
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
	group, err := s.GroupInterface.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if group.Status == constant.GroupStatusDismissed {
		return nil, constant.ErrDismissedAlready.Wrap()
	}
	members, err := s.GroupInterface.FindGroupMember(ctx, []string{group.GroupID}, nil, nil)
	if err != nil {
		return nil, err
	}
	memberMap := utils.SliceToMap(members, func(e *relation.GroupMemberModel) string {
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
				var requests []*relation.GroupRequestModel
				for _, userID := range req.InvitedUserIDs {
					requests = append(requests, &relation.GroupRequestModel{
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
		if err := s.GroupInterface.CreateSuperGroupMember(ctx, req.GroupID, req.InvitedUserIDs); err != nil {
			return nil, err
		}
		for _, userID := range req.InvitedUserIDs {
			chat.SuperGroupNotification(tracelog.GetOperationID(ctx), userID, userID)
		}
	} else {
		opUserID := tracelog.GetOpUserID(ctx)
		var groupMembers []*relation.GroupMemberModel
		for _, userID := range req.InvitedUserIDs {
			member := PbToDbGroupMember(userMap[userID])
			member.GroupID = req.GroupID
			member.RoleLevel = constant.GroupOrdinaryUsers
			member.OperatorUserID = opUserID
			member.InviterUserID = opUserID
			member.JoinSource = constant.JoinByInvitation
			if err := CallbackBeforeMemberJoinGroup(ctx, tracelog.GetOperationID(ctx), member, group.Ex); err != nil {
				return nil, err
			}
			groupMembers = append(groupMembers, member)
		}
		if err := s.GroupInterface.CreateGroup(ctx, nil, groupMembers); err != nil {
			return nil, err
		}
		chat.MemberInvitedNotification(tracelog.GetOperationID(ctx), req.GroupID, tracelog.GetOpUserID(ctx), req.Reason, req.InvitedUserIDs)
	}
	return resp, nil
}

func (s *groupServer) GetGroupAllMember(ctx context.Context, req *pbGroup.GetGroupAllMemberReq) (*pbGroup.GetGroupAllMemberResp, error) {
	resp := &pbGroup.GetGroupAllMemberResp{}
	group, err := s.GroupInterface.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if group.GroupType == constant.SuperGroup {
		return nil, constant.ErrArgs.Wrap("unsupported super group")
	}
	members, err := s.GroupInterface.FindGroupMember(ctx, []string{req.GroupID}, nil, nil)
	if err != nil {
		return nil, err
	}
	resp.Members = utils.Slice(members, func(e *relation.GroupMemberModel) *open_im_sdk.GroupMemberFullInfo {
		return DbToPbGroupMembersCMSResp(e)
	})
	return resp, nil
}

func (s *groupServer) GetGroupMemberList(ctx context.Context, req *pbGroup.GetGroupMemberListReq) (*pbGroup.GetGroupMemberListResp, error) {
	resp := &pbGroup.GetGroupMemberListResp{}
	total, members, err := s.GroupInterface.PageGroupMember(ctx, []string{req.GroupID}, nil, utils.If(req.Filter >= 0, []int32{req.Filter}, nil), req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}
	resp.Total = total
	resp.Members = utils.Slice(members, func(e *relation.GroupMemberModel) *open_im_sdk.GroupMemberFullInfo {
		return DbToPbGroupMembersCMSResp(e)
	})
	return resp, nil
}

func (s *groupServer) KickGroupMember(ctx context.Context, req *pbGroup.KickGroupMemberReq) (*pbGroup.KickGroupMemberResp, error) {
	resp := &pbGroup.KickGroupMemberResp{}
	group, err := s.GroupInterface.TakeGroup(ctx, req.GroupID)
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
		if err := s.GroupInterface.DeleteSuperGroupMember(ctx, req.GroupID, req.KickedUserIDs); err != nil {
			return nil, err
		}
		go func() {
			for _, userID := range req.KickedUserIDs {
				chat.SuperGroupNotification(tracelog.GetOperationID(ctx), userID, userID)
			}
		}()
	} else {
		members, err := s.GroupInterface.FindGroupMember(ctx, []string{req.GroupID}, append(req.KickedUserIDs, opUserID), nil)
		if err != nil {
			return nil, err
		}
		memberMap := make(map[string]*relation.GroupMemberModel)
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
		if err := s.GroupInterface.DeleteGroupMember(ctx, group.GroupID, req.KickedUserIDs); err != nil {
			return nil, err
		}
		chat.MemberKickedNotification(req, req.KickedUserIDs)
	}
	return resp, nil
}

func (s *groupServer) GetGroupMembersInfo(ctx context.Context, req *pbGroup.GetGroupMembersInfoReq) (*pbGroup.GetGroupMembersInfoResp, error) {
	resp := &pbGroup.GetGroupMembersInfoResp{}
	if len(req.Members) == 0 {
		return nil, constant.ErrArgs.Wrap("members empty")
	}
	if req.GroupID == "" {
		return nil, constant.ErrArgs.Wrap("groupID empty")
	}
	members, err := s.GroupInterface.FindGroupMember(ctx, []string{req.GroupID}, req.Members, nil)
	if err != nil {
		return nil, err
	}
	resp.Members = utils.Slice(members, func(e *relation.GroupMemberModel) *open_im_sdk.GroupMemberFullInfo {
		return DbToPbGroupMembersCMSResp(e)
	})
	return resp, nil
}

func (s *groupServer) GetGroupApplicationList(ctx context.Context, req *pbGroup.GetGroupApplicationListReq) (*pbGroup.GetGroupApplicationListResp, error) {
	resp := &pbGroup.GetGroupApplicationListResp{}
	total, groupRequests, err := s.GroupInterface.PageGroupRequestUser(ctx, req.FromUserID, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}
	resp.Total = total
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
	groups, err := s.GroupInterface.FindGroup(ctx, utils.Distinct(groupIDs))
	if err != nil {
		return nil, err
	}
	groupMap := utils.SliceToMap(groups, func(e *relation.GroupModel) string {
		return e.GroupID
	})
	if ids := utils.Single(utils.Keys(groupMap), groupIDs); len(ids) > 0 {
		return nil, constant.ErrGroupIDNotFound.Wrap(strings.Join(ids, ","))
	}
	groupMemberNumMap, err := s.GroupInterface.MapGroupMemberNum(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	owners, err := s.GroupInterface.FindGroupMember(ctx, groupIDs, nil, []int32{constant.GroupOwner})
	if err != nil {
		return nil, err
	}
	ownerMap := utils.SliceToMap(owners, func(e *relation.GroupMemberModel) string {
		return e.GroupID
	})
	resp.GroupRequests = utils.Slice(groupRequests, func(e *relation.GroupRequestModel) *open_im_sdk.GroupRequest {
		return DbToPbGroupRequest(e, userMap[e.UserID], DbToPbGroupInfo(groupMap[e.GroupID], ownerMap[e.GroupID].UserID, uint32(groupMemberNumMap[e.GroupID])))
	})
	return resp, nil
}

func (s *groupServer) GetGroupsInfo(ctx context.Context, req *pbGroup.GetGroupsInfoReq) (*pbGroup.GetGroupsInfoResp, error) {
	resp := &pbGroup.GetGroupsInfoResp{}
	if len(req.GroupIDs) == 0 {
		return nil, constant.ErrArgs.Wrap("groupID is empty")
	}
	groups, err := s.GroupInterface.FindGroup(ctx, req.GroupIDs)
	if err != nil {
		return nil, err
	}
	groupMemberNumMap, err := s.GroupInterface.MapGroupMemberNum(ctx, req.GroupIDs)
	if err != nil {
		return nil, err
	}
	owners, err := s.GroupInterface.FindGroupMember(ctx, req.GroupIDs, nil, []int32{constant.GroupOwner})
	if err != nil {
		return nil, err
	}
	ownerMap := utils.SliceToMap(owners, func(e *relation.GroupMemberModel) string {
		return e.GroupID
	})
	resp.GroupInfos = utils.Slice(groups, func(e *relation.GroupModel) *open_im_sdk.GroupInfo {
		return DbToPbGroupInfo(e, ownerMap[e.GroupID].UserID, uint32(groupMemberNumMap[e.GroupID]))
	})
	return resp, nil
}

func (s *groupServer) GroupApplicationResponse(ctx context.Context, req *pbGroup.GroupApplicationResponseReq) (*pbGroup.GroupApplicationResponseResp, error) {
	resp := &pbGroup.GroupApplicationResponseResp{}
	if !utils.Contain(req.HandleResult, constant.GroupResponseAgree, constant.GroupResponseRefuse) {
		return nil, constant.ErrArgs.Wrap("HandleResult unknown")
	}
	if !token_verify.IsAppManagerUid(ctx) {
		groupMember, err := s.GroupInterface.TakeGroupMember(ctx, req.GroupID, req.FromUserID)
		if err != nil {
			return nil, err
		}
		if !(groupMember.RoleLevel == constant.GroupOwner || groupMember.RoleLevel == constant.GroupAdmin) {
			return nil, constant.ErrNoPermission.Wrap("no group owner or admin")
		}
	}
	group, err := s.GroupInterface.TakeGroup(ctx, req.GroupID)
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
	if _, err := s.GroupInterface.TakeGroupMember(ctx, req.GroupID, req.FromUserID); err != nil {
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
	var member *relation.GroupMemberModel
	if req.HandleResult == constant.GroupResponseAgree {
		member = &relation.GroupMemberModel{
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
	group, err := s.GroupInterface.TakeGroup(ctx, req.GroupID)
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
		user, err := relation_conn.GetUserByUserID(tracelog.GetOpUserID(ctx))
		if err != nil {
			return nil, err
		}
		groupMember := PbToDbGroupMember(user)
		groupMember.GroupID = group.GroupID
		groupMember.RoleLevel = constant.GroupOrdinaryUsers
		groupMember.OperatorUserID = tracelog.GetOpUserID(ctx)
		groupMember.JoinSource = constant.JoinByInvitation
		groupMember.InviterUserID = tracelog.GetOpUserID(ctx)
		if err := CallbackBeforeMemberJoinGroup(ctx, tracelog.GetOperationID(ctx), groupMember, group.Ex); err != nil {
			return nil, err
		}
		if err := s.GroupInterface.CreateGroup(ctx, nil, []*relation.GroupMemberModel{groupMember}); err != nil {
			return nil, err
		}
		chat.MemberEnterDirectlyNotification(req.GroupID, tracelog.GetOpUserID(ctx), tracelog.GetOperationID(ctx))
		return resp, nil
	}
	groupRequest := relation.GroupRequestModel{
		UserID:     tracelog.GetOpUserID(ctx),
		ReqMsg:     req.ReqMessage,
		GroupID:    req.GroupID,
		JoinSource: req.JoinSource,
		ReqTime:    time.Now(),
	}
	if err := s.GroupInterface.CreateGroupRequest(ctx, []*relation.GroupRequestModel{&groupRequest}); err != nil {
		return nil, err
	}
	chat.JoinGroupApplicationNotification(ctx, req)
	return resp, nil
}

func (s *groupServer) QuitGroup(ctx context.Context, req *pbGroup.QuitGroupReq) (*pbGroup.QuitGroupResp, error) {
	resp := &pbGroup.QuitGroupResp{}
	group, err := s.GroupInterface.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if group.GroupType == constant.SuperGroup {
		if err := s.GroupInterface.DeleteSuperGroupMember(ctx, req.GroupID, []string{tracelog.GetOpUserID(ctx)}); err != nil {
			return nil, err
		}
		chat.SuperGroupNotification(tracelog.GetOperationID(ctx), tracelog.GetOpUserID(ctx), tracelog.GetOpUserID(ctx))
	} else {
		_, err := s.GroupInterface.TakeGroupMember(ctx, req.GroupID, tracelog.GetOpUserID(ctx))
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
		groupMember, err := s.GroupInterface.TakeGroupMember(ctx, req.GroupInfoForSet.GroupID, tracelog.GetOpUserID(ctx))
		if err != nil {
			return nil, err
		}
		if !(groupMember.RoleLevel == constant.GroupOwner || groupMember.RoleLevel == constant.GroupAdmin) {
			return nil, constant.ErrNoPermission.Wrap("no group owner or admin")
		}
	}
	group, err := s.TakeGroup(ctx, req.GroupInfoForSet.GroupID)
	if err != nil {
		return nil, err
	}
	if group.Status == constant.GroupStatusDismissed {
		return nil, utils.Wrap(constant.ErrDismissedAlready, "")
	}
	data := UpdateGroupInfoMap(req.GroupInfoForSet)
	if len(data) > 0 {
		return resp, nil
	}
	if err := s.GroupInterface.UpdateGroup(ctx, group.GroupID, data); err != nil {
		return nil, err
	}
	group, err = s.TakeGroup(ctx, req.GroupInfoForSet.GroupID)
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
	group, err := s.GroupInterface.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if group.Status == constant.GroupStatusDismissed {
		return nil, utils.Wrap(constant.ErrDismissedAlready, "")
	}
	if req.OldOwnerUserID == req.NewOwnerUserID {
		return nil, constant.ErrArgs.Wrap("OldOwnerUserID == NewOwnerUserID")
	}
	members, err := s.GroupInterface.FindGroupMember(ctx, []string{req.GroupID}, []string{req.OldOwnerUserID, req.NewOwnerUserID}, nil)
	if err != nil {
		return nil, err
	}
	memberMap := utils.SliceToMap(members, func(e *relation.GroupMemberModel) string { return e.UserID })
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
			oldOwner, err = s.GroupInterface.TakeGroupOwner(ctx, req.OldOwnerUserID)
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
	if err := s.GroupInterface.TransferGroupOwner(ctx, req.GroupID, req.OldOwnerUserID, req.NewOwnerUserID, newOwner.RoleLevel); err != nil {
		return nil, err
	}
	chat.GroupOwnerTransferredNotification(req)
	return resp, nil
}

func (s *groupServer) GetGroups(ctx context.Context, req *pbGroup.GetGroupsReq) (*pbGroup.GetGroupsResp, error) {
	resp := &pbGroup.GetGroupsResp{}
	var (
		groups []*relation.GroupModel
		err    error
	)
	if req.GroupID != "" {
		groups, err = s.GroupInterface.FindGroup(ctx, []string{req.GroupID})
		resp.Total = uint32(len(groups))
	} else {
		resp.Total, groups, err = s.GroupInterface.SearchGroup(ctx, req.GroupName, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	}
	if err != nil {
		return nil, err
	}
	groupIDs := utils.Slice(groups, func(e *relation.GroupModel) string {
		return e.GroupID
	})
	ownerMembers, err := s.GroupInterface.FindGroupMember(ctx, groupIDs, nil, []int32{constant.GroupOwner})
	if err != nil {
		return nil, err
	}
	ownerMemberMap := utils.SliceToMap(ownerMembers, func(e *relation.GroupMemberModel) string {
		return e.GroupID
	})
	if ids := utils.Single(groupIDs, utils.Keys(ownerMemberMap)); len(ids) > 0 {
		return nil, constant.ErrDB.Wrap("group not owner " + strings.Join(ids, ","))
	}
	groupMemberNumMap, err := s.GroupInterface.MapGroupMemberNum(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	resp.Groups = utils.Slice(groups, func(group *relation.GroupModel) *pbGroup.CMSGroup {
		member := ownerMemberMap[group.GroupID]
		return DbToPbCMSGroup(group, member.UserID, member.Nickname, uint32(groupMemberNumMap[group.GroupID]))
	})
	return resp, nil
}

func (s *groupServer) GetGroupMembersCMS(ctx context.Context, req *pbGroup.GetGroupMembersCMSReq) (*pbGroup.GetGroupMembersCMSResp, error) {
	resp := &pbGroup.GetGroupMembersCMSResp{}
	total, members, err := s.GroupInterface.SearchGroupMember(ctx, req.UserName, []string{req.GroupID}, nil, nil, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}
	resp.Total = total
	resp.Members = utils.Slice(members, func(e *relation.GroupMemberModel) *open_im_sdk.GroupMemberFullInfo {
		return DbToPbGroupMembersCMSResp(e)
	})
	return resp, nil
}

func (s *groupServer) GetUserReqApplicationList(ctx context.Context, req *pbGroup.GetUserReqApplicationListReq) (*pbGroup.GetUserReqApplicationListResp, error) {
	resp := &pbGroup.GetUserReqApplicationListResp{}
	user, err := GetPublicUserInfoOne(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	total, requests, err := s.GroupInterface.PageGroupRequestUser(ctx, req.UserID, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}
	resp.Total = total
	if len(requests) == 0 {
		return resp, nil
	}
	groupIDs := utils.Distinct(utils.Slice(requests, func(e *relation.GroupRequestModel) string {
		return e.GroupID
	}))
	groups, err := s.GroupInterface.FindGroup(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	groupMap := utils.SliceToMap(groups, func(e *relation.GroupModel) string {
		return e.GroupID
	})
	if ids := utils.Single(groupIDs, utils.Keys(groupMap)); len(ids) > 0 {
		return nil, constant.ErrGroupIDNotFound.Wrap(strings.Join(ids, ","))
	}
	owners, err := s.GroupInterface.FindGroupMember(ctx, groupIDs, nil, []int32{constant.GroupOwner})
	if err != nil {
		return nil, err
	}
	ownerMap := utils.SliceToMap(owners, func(e *relation.GroupMemberModel) string {
		return e.GroupID
	})
	if ids := utils.Single(groupIDs, utils.Keys(ownerMap)); len(ids) > 0 {
		return nil, constant.ErrData.Wrap("group no owner", strings.Join(ids, ","))
	}
	groupMemberNum, err := s.GroupInterface.MapGroupMemberNum(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	resp.GroupRequests = utils.Slice(requests, func(e *relation.GroupRequestModel) *open_im_sdk.GroupRequest {
		return DbToPbGroupRequest(e, user, DbToPbGroupInfo(groupMap[e.GroupID], ownerMap[e.GroupID].UserID, uint32(groupMemberNum[e.GroupID])))
	})
	return resp, nil
}

func (s *groupServer) DismissGroup(ctx context.Context, req *pbGroup.DismissGroupReq) (*pbGroup.DismissGroupResp, error) {
	resp := &pbGroup.DismissGroupResp{}
	if err := s.CheckGroupAdmin(ctx, req.GroupID); err != nil {
		return nil, err
	}
	group, err := s.GroupInterface.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if group.Status == constant.GroupStatusDismissed {
		return nil, constant.ErrArgs.Wrap("group status is dismissed")
	}
	if err := s.GroupInterface.DismissGroup(ctx, req.GroupID); err != nil {
		return nil, err
	}
	if group.GroupType == constant.SuperGroup {
		if err := s.GroupInterface.DeleteSuperGroup(ctx, group.GroupID); err != nil {
			return nil, err
		}
	} else {
		chat.GroupDismissedNotification(req)
	}
	return resp, nil
}

func (s *groupServer) MuteGroupMember(ctx context.Context, req *pbGroup.MuteGroupMemberReq) (*pbGroup.MuteGroupMemberResp, error) {
	resp := &pbGroup.MuteGroupMemberResp{}
	member, err := s.GroupInterface.TakeGroupMember(ctx, req.GroupID, req.UserID)
	if err != nil {
		return nil, err
	}
	if !(tracelog.GetOpUserID(ctx) == req.UserID || token_verify.IsAppManagerUid(ctx)) {
		opMember, err := s.GroupInterface.TakeGroupMember(ctx, req.GroupID, req.UserID)
		if err != nil {
			return nil, err
		}
		if opMember.RoleLevel <= member.RoleLevel {
			return nil, constant.ErrNoPermission.Wrap(fmt.Sprintf("self RoleLevel %d target %d", opMember.RoleLevel, member.RoleLevel))
		}
	}
	data := UpdateGroupMemberMutedTimeMap(time.Now().Add(time.Second * time.Duration(req.MutedSeconds)))
	if err := s.GroupInterface.UpdateGroupMember(ctx, member.GroupID, member.UserID, data); err != nil {
		return nil, err
	}
	chat.GroupMemberMutedNotification(tracelog.GetOperationID(ctx), tracelog.GetOpUserID(ctx), req.GroupID, req.UserID, req.MutedSeconds)
	return resp, nil
}

func (s *groupServer) CancelMuteGroupMember(ctx context.Context, req *pbGroup.CancelMuteGroupMemberReq) (*pbGroup.CancelMuteGroupMemberResp, error) {
	resp := &pbGroup.CancelMuteGroupMemberResp{}
	member, err := s.GroupInterface.TakeGroupMember(ctx, req.GroupID, req.UserID)
	if err != nil {
		return nil, err
	}
	if !(tracelog.GetOpUserID(ctx) == req.UserID || token_verify.IsAppManagerUid(ctx)) {
		opMember, err := s.GroupInterface.TakeGroupMember(ctx, req.GroupID, tracelog.GetOpUserID(ctx))
		if err != nil {
			return nil, err
		}
		if opMember.RoleLevel <= member.RoleLevel {
			return nil, constant.ErrNoPermission.Wrap(fmt.Sprintf("self RoleLevel %d target %d", opMember.RoleLevel, member.RoleLevel))
		}
	}
	data := UpdateGroupMemberMutedTimeMap(time.Unix(0, 0))
	if err := s.GroupInterface.UpdateGroupMember(ctx, member.GroupID, member.UserID, data); err != nil {
		return nil, err
	}
	chat.GroupMemberCancelMutedNotification(tracelog.GetOperationID(ctx), tracelog.GetOpUserID(ctx), req.GroupID, req.UserID)
	return resp, nil
}

func (s *groupServer) MuteGroup(ctx context.Context, req *pbGroup.MuteGroupReq) (*pbGroup.MuteGroupResp, error) {
	resp := &pbGroup.MuteGroupResp{}
	if err := s.CheckGroupAdmin(ctx, req.GroupID); err != nil {
		return nil, err
	}
	if err := s.GroupInterface.UpdateGroup(ctx, req.GroupID, UpdateGroupStatusMap(constant.GroupStatusMuted)); err != nil {
		return nil, err
	}
	chat.GroupMutedNotification(tracelog.GetOperationID(ctx), tracelog.GetOpUserID(ctx), req.GroupID)
	return resp, nil
}

func (s *groupServer) CancelMuteGroup(ctx context.Context, req *pbGroup.CancelMuteGroupReq) (*pbGroup.CancelMuteGroupResp, error) {
	resp := &pbGroup.CancelMuteGroupResp{}
	if err := s.CheckGroupAdmin(ctx, req.GroupID); err != nil {
		return nil, err
	}
	if err := s.GroupInterface.UpdateGroup(ctx, req.GroupID, UpdateGroupStatusMap(constant.GroupOk)); err != nil {
		return nil, err
	}
	chat.GroupCancelMutedNotification(tracelog.GetOperationID(ctx), tracelog.GetOpUserID(ctx), req.GroupID)
	return resp, nil
}

func (s *groupServer) SetGroupMemberNickname(ctx context.Context, req *pbGroup.SetGroupMemberNicknameReq) (*pbGroup.SetGroupMemberNicknameResp, error) {
	_, err := s.SetGroupMemberInfo(ctx, &pbGroup.SetGroupMemberInfoReq{GroupID: req.GroupID, UserID: req.UserID, Nickname: wrapperspb.String(req.Nickname)})
	if err != nil {
		return nil, err
	}
	return &pbGroup.SetGroupMemberNicknameResp{}, nil
}

func (s *groupServer) SetGroupMemberInfo(ctx context.Context, req *pbGroup.SetGroupMemberInfoReq) (*pbGroup.SetGroupMemberInfoResp, error) {
	resp := &pbGroup.SetGroupMemberInfoResp{}
	if req.RoleLevel != nil && req.RoleLevel.Value == constant.GroupOwner {
		return nil, constant.ErrNoPermission.Wrap("set group owner")
	}
	group, err := s.GroupInterface.TakeGroup(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if group.Status == constant.GroupStatusDismissed {
		return nil, constant.ErrArgs.Wrap("group status is dismissed")
	}
	member, err := s.GroupInterface.TakeGroupMember(ctx, req.GroupID, req.UserID)
	if err != nil {
		return nil, err
	}
	if tracelog.GetOpUserID(ctx) == req.UserID {
		if req.RoleLevel != nil {
			return nil, constant.ErrArgs.Wrap("update role level")
		}
	} else if !token_verify.IsAppManagerUid(ctx) {
		opMember, err := s.GroupInterface.TakeGroupMember(ctx, req.GroupID, tracelog.GetOpUserID(ctx))
		if err != nil {
			return nil, err
		}
		if opMember.RoleLevel <= member.RoleLevel {
			return nil, constant.ErrNoPermission.Wrap(fmt.Sprintf("self RoleLevel %d target %d", opMember.RoleLevel, member.RoleLevel))
		}
	}
	if err := CallbackBeforeSetGroupMemberInfo(ctx, req); err != nil {
		return nil, err
	}
	if err := s.GroupInterface.UpdateGroupMember(ctx, req.GroupID, req.UserID, UpdateGroupMemberMap(req)); err != nil {
		return nil, err
	}
	chat.GroupMemberInfoSetNotification(tracelog.GetOperationID(ctx), tracelog.GetOpUserID(ctx), req.GroupID, req.UserID)
	return resp, nil
}

func (s *groupServer) GetGroupAbstractInfo(ctx context.Context, req *pbGroup.GetGroupAbstractInfoReq) (*pbGroup.GetGroupAbstractInfoResp, error) {
	resp := &pbGroup.GetGroupAbstractInfoResp{}
	if len(req.GroupIDs) == 0 {
		return nil, constant.ErrArgs.Wrap("groupIDs empty")
	}
	if utils.Duplicate(req.GroupIDs) {
		return nil, constant.ErrArgs.Wrap("groupIDs duplicate")
	}
	groups, err := s.GroupInterface.FindGroup(ctx, req.GroupIDs)
	if err != nil {
		return nil, err
	}
	groupUserMap, err := s.GroupInterface.MapGroupMemberUserID(ctx, req.GroupIDs)
	if err != nil {
		return nil, err
	}
	resp.GroupAbstractInfos = utils.Slice(groups, func(e *relation.GroupModel) *pbGroup.GroupAbstractInfo {
		userIDs := groupUserMap[e.GroupID]
		utils.Sort(userIDs, true)
		bi := big.NewInt(0)
		bi.SetString(utils.Md5(strings.Join(userIDs, ";;"))[0:8], 16)
		return DbToPbGroupAbstractInfo(e.GroupID, int32(len(userIDs)), bi.Uint64())
	})
	return resp, nil
}

func (s *groupServer) GetUserInGroupMembers(ctx context.Context, req *pbGroup.GetUserInGroupMembersReq) (*pbGroup.GetUserInGroupMembersResp, error) {
	resp := &pbGroup.GetUserInGroupMembersResp{}
	if len(req.GroupIDs) == 0 {
		return nil, constant.ErrArgs.Wrap("groupIDs empty")
	}
	members, err := s.GroupInterface.FindGroupMember(ctx, []string{req.UserID}, req.GroupIDs, nil)
	if err != nil {
		return nil, err
	}
	resp.Members = utils.Slice(members, func(e *relation.GroupMemberModel) *open_im_sdk.GroupMemberFullInfo {
		return DbToPbGroupMembersCMSResp(e)
	})
	return resp, nil
}
