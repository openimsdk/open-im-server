package conversation

import (
	chat "Open_IM/internal/rpc/msg"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbConversation "Open_IM/pkg/proto/conversation"
	"Open_IM/pkg/utils"
	"context"
	"net"
	"strconv"
	"strings"

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
	if req.Conversation.ConversationType == constant.GroupChatType {
		groupInfo, err := imdb.GetGroupInfoByGroupID(req.Conversation.GroupID)
		if err != nil {
			log.NewError(req.OperationID, "GetGroupInfoByGroupID failed ", req.Conversation.GroupID, err.Error())
			resp.CommonResp = &pbConversation.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
			return resp, nil
		}
		if groupInfo.Status == constant.GroupStatusDismissed && !req.Conversation.IsNotInGroup {
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
		err = imdb.UpdateColumnsConversations(haveUserID, req.Conversation.ConversationID, map[string]interface{}{"unread_count": conversation.UnreadCount})

	}
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "UpdateColumnsConversations error", err.Error())
		resp.CommonResp = &pbConversation.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}
		return resp, nil
	}
	for _, v := range utils.DifferenceString(haveUserID, req.UserIDList) {
		conversation.OwnerUserID = v
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
		for _, v := range req.UserIDList {
			chat.ConversationChangeNotification(req.OperationID, v)
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
	srv := grpc.NewServer()
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
		return
	}
	log.NewInfo("0", "RegisterConversationServer ok ", rpc.etcdSchema, strings.Join(rpc.etcdAddr, ","), rpcRegisterIP, rpc.rpcPort, rpc.rpcRegisterName)
	err = srv.Serve(listener)
	if err != nil {
		log.NewError("0", "Serve failed ", err.Error())
		return
	}
	log.NewInfo("0", "rpc conversation ok")
}
