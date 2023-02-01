package group

import (
	"Open_IM/internal/rpc/fault_tolerant"
	chat "Open_IM/internal/rpc/msg"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/cache"
	"Open_IM/pkg/common/db/controller"
	"Open_IM/pkg/common/db/relation"
	"Open_IM/pkg/common/db/table"
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
	"errors"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"net"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"
)

type groupServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
	controller.GroupInterface

	userRpc         pbUser.UserClient
	conversationRpc pbConversation.ConversationClient
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
	etcdClient, err := getcdv3.NewEtcdConn(g.rpcRegisterName, strings.Join(g.etcdAddr, ","), config.Config.RpcRegisterIP, config.Config.Etcd.UserName, config.Config.Etcd.Password, port, ttl)
	if err != nil {
		panic("NewEtcdConn failed" + err.Error())
	}
	err = etcdClient.RegisterEtcd("", g.rpcRegisterName)
	if err != nil {
		panic("NewEtcdConn failed" + err.Error())
	}
	etcdClient.SetDefaultEtcdConfig(config.Config.RpcRegisterName.OpenImUserName, config.Config.RpcPort.OpenImUserPort)
	conn := etcdClient.GetConn("", config.Config.RpcRegisterName.OpenImUserName)
	g.userRpc = pbUser.NewUserClient(conn)

	etcdClient.SetDefaultEtcdConfig(config.Config.RpcRegisterName.OpenImConversationName, config.Config.RpcPort.OpenImConversationPort)
	conn = etcdClient.GetConn("", config.Config.RpcRegisterName.OpenImConversationName)
	g.conversationRpc = pbConversation.NewConversationClient(conn)

	//mysql init
	var mysql relation.Mysql
	var mongo unrelation.Mongo
	var groupModel table.GroupModel
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
	var userIDs []string
	for _, userID := range req.InitMembers {
		userIDs = append(userIDs, userID)
	}
	for _, userID := range req.AdminUserIDs {
		userIDs = append(userIDs, userID)
	}
	userIDs = append(userIDs, req.OwnerUserID)
	if utils.IsDuplicateID(userIDs) {
		return nil, constant.ErrArgs.Wrap("group member repeated")
	}
	users, err := getUsersInfo(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	userMap := make(map[string]*open_im_sdk.UserInfo)
	for i, user := range users {
		userMap[user.UserID] = users[i]
	}
	for _, userID := range userIDs {
		if userMap[userID] == nil {
			return nil, constant.ErrUserIDNotFound.Wrap(userID)
		}
	}
	if err := callbackBeforeCreateGroup(ctx, req); err != nil {
		return nil, err
	}
	var group table.GroupModel
	var groupMembers []*table.GroupMemberModel
	utils.CopyStructFields(&group, req.GroupInfo)
	group.GroupID = genGroupID(ctx, req.GroupInfo.GroupID)
	if req.GroupInfo.GroupType == constant.SuperGroup {
		if err := s.GroupInterface.CreateSuperGroup(ctx, group.GroupID, userIDs); err != nil {
			return nil, err
		}
	} else {
		joinGroup := func(userID string, roleLevel int32) error {
			user := userMap[userID]
			groupMember := &table.GroupMemberModel{GroupID: group.GroupID, RoleLevel: roleLevel, OperatorUserID: tracelog.GetOpUserID(ctx), JoinSource: constant.JoinByInvitation, InviterUserID: tracelog.GetOpUserID(ctx)}
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
	if err := s.GroupInterface.CreateGroup(ctx, []*table.GroupModel{&group}, groupMembers); err != nil {
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
	var groupIDs []string
	for _, group := range groups {
		groupIDs = append(groupIDs, group.GroupID)
	}
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
	if utils.IsDuplicateID(req.InvitedUserIDs) {
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
	memberMap := make(map[string]*table.GroupMemberModel)
	for i, member := range members {
		memberMap[member.GroupID] = members[i]
	}
	for _, userID := range req.InvitedUserIDs {
		if _, ok := memberMap[userID]; ok {
			return nil, constant.ErrArgs.Wrap("user in group " + userID)
		}
	}
	userMap, err := getUserMap(ctx, req.InvitedUserIDs)
	if err != nil {
		return nil, err
	}
	for _, userID := range req.InvitedUserIDs {
		if _, ok := userMap[userID]; !ok {
			return nil, constant.ErrUserIDNotFound.Wrap(userID)
		}
	}
	if group.NeedVerification == constant.AllNeedVerification {
		if !token_verify.IsAppManagerUid(ctx) {
			opUserID := tracelog.GetOpUserID(ctx)
			member, ok := memberMap[opUserID]
			if !ok {
				return nil, constant.ErrNoPermission.Wrap("not in group")
			}
			if !(member.RoleLevel == constant.GroupOwner || member.RoleLevel == constant.GroupAdmin) {
				var requests []*table.GroupRequestModel
				for _, userID := range req.InvitedUserIDs {
					requests = append(requests, &table.GroupRequestModel{
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
		var groupMembers []*table.GroupMemberModel
		for _, userID := range req.InvitedUserIDs {
			user := userMap[userID]
			var member table.GroupMemberModel
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
		var userIDs []string
		for _, member := range members {
			userIDs = append(userIDs, member.UserID)
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
		memberMap := make(map[string]*table.GroupMemberModel)
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
	userMap, err := getUserMap(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	for _, userID := range userIDs {
		if _, ok := userMap[userID]; !ok {
			return nil, constant.ErrUserIDNotFound.Wrap(userID)
		}
	}
	groups, err := s.GroupInterface.FindGroupsByID(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	groupMap := make(map[string]*table.GroupModel)
	for i, group := range groups {
		groupMap[group.GroupID] = groups[i]
	}

	for _, gr := range groupRequests {
		groupRequest := open_im_sdk.GroupRequest{UserInfo: &open_im_sdk.PublicUserInfo{}, GroupInfo: &open_im_sdk.GroupInfo{}}
		utils.CopyStructFields(&groupRequest, gr)
		getUserMap()

	}

	var errResult error
	tracelog.SetCtxInfo(ctx, "GetRecvGroupApplicationList", nil, " FromUserID: ", req.FromUserID, "GroupApplicationList: ", groupRequests)
	for _, v := range groupRequests {
		node := open_im_sdk.GroupRequest{UserInfo: &open_im_sdk.PublicUserInfo{}, GroupInfo: &open_im_sdk.GroupInfo{}}
		err := FillGroupInfoByGroupID(tracelog.GetOperationID(ctx), v.GroupID, node.GroupInfo)
		if err != nil {
			if !errors.Is(errors.Unwrap(err), constant.ErrDismissedAlready) {
				errResult = err
			}
			continue
		}
		tracelog.SetCtxInfo(ctx, "FillGroupInfoByGroupID ", nil, " groupID: ", v.GroupID, " groupInfo: ", node.GroupInfo)
		err = FillPublicUserInfoByUserID(tracelog.GetOperationID(ctx), v.UserID, node.UserInfo)
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
	tracelog.SetRpcRespInfo(ctx, utils.GetSelfFuncName(), resp.String())
	return resp, nil
}

func (s *groupServer) GetGroupsInfo(ctx context.Context, req *pbGroup.GetGroupsInfoReq) (*pbGroup.GetGroupsInfoResp, error) {
	resp := &pbGroup.GetGroupsInfoResp{}
	groupsInfoList := make([]*open_im_sdk.GroupInfo, 0)
	for _, groupID := range req.GroupIDList {
		groupInfoFromRedis, err := rocksCache.GetGroupInfoFromCache(ctx, groupID)
		if err != nil {
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
		tracelog.SetCtxInfo(ctx, utils.GetSelfFuncName(), err, "groupID", groupID, "userID", userID)
	}()
	if !token_verify.IsManagerUserID(userID) && !relation.IsGroupOwnerAdmin(groupID, userID) {
		return utils.Wrap(constant.ErrNoPermission, utils.GetSelfFuncName())
	}
	return nil
}

func (s *groupServer) GroupApplicationResponse(ctx context.Context, req *pbGroup.GroupApplicationResponseReq) (*pbGroup.GroupApplicationResponseResp, error) {
	resp := &pbGroup.GroupApplicationResponseResp{}

	if err := CheckPermission(ctx, req.GroupID, tracelog.GetOpUserID(ctx)); err != nil {
		return nil, err
	}
	groupRequest := getDBGroupRequest(ctx, req)
	if err := (&table.GroupRequestModel{}).Update(ctx, []*table.GroupRequestModel{groupRequest}); err != nil {
		return nil, err
	}
	groupInfo, err := rocksCache.GetGroupInfoFromCache(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if req.HandleResult == constant.GroupResponseAgree {
		member, err := getDBGroupMember(ctx, req.GroupID, req.FromUserID)
		if err != nil {
			return nil, err
		}
		err = CallbackBeforeMemberJoinGroup(ctx, tracelog.GetOperationID(ctx), member, groupInfo.Ex)
		if err != nil {
			return nil, err
		}
		err = (&table.GroupMemberModel{}).Create(ctx, []*table.GroupMemberModel{member})
		if err != nil {
			return nil, err
		}
		etcdCacheConn, err := fault_tolerant.GetDefaultConn(config.Config.RpcRegisterName.OpenImCacheName, tracelog.GetOperationID(ctx))
		if err != nil {
			return nil, err
		}
		cacheClient := pbCache.NewCacheClient(etcdCacheConn)
		cacheResp, err := cacheClient.DelGroupMemberIDListFromCache(context.Background(), &pbCache.DelGroupMemberIDListFromCacheReq{OperationID: tracelog.GetOperationID(ctx), GroupID: req.GroupID})
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
		return nil, constant.ErrArgs.Wrap()
	}
	return resp, nil
}

func (s *groupServer) JoinGroup(ctx context.Context, req *pbGroup.JoinGroupReq) (*pbGroup.JoinGroupResp, error) {
	resp := &pbGroup.JoinGroupResp{}

	if _, err := relation.GetUserByUserID(tracelog.GetOpUserID(ctx)); err != nil {
		return nil, err
	}
	groupInfo, err := rocksCache.GetGroupInfoFromCache(ctx, req.GroupID)
	if err != nil {
		return nil, err
	}
	if groupInfo.Status == constant.GroupStatusDismissed {
		return nil, utils.Wrap(constant.ErrDismissedAlready, "")
	}

	if groupInfo.NeedVerification == constant.Directly {
		if groupInfo.GroupType != constant.SuperGroup {
			us, err := relation.GetUserByUserID(tracelog.GetOpUserID(ctx))
			if err != nil {
				return nil, err
			}
			//to group member
			groupMember := table.GroupMemberModel{GroupID: req.GroupID, RoleLevel: constant.GroupOrdinaryUsers, OperatorUserID: tracelog.GetOpUserID(ctx)}
			utils.CopyStructFields(&groupMember, us)
			if err := CallbackBeforeMemberJoinGroup(ctx, tracelog.GetOperationID(ctx), &groupMember, groupInfo.Ex); err != nil {
				return nil, err
			}
			if err := s.DelGroupAndUserCache(ctx, req.GroupID, []string{tracelog.GetOpUserID(ctx)}); err != nil {
				return nil, err
			}
			err = relation.InsertIntoGroupMember(groupMember)
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
			reqPb.OperationID = tracelog.GetOperationID(ctx)
			c.OwnerUserID = tracelog.GetOpUserID(ctx)
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
			tracelog.SetCtxInfo(ctx, "SetConversation", err, "req", reqPb, "resp", respPb)
			chat.MemberEnterDirectlyNotification(req.GroupID, tracelog.GetOpUserID(ctx), tracelog.GetOperationID(ctx))
			return resp, nil
		} else {
			constant.SetErrorForResp(constant.ErrGroupTypeNotSupport, resp.CommonResp)
			return resp, nil
		}
	}
	var groupRequest table.GroupRequestModel
	groupRequest.UserID = tracelog.GetOpUserID(ctx)
	groupRequest.ReqMsg = req.ReqMessage
	groupRequest.GroupID = req.GroupID
	groupRequest.JoinSource = req.JoinSource
	err = relation.InsertIntoGroupRequest(groupRequest)
	if err != nil {
		return nil, err
	}
	chat.JoinGroupApplicationNotification(ctx, req)
	return resp, nil
}

func (s *groupServer) QuitGroup(ctx context.Context, req *pbGroup.QuitGroupReq) (*pbGroup.QuitGroupResp, error) {
	resp := &pbGroup.QuitGroupResp{}

	groupInfo, err := relation.GetGroupInfoByGroupID(req.GroupID)
	if err != nil {
		return nil, err
	}
	if groupInfo.GroupType != constant.SuperGroup {
		_, err = rocksCache.GetGroupMemberInfoFromCache(ctx, req.GroupID, tracelog.GetOpUserID(ctx))
		if err != nil {
			return nil, err
		}
		if err := s.DelGroupAndUserCache(ctx, req.GroupID, []string{tracelog.GetOpUserID(ctx)}); err != nil {
			return nil, err
		}
		err = relation.DeleteGroupMemberByGroupIDAndUserID(req.GroupID, tracelog.GetOpUserID(ctx))
		if err != nil {
			return nil, err
		}
	} else {
		okUserIDList := []string{tracelog.GetOpUserID(ctx)}
		if err := db.DB.RemoverUserFromSuperGroup(req.GroupID, okUserIDList); err != nil {
			return nil, err
		}
	}

	if groupInfo.GroupType != constant.SuperGroup {
		_ = rocksCache.DelGroupMemberInfoFromCache(ctx, req.GroupID, tracelog.GetOpUserID(ctx))
		chat.MemberQuitNotification(req)
	} else {
		_ = rocksCache.DelJoinedSuperGroupIDListFromCache(ctx, tracelog.GetOpUserID(ctx))
		_ = rocksCache.DelGroupMemberListHashFromCache(ctx, req.GroupID)
		chat.SuperGroupNotification(tracelog.GetOperationID(ctx), tracelog.GetOpUserID(ctx), tracelog.GetOpUserID(ctx))
	}
	return resp, nil
}

func (s *groupServer) SetGroupInfo(ctx context.Context, req *pbGroup.SetGroupInfoReq) (*pbGroup.SetGroupInfoResp, error) {
	resp := &pbGroup.SetGroupInfoResp{}

	if !hasAccess(req) {
		return nil, utils.Wrap(constant.ErrIdentity, "")
	}
	group, err := relation.GetGroupInfoByGroupID(req.GroupInfoForSet.GroupID)
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
		if err := relation.UpdateGroupInfoDefaultZero(req.GroupInfoForSet.GroupID, m); err != nil {
			return nil, err
		}
	}
	if req.GroupInfoForSet.LookMemberInfo != nil {
		changedType = changedType | (1 << 5)
		m := make(map[string]interface{})
		m["look_member_info"] = req.GroupInfoForSet.LookMemberInfo.Value
		if err := relation.UpdateGroupInfoDefaultZero(req.GroupInfoForSet.GroupID, m); err != nil {
			return nil, err
		}
	}
	if req.GroupInfoForSet.ApplyMemberFriend != nil {
		changedType = changedType | (1 << 6)
		m := make(map[string]interface{})
		m["apply_member_friend"] = req.GroupInfoForSet.ApplyMemberFriend.Value
		if err := relation.UpdateGroupInfoDefaultZero(req.GroupInfoForSet.GroupID, m); err != nil {
			return nil, err
		}
	}
	//only administrators can set group information
	var groupInfo table.GroupModel
	utils.CopyStructFields(&groupInfo, req.GroupInfoForSet)
	if req.GroupInfoForSet.Notification != "" {
		groupInfo.NotificationUserID = tracelog.GetOpUserID(ctx)
		groupInfo.NotificationUpdateTime = time.Now()
	}
	if err := rocksCache.DelGroupInfoFromCache(ctx, req.GroupInfoForSet.GroupID); err != nil {
		return nil, err
	}
	err = relation.SetGroupInfo(groupInfo)
	if err != nil {
		return nil, err
	}
	if changedType != 0 {
		chat.GroupInfoSetNotification(tracelog.GetOperationID(ctx), tracelog.GetOpUserID(ctx), req.GroupInfoForSet.GroupID, groupName, notification, introduction, faceURL, req.GroupInfoForSet.NeedVerification)
	}
	if req.GroupInfoForSet.Notification != "" {
		//get group member user id
		var conversationReq pbConversation.ModifyConversationFieldReq
		conversation := pbConversation.Conversation{
			OwnerUserID:      tracelog.GetOpUserID(ctx),
			ConversationID:   utils.GetConversationIDBySessionType(req.GroupInfoForSet.GroupID, constant.GroupChatType),
			ConversationType: constant.GroupChatType,
			GroupID:          req.GroupInfoForSet.GroupID,
		}
		conversationReq.Conversation = &conversation
		conversationReq.OperationID = tracelog.GetOperationID(ctx)
		conversationReq.FieldType = constant.FieldGroupAtType
		conversation.GroupAtType = constant.GroupNotification
		conversationReq.UserIDList = cacheResp.UserIDList

		_, err := s.conversationRpc.ModifyConversationField(ctx, &conversationReq)
		tracelog.SetCtxInfo(ctx, "ModifyConversationField", err, "req", &conversationReq, "resp", conversationReply)
	}
	return resp, nil
}

func (s *groupServer) TransferGroupOwner(ctx context.Context, req *pbGroup.TransferGroupOwnerReq) (*pbGroup.TransferGroupOwnerResp, error) {
	resp := &pbGroup.TransferGroupOwnerResp{}

	groupInfo, err := relation.GetGroupInfoByGroupID(req.GroupID)
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

	groupMemberInfo := table.GroupMemberModel{GroupID: req.GroupID, UserID: req.OldOwnerUserID, RoleLevel: constant.GroupOrdinaryUsers}
	err = relation.UpdateGroupMemberInfo(groupMemberInfo)
	if err != nil {
		return nil, err
	}
	groupMemberInfo = table.GroupMemberModel{GroupID: req.GroupID, UserID: req.NewOwnerUserID, RoleLevel: constant.GroupOwner}
	err = relation.UpdateGroupMemberInfo(groupMemberInfo)
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
		groupInfoDB, err := relation.GetGroupInfoByGroupID(req.GroupID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return resp, nil
			}
			return nil, err
		}
		resp.GroupNum = 1
		groupInfo := &open_im_sdk.GroupInfo{}
		utils.CopyStructFields(groupInfo, groupInfoDB)
		groupMember, err := relation.GetGroupOwnerInfoByGroupID(req.GroupID)
		if err != nil {
			return nil, err
		}
		memberNum, err := relation.GetGroupMembersCount(req.GroupID, "")
		if err != nil {
			return nil, err
		}
		groupInfo.MemberCount = uint32(memberNum)
		groupInfo.CreateTime = uint32(groupInfoDB.CreateTime.Unix())
		resp.Groups = append(resp.Groups, &pbGroup.CMSGroup{GroupInfo: groupInfo, GroupOwnerUserName: groupMember.Nickname, GroupOwnerUserID: groupMember.UserID})
	} else {
		groups, count, err := relation.GetGroupsByName(req.GroupName, req.Pagination.PageNumber, req.Pagination.ShowNumber)
		if err != nil {
			tracelog.SetCtxInfo(ctx, "GetGroupsByName", err, "GroupName", req.GroupName, "PageNumber", req.Pagination.PageNumber, "ShowNumber", req.Pagination.ShowNumber)
		}
		for _, v := range groups {
			group := &pbGroup.CMSGroup{GroupInfo: &open_im_sdk.GroupInfo{}}
			utils.CopyStructFields(group.GroupInfo, v)
			groupMember, err := relation.GetGroupOwnerInfoByGroupID(v.GroupID)
			if err != nil {
				tracelog.SetCtxInfo(ctx, "GetGroupOwnerInfoByGroupID", err, "GroupID", v.GroupID)
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
	groupMemberInfo := table.GroupMemberModel{GroupID: req.GroupID, UserID: req.UserID}
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

	groupMemberInfo := table.GroupMemberModel{GroupID: req.GroupID, UserID: req.UserID}
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
	groupMemberInfo := table.GroupMemberModel{}
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
	groupMember := table.GroupMemberModel{
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
