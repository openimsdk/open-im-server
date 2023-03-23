package conversation

import (
	chat "Open_IM/internal/rpc/msg"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	rocksCache "Open_IM/pkg/common/db/rocks_cache"
	"Open_IM/pkg/common/log"
	promePkg "Open_IM/pkg/common/prometheus"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbConversation "Open_IM/pkg/proto/conversation"
	"Open_IM/pkg/utils"
	"context"
	"net"
	"strconv"
	"strings"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"Open_IM/pkg/common/config"

	"google.golang.org/grpc"
)

type rpcConversation struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func (rpc *rpcConversation) ModifyConversationField(c context.Context, req *pbConversation.ModifyConversationFieldReq) (*pbConversation.ModifyConversationFieldResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &pbConversation.ModifyConversationFieldResp{}
	var err error
	isSyncConversation := true
	if req.Conversation.ConversationType == constant.GroupChatType {
		groupInfo, err := imdb.GetGroupInfoByGroupID(req.Conversation.GroupID)
		if err != nil {
			log.NewError(req.OperationID, "GetGroupInfoByGroupID failed ", req.Conversation.GroupID, err.Error())
			resp.CommonResp = &pbConversation.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
			return resp, nil
		}
		if groupInfo.Status == constant.GroupStatusDismissed && !req.Conversation.IsNotInGroup && req.FieldType != constant.FieldUnread {
			errMsg := "group status is dismissed"
			resp.CommonResp = &pbConversation.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}
			return resp, nil
		}
	}
	var conversation db.Conversation
	if err := utils.CopyStructFields(&conversation, req.Conversation); err != nil {
		log.NewDebug(req.OperationID, utils.GetSelfFuncName(), "CopyStructFields failed", *req.Conversation, err.Error())
	}
	haveUserID, _ := imdb.GetExistConversationUserIDList(req.UserIDList, req.Conversation.ConversationID)
	switch req.FieldType {
	case constant.FieldRecvMsgOpt:
		for _, v := range req.UserIDList {
			if err = db.DB.SetSingleConversationRecvMsgOpt(v, req.Conversation.ConversationID, req.Conversation.RecvMsgOpt); err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "cache failed, rpc return", err.Error())
				resp.CommonResp = &pbConversation.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
				return resp, nil
			}
			if req.Conversation.ConversationType == constant.SuperGroupChatType {
				if req.Conversation.RecvMsgOpt == constant.ReceiveNotNotifyMessage {
					if err = db.DB.SetSuperGroupUserReceiveNotNotifyMessage(req.Conversation.GroupID, v); err != nil {
						log.NewError(req.OperationID, utils.GetSelfFuncName(), "cache failed, rpc return", err.Error(), req.Conversation.GroupID, v)
						resp.CommonResp = &pbConversation.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
						return resp, nil
					}
				} else {
					if err = db.DB.SetSuperGroupUserReceiveNotifyMessage(req.Conversation.GroupID, v); err != nil {
						log.NewError(req.OperationID, utils.GetSelfFuncName(), "cache failed, rpc return", err.Error(), req.Conversation.GroupID, v)
						resp.CommonResp = &pbConversation.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
						return resp, nil
					}
				}
			}
		}
		err = imdb.UpdateColumnsConversations(haveUserID, req.Conversation.ConversationID, map[string]interface{}{"recv_msg_opt": conversation.RecvMsgOpt})
	case constant.FieldGroupAtType:
		err = imdb.UpdateColumnsConversations(haveUserID, req.Conversation.ConversationID, map[string]interface{}{"group_at_type": conversation.GroupAtType})
	case constant.FieldIsNotInGroup:
		err = imdb.UpdateColumnsConversations(haveUserID, req.Conversation.ConversationID, map[string]interface{}{"is_not_in_group": conversation.IsNotInGroup})
	case constant.FieldIsPinned:
		err = imdb.UpdateColumnsConversations(haveUserID, req.Conversation.ConversationID, map[string]interface{}{"is_pinned": conversation.IsPinned})
	case constant.FieldIsPrivateChat:
		err = imdb.UpdateColumnsConversations(haveUserID, req.Conversation.ConversationID, map[string]interface{}{"is_private_chat": conversation.IsPrivateChat})
	case constant.FieldEx:
		err = imdb.UpdateColumnsConversations(haveUserID, req.Conversation.ConversationID, map[string]interface{}{"ex": conversation.Ex})
	case constant.FieldAttachedInfo:
		err = imdb.UpdateColumnsConversations(haveUserID, req.Conversation.ConversationID, map[string]interface{}{"attached_info": conversation.AttachedInfo})
	case constant.FieldUnread:
		isSyncConversation = false
		err = imdb.UpdateColumnsConversations(haveUserID, req.Conversation.ConversationID, map[string]interface{}{"update_unread_count_time": conversation.UpdateUnreadCountTime})
	case constant.FieldBurnDuration:
		err = imdb.UpdateColumnsConversations(haveUserID, req.Conversation.ConversationID, map[string]interface{}{"burn_duration": conversation.BurnDuration})
	}
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "UpdateColumnsConversations error", err.Error())
		resp.CommonResp = &pbConversation.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
		return resp, nil
	}
	for _, v := range utils.DifferenceString(haveUserID, req.UserIDList) {
		conversation.OwnerUserID = v
		err = rocksCache.DelUserConversationIDListFromCache(v)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), v, req.Conversation.ConversationID, err.Error())
		}
		err := imdb.SetOneConversation(conversation)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "SetConversation error", err.Error())
			resp.CommonResp = &pbConversation.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
			return resp, nil
		}
	}

	// notification
	if req.Conversation.ConversationType == constant.SingleChatType && req.FieldType == constant.FieldIsPrivateChat {
		//sync peer user conversation if conversation is singleChatType
		if err := syncPeerUserConversation(req.Conversation, req.OperationID); err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "syncPeerUserConversation", err.Error())
			resp.CommonResp = &pbConversation.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
			return resp, nil
		}

	} else {
		if isSyncConversation {
			for _, v := range req.UserIDList {
				if err = rocksCache.DelConversationFromCache(v, req.Conversation.ConversationID); err != nil {
					log.NewError(req.OperationID, utils.GetSelfFuncName(), v, req.Conversation.ConversationID, err.Error())
				}
				chat.ConversationChangeNotification(req.OperationID, v)
			}
		} else {
			for _, v := range req.UserIDList {
				if err = rocksCache.DelConversationFromCache(v, req.Conversation.ConversationID); err != nil {
					log.NewError(req.OperationID, utils.GetSelfFuncName(), v, req.Conversation.ConversationID, err.Error())
				}
				chat.ConversationUnreadChangeNotification(req.OperationID, v, req.Conversation.ConversationID, conversation.UpdateUnreadCountTime)
			}
		}

	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "rpc return", resp.String())
	resp.CommonResp = &pbConversation.CommonResp{}
	return resp, nil
}
func syncPeerUserConversation(conversation *pbConversation.Conversation, operationID string) error {
	peerUserConversation := db.Conversation{
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
	err = rocksCache.DelUserConversationIDListFromCache(conversation.UserID)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "DelConversationFromCache failed", err.Error(), conversation.OwnerUserID, conversation.ConversationID)
	}
	err = rocksCache.DelConversationFromCache(conversation.UserID, utils.GetConversationIDBySessionType(conversation.OwnerUserID, constant.SingleChatType))
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "DelConversationFromCache failed", err.Error(), conversation.OwnerUserID, conversation.ConversationID)
	}
	err = rocksCache.DelConversationFromCache(conversation.OwnerUserID, conversation.ConversationID)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "DelConversationFromCache failed", err.Error(), conversation.OwnerUserID, conversation.ConversationID)
	}
	chat.ConversationSetPrivateNotification(operationID, conversation.OwnerUserID, conversation.UserID, conversation.IsPrivateChat)
	return nil
}
func NewRpcConversationServer(port int) *rpcConversation {
	log.NewPrivateLog(constant.LogFileName)
	return &rpcConversation{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImConversationName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (rpc *rpcConversation) Run() {
	log.NewInfo("0", "rpc conversation start...")

	listenIP := ""
	if config.Config.ListenIP == "" {
		listenIP = "0.0.0.0"
	} else {
		listenIP = config.Config.ListenIP
	}
	address := listenIP + ":" + strconv.Itoa(rpc.rpcPort)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic("listening err:" + err.Error() + rpc.rpcRegisterName)
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
	pbConversation.RegisterConversationServer(srv, rpc)
	rpcRegisterIP := config.Config.RpcRegisterIP
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error("", "GetLocalIP failed ", err.Error())
		}
	}
	log.NewInfo("", "rpcRegisterIP", rpcRegisterIP)
	err = getcdv3.RegisterEtcd(rpc.etcdSchema, strings.Join(rpc.etcdAddr, ","), rpcRegisterIP, rpc.rpcPort, rpc.rpcRegisterName, 10)
	if err != nil {
		log.NewError("0", "RegisterEtcd failed ", err.Error(),
			rpc.etcdSchema, strings.Join(rpc.etcdAddr, ","), rpcRegisterIP, rpc.rpcPort, rpc.rpcRegisterName)
		panic(utils.Wrap(err, "register conversation module  rpc to etcd err"))
	}
	log.NewInfo("0", "RegisterConversationServer ok ", rpc.etcdSchema, strings.Join(rpc.etcdAddr, ","), rpcRegisterIP, rpc.rpcPort, rpc.rpcRegisterName)
	err = srv.Serve(listener)
	if err != nil {
		log.NewError("0", "Serve failed ", err.Error())
		return
	}
	log.NewInfo("0", "rpc conversation ok")
}
