package user

import (
	chat "Open_IM/internal/rpc/msg"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/controller"
	"Open_IM/pkg/common/db/relation"
	"Open_IM/pkg/common/log"
	promePkg "Open_IM/pkg/common/prometheus"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/getcdv3"
	pbConversation "Open_IM/pkg/proto/conversation"
	pbFriend "Open_IM/pkg/proto/friend"
	sdkws "Open_IM/pkg/proto/sdk_ws"
	pbUser "Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"net"
	"strconv"
	"strings"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	utils2 "Open_IM/internal/utils"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type userServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
	controller.UserInterface
}

func NewUserServer(port int) *userServer {
	log.NewPrivateLog(constant.LogFileName)
	u := userServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImUserName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
	//mysql init
	var mysql relation.Mysql
	var model relation.User
	err := mysql.InitConn().AutoMigrateModel(&model)
	if err != nil {
		panic("db init err:" + err.Error())
	}
	if mysql.GormConn() != nil {
		model.DB = mysql.GormConn()
	} else {
		panic("db init err:" + "conn is nil")
	}
	u.UserInterface = controller.NewUserController(model.DB)
	return &u
}

func (s *userServer) Run() {
	log.NewInfo("0", "rpc user start...")

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
	log.NewInfo("0", "listen network success, address ", address, listener)
	defer listener.Close()
	//grpc server
	var grpcOpts []grpc.ServerOption
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
	pbUser.RegisterUserServer(srv, s)
	rpcRegisterIP := config.Config.RpcRegisterIP
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error("", "GetLocalIP failed ", err.Error())
		}
	}
	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName, 10)
	if err != nil {
		log.NewError("0", "RegisterEtcd failed ", err.Error(), s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName)
		panic(utils.Wrap(err, "register user module  rpc to etcd err"))
	}
	err = srv.Serve(listener)
	if err != nil {
		log.NewError("0", "Serve failed ", err.Error())
		return
	}
	log.NewInfo("0", "rpc  user success")
}

func syncPeerUserConversation(conversation *pbConversation.Conversation, operationID string) error {
	peerUserConversation := imdb.Conversation{
		OwnerUserID:      conversation.UserID,
		ConversationID:   utils.GetConversationIDBySessionType(conversation.OwnerUserID, constant.SingleChatType),
		ConversationType: constant.SingleChatType,
		UserID:           conversation.OwnerUserID,
		GroupID:          "",
		RecvMsgOpt:       0,
		UnreadCount:      0,
		DraftTextTime:    0,
		IsPinned:         false,
		IsPrivateChat:    conversation.IsPrivateChat,
		AttachedInfo:     "",
		Ex:               "",
	}
	err := imdb.PeerUserSetConversation(peerUserConversation)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "SetConversation error", err.Error())
		return err
	}
	chat.ConversationSetPrivateNotification(operationID, conversation.OwnerUserID, conversation.UserID, conversation.IsPrivateChat)
	return nil
}

func (s *userServer) GetUsersInfo(ctx context.Context, req *pbUser.GetUsersInfoReq) (*pbUser.GetUsersInfoResp, error) {
	resp := &pbUser.GetUsersInfoResp{}
	users, err := s.Find(ctx, req.UserIDs)
	if err != nil {
		return nil, err
	}
	for _, v := range users {
		n, err := utils2.NewDBUser(v).Convert()
		if err != nil {
			return nil, err
		}
		resp.UsersInfo = append(resp.UsersInfo, n)
	}
	return resp, nil
}

func (s *userServer) BatchSetConversations(ctx context.Context, req *pbUser.BatchSetConversationsReq) (*pbUser.BatchSetConversationsResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	if req.NotificationType == 0 {
		req.NotificationType = constant.ConversationOptChangeNotification
	}
	resp := &pbUser.BatchSetConversationsResp{}
	for _, v := range req.Conversations {
		conversation := imdb.Conversation{}
		if err := utils.CopyStructFields(&conversation, v); err != nil {
			log.NewDebug(req.OperationID, utils.GetSelfFuncName(), v.String(), "CopyStructFields failed", err.Error())
		}
		//redis op
		if err := db.DB.SetSingleConversationRecvMsgOpt(req.OwnerUserID, v.ConversationID, v.RecvMsgOpt); err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "cache failed, rpc return", err.Error())
			resp.CommonResp = &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
			return resp, nil
		}

		isUpdate, err := imdb.SetConversation(conversation)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "SetConversation error", err.Error())
			resp.Failed = append(resp.Failed, v.ConversationID)
			continue
		}
		if isUpdate {
			err = rocksCache.DelConversationFromCache(v.OwnerUserID, v.ConversationID)
		} else {
			err = rocksCache.DelUserConversationIDListFromCache(v.OwnerUserID)
		}
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), v.ConversationID, v.OwnerUserID)
		}
		resp.Success = append(resp.Success, v.ConversationID)
		// if is set private msg operation，then peer user need to sync and set tips\
		if v.ConversationType == constant.SingleChatType && req.NotificationType == constant.ConversationPrivateChatNotification {
			if err := syncPeerUserConversation(v, req.OperationID); err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "syncPeerUserConversation", err.Error())
			}
		}
	}
	chat.ConversationChangeNotification(req.OperationID, req.OwnerUserID)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc return", resp.String())
	resp.CommonResp = &pbUser.CommonResp{}
	return resp, nil
}

func (s *userServer) GetAllConversations(ctx context.Context, req *pbUser.GetAllConversationsReq) (*pbUser.GetAllConversationsResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &pbUser.GetAllConversationsResp{Conversations: []*pbConversation.Conversation{}}
	conversations, err := rocksCache.GetUserAllConversationList(req.OwnerUserID)
	log.NewDebug(req.OperationID, "conversations: ", conversations)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetConversations error", err.Error())
		resp.CommonResp = &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
		return resp, nil
	}
	if err = utils.CopyStructFields(&resp.Conversations, conversations); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields error", err.Error())
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc return", resp.String())
	resp.CommonResp = &pbUser.CommonResp{}
	return resp, nil
}

func (s *userServer) GetConversation(ctx context.Context, req *pbUser.GetConversationReq) (*pbUser.GetConversationResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &pbUser.GetConversationResp{Conversation: &pbConversation.Conversation{}}
	conversation, err := rocksCache.GetConversationFromCache(req.OwnerUserID, req.ConversationID)
	log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "conversation", conversation)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetConversation error", err.Error())
		resp.CommonResp = &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
		return resp, nil
	}
	if err := utils.CopyStructFields(resp.Conversation, &conversation); err != nil {
		log.Debug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields error", conversation, err.Error())
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	resp.CommonResp = &pbUser.CommonResp{}
	return resp, nil
}

func (s *userServer) GetConversations(ctx context.Context, req *pbUser.GetConversationsReq) (*pbUser.GetConversationsResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &pbUser.GetConversationsResp{Conversations: []*pbConversation.Conversation{}}
	conversations, err := rocksCache.GetConversationsFromCache(req.OwnerUserID, req.ConversationIDs)
	log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "conversations", conversations)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetConversations error", err.Error())
		resp.CommonResp = &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
		return resp, nil
	}
	if err := utils.CopyStructFields(&resp.Conversations, conversations); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", conversations, err.Error())
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	resp.CommonResp = &pbUser.CommonResp{}
	return resp, nil
}

func (s *userServer) SetConversation(ctx context.Context, req *pbUser.SetConversationReq) (*pbUser.SetConversationResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &pbUser.SetConversationResp{}
	if req.NotificationType == 0 {
		req.NotificationType = constant.ConversationOptChangeNotification
	}
	if req.Conversation.ConversationType == constant.GroupChatType {
		groupInfo, err := imdb.GetGroupInfoByGroupID(req.Conversation.GroupID)
		if err != nil {
			log.NewError(req.OperationID, "GetGroupInfoByGroupID failed ", req.Conversation.GroupID, err.Error())
			resp.CommonResp = &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
			return resp, nil
		}
		if groupInfo.Status == constant.GroupStatusDismissed && !req.Conversation.IsNotInGroup {
			log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "group status is dismissed", groupInfo)
			errMsg := "group status is dismissed"
			resp.CommonResp = &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}
			return resp, nil
		}
	}
	var conversation imdb.Conversation
	if err := utils.CopyStructFields(&conversation, req.Conversation); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", *req.Conversation, err.Error())
	}
	if err := db.DB.SetSingleConversationRecvMsgOpt(req.Conversation.OwnerUserID, req.Conversation.ConversationID, req.Conversation.RecvMsgOpt); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "cache failed, rpc return", err.Error())
		resp.CommonResp = &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
		return resp, nil
	}
	isUpdate, err := imdb.SetConversation(conversation)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "SetConversation error", err.Error())
		resp.CommonResp = &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
		return resp, nil
	}
	if isUpdate {
		err = rocksCache.DelConversationFromCache(req.Conversation.OwnerUserID, req.Conversation.ConversationID)
	} else {
		err = rocksCache.DelUserConversationIDListFromCache(req.Conversation.OwnerUserID)
	}
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.Conversation.ConversationID, req.Conversation.OwnerUserID)
	}

	// notification
	if req.Conversation.ConversationType == constant.SingleChatType && req.NotificationType == constant.ConversationPrivateChatNotification {
		//sync peer user conversation if conversation is singleChatType
		if err := syncPeerUserConversation(req.Conversation, req.OperationID); err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "syncPeerUserConversation", err.Error())
			resp.CommonResp = &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
			return resp, nil
		}
	} else {
		chat.ConversationChangeNotification(req.OperationID, req.Conversation.OwnerUserID)
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc return", resp.String())
	resp.CommonResp = &pbUser.CommonResp{}
	return resp, nil
}

func (s *userServer) SetRecvMsgOpt(ctx context.Context, req *pbUser.SetRecvMsgOptReq) (*pbUser.SetRecvMsgOptResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &pbUser.SetRecvMsgOptResp{}
	var conversation imdb.Conversation
	if err := utils.CopyStructFields(&conversation, req); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", *req, err.Error())
	}
	if err := db.DB.SetSingleConversationRecvMsgOpt(req.OwnerUserID, req.ConversationID, req.RecvMsgOpt); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "cache failed, rpc return", err.Error())
		resp.CommonResp = &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
		return resp, nil
	}
	stringList := strings.Split(req.ConversationID, "_")
	if len(stringList) > 1 {
		switch stringList[0] {
		case "single":
			conversation.UserID = stringList[1]
			conversation.ConversationType = constant.SingleChatType
		case "group":
			conversation.GroupID = stringList[1]
			conversation.ConversationType = constant.GroupChatType
		}
	}
	isUpdate, err := imdb.SetRecvMsgOpt(conversation)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "SetConversation error", err.Error())
		resp.CommonResp = &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
		return resp, nil
	}
	if isUpdate {
		err = rocksCache.DelConversationFromCache(conversation.OwnerUserID, conversation.ConversationID)
	} else {
		err = rocksCache.DelUserConversationIDListFromCache(conversation.OwnerUserID)
	}
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), conversation.ConversationID, err.Error())
	}
	chat.ConversationChangeNotification(req.OperationID, req.OwnerUserID)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	resp.CommonResp = &pbUser.CommonResp{}
	return resp, nil
}

func (s *userServer) GetAllUserID(_ context.Context, req *pbUser.GetAllUserIDReq) (*pbUser.GetAllUserIDResp, error) {
	log.NewInfo(req.OperationID, "GetAllUserID args ", req.String())
	if !token_verify.IsManagerUserID(req.OpUserID) {
		log.NewError(req.OperationID, "IsManagerUserID false ", req.OpUserID)
		return &pbUser.GetAllUserIDResp{CommonResp: &pbUser.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}}, nil
	}
	uidList, err := imdb.SelectAllUserID()
	if err != nil {
		log.NewError(req.OperationID, "SelectAllUserID false ", err.Error())
		return &pbUser.GetAllUserIDResp{CommonResp: &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	} else {
		log.NewInfo(req.OperationID, "GetAllUserID rpc return ", pbUser.GetAllUserIDResp{CommonResp: &pbUser.CommonResp{}, UserIDList: uidList})
		return &pbUser.GetAllUserIDResp{CommonResp: &pbUser.CommonResp{}, UserIDList: uidList}, nil
	}
}

func (s *userServer) AccountCheck(_ context.Context, req *pbUser.AccountCheckReq) (*pbUser.AccountCheckResp, error) {
	log.NewInfo(req.OperationID, "AccountCheck args ", req.String())
	if !token_verify.IsManagerUserID(req.OpUserID) {
		log.NewError(req.OperationID, "IsManagerUserID false ", req.OpUserID)
		return &pbUser.AccountCheckResp{CommonResp: &pbUser.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}}, nil
	}
	uidList, err := imdb.SelectSomeUserID(req.CheckUserIDList)
	log.NewDebug(req.OperationID, "from db uid list is:", uidList)
	if err != nil {
		log.NewError(req.OperationID, "SelectSomeUserID failed ", err.Error(), req.CheckUserIDList)
		return &pbUser.AccountCheckResp{CommonResp: &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	} else {
		var r []*pbUser.AccountCheckResp_SingleUserStatus
		for _, v := range req.CheckUserIDList {
			temp := new(pbUser.AccountCheckResp_SingleUserStatus)
			temp.UserID = v
			if utils.IsContain(v, uidList) {
				temp.AccountStatus = constant.Registered
			} else {
				temp.AccountStatus = constant.UnRegistered
			}
			r = append(r, temp)
		}
		resp := pbUser.AccountCheckResp{CommonResp: &pbUser.CommonResp{ErrCode: 0, ErrMsg: ""}, ResultList: r}
		log.NewInfo(req.OperationID, "AccountCheck rpc return ", resp.String())
		return &resp, nil
	}

}

func (s *userServer) UpdateUserInfo(ctx context.Context, req *pbUser.UpdateUserInfoReq) (*pbUser.UpdateUserInfoResp, error) {
	resp := pbUser.UpdateUserInfoResp{}
	err := token_verify.CheckAccessV3(ctx, req.UserInfo.UserID)
	if err != nil {
		return nil, err
	}

	oldNickname := ""
	if req.UserInfo.Nickname != "" {
		u, err := s.Take(ctx, req.UserInfo.UserID)
		if err != nil {
			return nil, err
		}
		oldNickname = u.Nickname
	}

	user, err := utils2.NewPBUser(req.UserInfo).Convert()
	if err != nil {
		return nil, err
	}
	err = s.Update(ctx, []*relation.User{user})
	if err != nil {
		return nil, err
	}
	etcdConn, err := getcdv3.GetConn(ctx, config.Config.RpcRegisterName.OpenImFriendName)
	if err != nil {
		return nil, err
	}

	client := pbFriend.NewFriendClient(etcdConn)
	newReq := &pbFriend.GetFriendListReq{
		CommID: &pbFriend.CommID{OperationID: req.OperationID, FromUserID: req.UserInfo.UserID, OpUserID: req.OpUserID},
	}

	rpcResp, err := client.GetFriendList(context.Background(), newReq)
	if err != nil {
		log.NewError(req.OperationID, "GetFriendList failed ", err.Error(), newReq)
		return &pbUser.UpdateUserInfoResp{CommonResp: &pbUser.CommonResp{ErrCode: 500, ErrMsg: err.Error()}}, nil
	}
	for _, v := range rpcResp.FriendInfoList {
		log.Info(req.OperationID, "UserInfoUpdatedNotification ", req.UserInfo.UserID, v.FriendUser.UserID)
		//	chat.UserInfoUpdatedNotification(req.OperationID, req.UserInfo.UserID, v.FriendUser.UserID)
		chat.FriendInfoUpdatedNotification(req.OperationID, req.UserInfo.UserID, v.FriendUser.UserID, req.OpUserID)
	}
	if err := rocksCache.DelUserInfoFromCache(user.UserID); err != nil {
		log.NewError(req.OperationID, "GetFriendList failed ", err.Error(), newReq)
		return &pbUser.UpdateUserInfoResp{CommonResp: &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}
	//chat.UserInfoUpdatedNotification(req.OperationID, req.OpUserID, req.UserInfo.UserID)
	chat.UserInfoUpdatedNotification(req.OperationID, req.OpUserID, req.UserInfo.UserID)
	log.Info(req.OperationID, "UserInfoUpdatedNotification ", req.UserInfo.UserID, req.OpUserID)
	if req.UserInfo.FaceURL != "" {
		s.SyncJoinedGroupMemberFaceURL(req.UserInfo.UserID, req.UserInfo.FaceURL, req.OperationID, req.OpUserID)
	}
	if req.UserInfo.Nickname != "" {
		s.SyncJoinedGroupMemberNickname(req.UserInfo.UserID, req.UserInfo.Nickname, oldNickname, req.OperationID, req.OpUserID)
	}
	return &pbUser.UpdateUserInfoResp{CommonResp: &pbUser.CommonResp{}}, nil
}
func (s *userServer) SetGlobalRecvMessageOpt(ctx context.Context, req *pbUser.SetGlobalRecvMessageOptReq) (*pbUser.SetGlobalRecvMessageOptResp, error) {
	log.NewInfo(req.OperationID, "SetGlobalRecvMessageOpt args ", req.String())

	var user imdb.User
	user.UserID = req.UserID
	m := make(map[string]interface{}, 1)

	log.NewDebug(req.OperationID, utils.GetSelfFuncName(), req.GlobalRecvMsgOpt, "set GlobalRecvMsgOpt")
	m["global_recv_msg_opt"] = req.GlobalRecvMsgOpt
	err := db.DB.SetUserGlobalMsgRecvOpt(user.UserID, req.GlobalRecvMsgOpt)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "SetGlobalRecvMessageOpt failed ", err.Error(), user)
		return &pbUser.SetGlobalRecvMessageOptResp{CommonResp: &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	err = imdb.UpdateUserInfoByMap(user, m)
	if err != nil {
		log.NewError(req.OperationID, "SetGlobalRecvMessageOpt failed ", err.Error(), user)
		return &pbUser.SetGlobalRecvMessageOptResp{CommonResp: &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	if err := rocksCache.DelUserInfoFromCache(user.UserID); err != nil {
		log.NewError(req.OperationID, "DelUserInfoFromCache failed ", err.Error(), req.String())
		return &pbUser.SetGlobalRecvMessageOptResp{CommonResp: &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}
	chat.UserInfoUpdatedNotification(req.OperationID, req.UserID, req.UserID)
	return &pbUser.SetGlobalRecvMessageOptResp{CommonResp: &pbUser.CommonResp{}}, nil
}

func (s *userServer) SyncJoinedGroupMemberFaceURL(userID string, faceURL string, operationID string, opUserID string) {
	joinedGroupIDList, err := rocksCache.GetJoinedGroupIDListFromCache(userID)
	if err != nil {
		log.NewWarn(operationID, "GetJoinedGroupIDListByUserID failed ", userID, err.Error())
		return
	}
	for _, groupID := range joinedGroupIDList {
		groupMemberInfo := imdb.GroupMember{UserID: userID, GroupID: groupID, FaceURL: faceURL}
		if err := imdb.UpdateGroupMemberInfo(groupMemberInfo); err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), groupMemberInfo)
			continue
		}
		//if err := rocksCache.DelAllGroupMembersInfoFromCache(groupID); err != nil {
		//	log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), groupID)
		//	continue
		//}
		if err := rocksCache.DelGroupMemberInfoFromCache(groupID, userID); err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), groupID, userID)
			continue
		}
		chat.GroupMemberInfoSetNotification(operationID, opUserID, groupID, userID)
	}
}

func (s *userServer) SyncJoinedGroupMemberNickname(userID string, newNickname, oldNickname string, operationID string, opUserID string) {
	joinedGroupIDList, err := imdb.GetJoinedGroupIDListByUserID(userID)
	if err != nil {
		log.NewWarn(operationID, "GetJoinedGroupIDListByUserID failed ", userID, err.Error())
		return
	}
	for _, v := range joinedGroupIDList {
		member, err := imdb.GetGroupMemberInfoByGroupIDAndUserID(v, userID)
		if err != nil {
			log.NewWarn(operationID, "GetGroupMemberInfoByGroupIDAndUserID failed ", err.Error(), v, userID)
			continue
		}
		if member.Nickname == oldNickname {
			groupMemberInfo := imdb.GroupMember{UserID: userID, GroupID: v, Nickname: newNickname}
			if err := imdb.UpdateGroupMemberInfo(groupMemberInfo); err != nil {
				log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), groupMemberInfo)
				continue
			}
			//if err := rocksCache.DelAllGroupMembersInfoFromCache(v); err != nil {
			//	log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), v)
			//	continue
			//}
			if err := rocksCache.DelGroupMemberInfoFromCache(v, userID); err != nil {
				log.NewError(operationID, utils.GetSelfFuncName(), err.Error(), v)
			}
			chat.GroupMemberInfoSetNotification(operationID, opUserID, v, userID)
		}
	}
}

func (s *userServer) GetUsers(ctx context.Context, req *pbUser.GetUsersReq) (*pbUser.GetUsersResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	var usersDB []imdb.User
	var err error
	resp := &pbUser.GetUsersResp{CommonResp: &pbUser.CommonResp{}, Pagination: &sdkws.ResponsePagination{CurrentPage: req.Pagination.PageNumber, ShowNumber: req.Pagination.ShowNumber}}
	if req.UserID != "" {
		userDB, err := imdb.GetUserByUserID(req.UserID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return resp, nil
			}
			log.NewError(req.OperationID, utils.GetSelfFuncName(), req.UserID, err.Error())
			resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg
			return resp, nil
		}
		usersDB = append(usersDB, *userDB)
		resp.TotalNums = 1
	} else if req.UserName != "" {
		usersDB, err = imdb.GetUserByName(req.UserName, req.Pagination.ShowNumber, req.Pagination.PageNumber)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), req.UserName, req.Pagination.ShowNumber, req.Pagination.PageNumber, err.Error())
			resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg
			return resp, nil
		}
		resp.TotalNums, err = imdb.GetUsersCount(req.UserName)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), req.UserName, err.Error())
			resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resp.CommonResp.ErrMsg = err.Error()
			return resp, nil
		}

	} else if req.Content != "" {
		var count int64
		usersDB, count, err = imdb.GetUsersByNameAndID(req.Content, req.Pagination.ShowNumber, req.Pagination.PageNumber)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUsers failed", req.Pagination.ShowNumber, req.Pagination.PageNumber, err.Error())
			resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resp.CommonResp.ErrMsg = err.Error()
			return resp, nil
		}
		resp.TotalNums = int32(count)
	} else {
		usersDB, err = imdb.GetUsers(req.Pagination.ShowNumber, req.Pagination.PageNumber)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUsers failed", req.Pagination.ShowNumber, req.Pagination.PageNumber, err.Error())
			resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resp.CommonResp.ErrMsg = err.Error()
			return resp, nil
		}
		resp.TotalNums, err = imdb.GetTotalUserNum()
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error())
			resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
			resp.CommonResp.ErrMsg = err.Error()
			return resp, nil
		}
	}
	for _, userDB := range usersDB {
		var user sdkws.UserInfo
		utils.CopyStructFields(&user, userDB)
		user.CreateTime = uint32(userDB.CreateTime.Unix())
		user.BirthStr = utils.TimeToString(userDB.Birth)
		resp.UserList = append(resp.UserList, &pbUser.CmsUser{User: &user})
	}

	var userIDList []string
	for _, v := range resp.UserList {
		userIDList = append(userIDList, v.User.UserID)
	}
	isBlockUser, err := imdb.UsersIsBlock(userIDList)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), userIDList)
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}

	for _, v := range resp.UserList {
		if utils.IsContain(v.User.UserID, isBlockUser) {
			v.IsBlock = true
		}
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

func (s *userServer) AddUser(ctx context.Context, req *pbUser.AddUserReq) (*pbUser.AddUserResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &pbUser.AddUserResp{CommonResp: &pbUser.CommonResp{}}
	err := imdb.AddUser(req.UserInfo.UserID, req.UserInfo.PhoneNumber, req.UserInfo.Nickname, req.UserInfo.Email, req.UserInfo.Gender, req.UserInfo.FaceURL, req.UserInfo.BirthStr)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "AddUser", err.Error(), req.String())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	return resp, nil
}
