package MessageCMS

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db"
	"context"

	//"Open_IM/pkg/common/constant"
	//"Open_IM/pkg/common/db"

	"Open_IM/pkg/common/log"

	//cp "Open_IM/pkg/common/utils"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbMessageCMS "Open_IM/pkg/proto/message_cms"

	"Open_IM/pkg/utils"
	//"context"
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
	log.NewPrivateLog("Statistics")
	return &messageCMSServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImStatisticsName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (s *messageCMSServer) Run() {
	log.NewInfo("0", "Statistics rpc start ")
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
	pbMessageCMS.RegisterMessageServer(srv, s)
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
	log.NewInfo("0", "statistics rpc success")
}

func (s *messageCMSServer) BoradcastMessage(_ context.Context, req *pbMessageCMS.BoradcastMessageReq) (*pbMessageCMS.BoradcastMessageResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "BoradcastMessage", req.String())
	resp := &pbMessageCMS.BoradcastMessageResp{}
	return resp, nil
}

func (s *messageCMSServer) GetChatLogs(_ context.Context, req *pbMessageCMS.GetChatLogsReq) (*pbMessageCMS.GetChatLogsResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "GetChatLogs", req.String())
	resp := &pbMessageCMS.GetChatLogsResp{}
	chatLog := db.ChatLog{
		Content: req.Content,
	}
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
