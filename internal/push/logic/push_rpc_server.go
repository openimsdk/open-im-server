package logic

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	"Open_IM/pkg/proto/push"
	"Open_IM/pkg/utils"
	"context"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"strings"
)

type RPCServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func (r *RPCServer) Init(rpcPort int) {
	r.rpcPort = rpcPort
	r.rpcRegisterName = config.Config.RpcRegisterName.OpenImPushName
	r.etcdSchema = config.Config.Etcd.EtcdSchema
	r.etcdAddr = config.Config.Etcd.EtcdAddr
}
func (r *RPCServer) run() {
	listenIP := ""
	if config.Config.ListenIP == "" {
		listenIP = "0.0.0.0"
	} else {
		listenIP = config.Config.ListenIP
	}
	address := listenIP + ":" + strconv.Itoa(r.rpcPort)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic("listening err:" + err.Error() + r.rpcRegisterName)
	}
	defer listener.Close()
	srv := grpc.NewServer()
	defer srv.GracefulStop()
	pbPush.RegisterPushMsgServiceServer(srv, r)
	rpcRegisterIP := config.Config.RpcRegisterIP
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error("", "GetLocalIP failed ", err.Error())
		}
	}

	err = getcdv3.RegisterEtcd(r.etcdSchema, strings.Join(r.etcdAddr, ","), rpcRegisterIP, r.rpcPort, r.rpcRegisterName, 10)
	if err != nil {
		log.Error("", "register push module  rpc to etcd err", err.Error(), r.etcdSchema, strings.Join(r.etcdAddr, ","), rpcRegisterIP, r.rpcPort, r.rpcRegisterName)
	}
	err = srv.Serve(listener)
	if err != nil {
		log.Error("", "push module rpc start err", err.Error())
		return
	}
}
func (r *RPCServer) PushMsg(_ context.Context, pbData *pbPush.PushMsgReq) (*pbPush.PushMsgResp, error) {
	//Call push module to send message to the user
	switch pbData.MsgData.SessionType {
	case constant.SuperGroupChatType:
		MsgToSuperGroupUser(pbData)
	default:
		MsgToUser(pbData)
	}
	return &pbPush.PushMsgResp{
		ResultCode: 0,
	}, nil

}
