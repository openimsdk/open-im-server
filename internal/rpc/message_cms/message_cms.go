package messageCMS

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/http"
	"context"

	"Open_IM/pkg/common/log"

	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbMessageCMS "Open_IM/pkg/proto/message_cms"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"

	"Open_IM/pkg/utils"

	"net"
	"strconv"
	"strings"

	"google.golang.org/grpc"
)

type messageCMSServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func NewMessageCMSServer(port int) *messageCMSServer {
	log.NewPrivateLog("MessageCMS")
	return &messageCMSServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImMessageCMSName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (s *messageCMSServer) Run() {
	log.NewInfo("0", "messageCMS rpc start ")
	ip := utils.ServerIP
	registerAddress := ip + ":" + strconv.Itoa(s.rpcPort)
	//listener network
	listener, err := net.Listen("tcp", registerAddress)
	if err != nil {
		log.NewError("0", "Listen failed ", err.Error(), registerAddress)
		return
	}
	log.NewInfo("0", "listen network success, ", registerAddress, listener)
	defer listener.Close()
	//grpc server
	srv := grpc.NewServer()
	defer srv.GracefulStop()
	//Service registers with etcd
	pbMessageCMS.RegisterMessageCMSServer(srv, s)
	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), ip, s.rpcPort, s.rpcRegisterName, 10)
	if err != nil {
		log.NewError("0", "RegisterEtcd failed ", err.Error())
		return
	}
	err = srv.Serve(listener)
	if err != nil {
		log.NewError("0", "Serve failed ", err.Error())
		return
	}
	log.NewInfo("0", "message cms rpc success")
}

func (s *messageCMSServer) BoradcastMessage(_ context.Context, req *pbMessageCMS.BoradcastMessageReq) (*pbMessageCMS.BoradcastMessageResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "BoradcastMessage", req.String())
	resp := &pbMessageCMS.BoradcastMessageResp{}
	return resp, http.WarpError(constant.ErrDB)
}

func (s *messageCMSServer) GetChatLogs(_ context.Context, req *pbMessageCMS.GetChatLogsReq) (*pbMessageCMS.GetChatLogsResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "GetChatLogs", req.String())
	resp := &pbMessageCMS.GetChatLogsResp{}
	chatLog := db.ChatLog{
		Content: req.Content,
	}
	chatLogs, err := imdb.GetChatLog(chatLog, req.Pagination.ShowNumber, req.Pagination.PageNumber)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetChatLog", err.Error())
		return resp, http.WarpError(constant.ErrDB)
	}
	for _, chatLog := range chatLogs {
		pbChatLog := &pbMessageCMS.ChatLogs{
			SessionType:     chatLog.SessionType,
			ContentType:     chatLog.ContentType,
			SenderNickName:  chatLog.SenderNickname,
			SenderId:        chatLog.SendID,
			SearchContent:   req.Content,
			WholeContent:    chatLog.Content,
			Date:            chatLog.CreateTime.String(),
		}
		switch chatLog.SessionType {
		case constant.SingleChatType:
			recvUser, err := imdb.GetUserByUserID(chatLog.RecvID)
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetUserByUserID failed", err.Error())
				continue
			}
			pbChatLog.ReciverId = recvUser.UserID
			pbChatLog.ReciverNickName = recvUser.Nickname
		case constant.GroupChatType:
			group, err := imdb.GetGroupById(chatLog.RecvID)
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetGroupById failed")
				continue
			}
			pbChatLog.GroupId = group.GroupID
			pbChatLog.GroupName = group.GroupName
		}
		resp.ChatLogs = append(resp.ChatLogs, pbChatLog)
	}
	resp.Pagination = &open_im_sdk.ResponsePagination{
		CurrentPage: req.Pagination.PageNumber,
		ShowNumber:  req.Pagination.ShowNumber,
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp output: ", resp.String())
	return resp, nil
}

func (s *messageCMSServer) MassSendMessage(_ context.Context, req *pbMessageCMS.MassSendMessageReq) (*pbMessageCMS.MassSendMessageResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "MassSendMessage", req.String())
	resp := &pbMessageCMS.MassSendMessageResp{}
	return resp, nil
}

func (s *messageCMSServer) WithdrawMessage(_ context.Context, req *pbMessageCMS.WithdrawMessageReq) (*pbMessageCMS.WithdrawMessageResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "WithdrawMessage", req.String())
	resp := &pbMessageCMS.WithdrawMessageResp{}

	return resp, nil
}
