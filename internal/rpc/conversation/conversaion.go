package conversation

import (
	"Open_IM/internal/common/check"
	chat "Open_IM/internal/rpc/msg"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/cache"
	"Open_IM/pkg/common/db/controller"
	"Open_IM/pkg/common/db/relation"
	"Open_IM/pkg/common/db/table"
	"Open_IM/pkg/common/db/unrelation"
	"Open_IM/pkg/common/log"
	promePkg "Open_IM/pkg/common/prometheus"
	"Open_IM/pkg/getcdv3"
	pbConversation "Open_IM/pkg/proto/conversation"
	pbUser "Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"
	"context"
	"github.com/dtm-labs/rockscache"
	"net"
	"strconv"
	"strings"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"Open_IM/pkg/common/config"

	"google.golang.org/grpc"
)

type conversationServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
	groupChecker    *check.GroupChecker
	controller.ConversationInterface
}

func NewConversationServer(port int) *conversationServer {
	log.NewPrivateLog(constant.LogFileName)
	c := conversationServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImConversationName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
		groupChecker:    check.NewGroupChecker(),
	}
	var cDB relation.Conversation
	var cCache cache.ConversationCache
	//mysql init
	var mysql relation.Mysql
	err := mysql.InitConn().AutoMigrateModel(&table.ConversationModel{})
	if err != nil {
		panic("db init err:" + err.Error())
	}
	if mysql.GormConn() != nil {
		//get gorm model
		cDB = relation.NewConversationGorm(mysql.GormConn())
	} else {
		panic("db init err:" + "conn is nil")
	}
	//redis init
	var redis cache.RedisClient
	redis.InitRedis()
	rcClient := rockscache.NewClient(redis.GetClient(), rockscache.Options{
		RandomExpireAdjustment: 0.2,
		DisableCacheRead:       false,
		DisableCacheDelete:     false,
		StrongConsistency:      true,
	})
	cCache = cache.NewConversationRedis(rcClient)

	database := controller.NewConversationDataBase(cDB, cCache)
	c.ConversationInterface = controller.NewConversationController(database)
	return &c
}

func (c *conversationServer) Run() {
	log.NewInfo("0", "rpc conversation start...")

	listenIP := ""
	if config.Config.ListenIP == "" {
		listenIP = "0.0.0.0"
	} else {
		listenIP = config.Config.ListenIP
	}
	address := listenIP + ":" + strconv.Itoa(c.rpcPort)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic("listening err:" + err.Error() + c.rpcRegisterName)
	}
	log.NewInfo("0", "listen network success, ", address, listener)
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

	//service registers with etcd
	pbConversation.RegisterConversationServer(srv, c)
	rpcRegisterIP := config.Config.RpcRegisterIP
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error("", "GetLocalIP failed ", err.Error())
		}
	}
	log.NewInfo("", "rpcRegisterIP", rpcRegisterIP)
	err = getcdv3.RegisterEtcd(c.etcdSchema, strings.Join(c.etcdAddr, ","), rpcRegisterIP, c.rpcPort, c.rpcRegisterName, 10, "")
	if err != nil {
		log.NewError("0", "RegisterEtcd failed ", err.Error(),
			c.etcdSchema, strings.Join(c.etcdAddr, ","), rpcRegisterIP, c.rpcPort, c.rpcRegisterName)
		panic(utils.Wrap(err, "register conversation module  rpc to etcd err"))
	}
	log.NewInfo("0", "RegisterConversationServer ok ", c.etcdSchema, strings.Join(c.etcdAddr, ","), rpcRegisterIP, c.rpcPort, c.rpcRegisterName)
	err = srv.Serve(listener)
	if err != nil {
		log.NewError("0", "Serve failed ", err.Error())
		return
	}
	log.NewInfo("0", "rpc conversation ok")
}

func (c *conversationServer) GetConversation(ctx context.Context, req *pbConversation.GetConversationReq) (*pbConversation.GetConversationResp, error) {
	resp := &pbConversation.GetConversationResp{Conversation: &pbConversation.Conversation{}}
	conversations, err := c.ConversationInterface.FindConversations(ctx, req.OwnerUserID, []string{req.ConversationID})
	if err != nil {
		return nil, err
	}
	if len(conversations) > 0 {
		if err := utils.CopyStructFields(resp.Conversation, &conversations[0]); err != nil {
			return nil, err
		}
		return resp, nil
	}
	return nil, nil
}

func (c *conversationServer) GetAllConversations(ctx context.Context, req *pbConversation.GetAllConversationsReq) (*pbConversation.GetAllConversationsResp, error) {
	resp := &pbConversation.GetAllConversationsResp{Conversations: []*pbConversation.Conversation{}}
	conversations, err := c.ConversationInterface.GetUserAllConversation(ctx, req.OwnerUserID)
	if err != nil {
		return nil, err
	}
	if err := utils.CopyStructFields(&resp.Conversations, conversations); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *conversationServer) GetConversations(ctx context.Context, req *pbConversation.GetConversationsReq) (*pbConversation.GetConversationsResp, error) {
	resp := &pbConversation.GetConversationsResp{Conversations: []*pbConversation.Conversation{}}
	conversations, err := c.ConversationInterface.FindConversations(ctx, req.OwnerUserID, req.ConversationIDs)
	if err != nil {
		return nil, err
	}
	if err := utils.CopyStructFields(&resp.Conversations, conversations); err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *conversationServer) BatchSetConversations(ctx context.Context, req *pbConversation.BatchSetConversationsReq) (*pbConversation.BatchSetConversationsResp, error) {
	resp := &pbConversation.BatchSetConversationsResp{}
	var conversations []*table.ConversationModel
	if err := utils.CopyStructFields(&conversations, req.Conversations); err != nil {
		return nil, err
	}
	err := c.ConversationInterface.SetUserConversations(ctx, req.OwnerUserID, conversations)
	if err != nil {
		return nil, err
	}
	chat.ConversationChangeNotification(ctx, req.OwnerUserID)
	return resp, nil
}

func (c *conversationServer) SetConversation(ctx context.Context, req *pbConversation.SetConversationReq) (*pbConversation.SetConversationResp, error) {
	panic("implement me")
}

func (c *conversationServer) SetRecvMsgOpt(ctx context.Context, req *pbConversation.SetRecvMsgOptReq) (*pbConversation.SetRecvMsgOptResp, error) {
	panic("implement me")
}

func (c *conversationServer) ModifyConversationField(ctx context.Context, req *pbConversation.ModifyConversationFieldReq) (*pbConversation.ModifyConversationFieldResp, error) {
	resp := &pbConversation.ModifyConversationFieldResp{}
	var err error
	isSyncConversation := true
	if req.Conversation.ConversationType == constant.GroupChatType {
		groupInfo, err := c.groupChecker.GetGroupInfo(req.Conversation.GroupID)
		if err != nil {
			return nil, err
		}
		if groupInfo.Status == constant.GroupStatusDismissed && req.FieldType != constant.FieldUnread {
			return nil, err
		}
	}
	var conversation table.ConversationModel
	if err := utils.CopyStructFields(&conversation, req.Conversation); err != nil {
		return nil, err
	}
	if req.FieldType == constant.FieldIsPrivateChat {
		err := c.ConversationInterface.SyncPeerUserPrivateConversationTx(ctx, req.Conversation)
		if err != nil {
			return nil, err
		}
		chat.ConversationSetPrivateNotification(req.OperationID, req.Conversation.OwnerUserID, req.Conversation.UserID, req.Conversation.IsPrivateChat)
		return resp, nil
	}
	//haveUserID, err := c.ConversationInterface.GetUserIDExistConversation(ctx, req.UserIDList, req.Conversation.ConversationID)
	//if err != nil {
	//	return nil, err
	//}
	filedMap := make(map[string]interface{})
	switch req.FieldType {
	case constant.FieldRecvMsgOpt:
		filedMap["recv_msg_opt"] = req.Conversation.RecvMsgOpt
	case constant.FieldGroupAtType:
		filedMap["group_at_type"] = req.Conversation.GroupAtType
	case constant.FieldIsNotInGroup:
		filedMap["is_not_in_group"] = req.Conversation.IsNotInGroup
	case constant.FieldIsPinned:
		filedMap["is_pinned"] = req.Conversation.IsPinned
	case constant.FieldEx:
		filedMap["ex"] = req.Conversation.Ex
	case constant.FieldAttachedInfo:
		filedMap["attached_info"] = req.Conversation.AttachedInfo
	case constant.FieldUnread:
		isSyncConversation = false
		filedMap["update_unread_count_time"] = req.Conversation.UpdateUnreadCountTime
	case constant.FieldBurnDuration:
		filedMap["burn_duration"] = req.Conversation.BurnDuration
	}
	c.ConversationInterface.SetUsersConversationFiledTx(ctx, req.UserIDList, &conversation, filedMap)
	err = c.ConversationInterface.UpdateUsersConversationFiled(ctx, haveUserID, req.Conversation.ConversationID, filedMap)
	if err != nil {
		return nil, err
	}
	var conversations []*pbConversation.Conversation
	for _, v := range utils.DifferenceString(haveUserID, req.UserIDList) {
		temp := new(pbConversation.Conversation)
		_ = utils.CopyStructFields(temp, req.Conversation)
		temp.OwnerUserID = v
		conversations = append(conversations, temp)
	}
	err = c.ConversationInterface.CreateConversation(ctx, conversations)
	if err != nil {
		return nil, err
	}
	if isSyncConversation {
		for _, v := range req.UserIDList {
			chat.ConversationChangeNotification(req.OperationID, v)
		}
	} else {
		for _, v := range req.UserIDList {
			chat.ConversationUnreadChangeNotification(req.OperationID, v, req.Conversation.ConversationID, req.Conversation.UpdateUnreadCountTime)
		}
	}
	return resp, nil
}
