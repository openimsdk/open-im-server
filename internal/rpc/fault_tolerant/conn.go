package fault_tolerant

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/middleware"
	"Open_IM/pkg/utils"
	"github.com/OpenIMSDK/getcdv3"
	"google.golang.org/grpc"
	"strings"
)

func GetConfigConn(serviceName string, operationID string) *grpc.ClientConn {
	rpcRegisterIP := config.Config.RpcRegisterIP
	var err error
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error(operationID, "GetLocalIP failed ", err.Error())
			return nil
		}
	}

	var configPortList []int
	//1
	if config.Config.RpcRegisterName.OpenImUserName == serviceName {
		configPortList = config.Config.RpcPort.OpenImUserPort
	}
	//2
	if config.Config.RpcRegisterName.OpenImFriendName == serviceName {
		configPortList = config.Config.RpcPort.OpenImFriendPort
	}
	//3
	if config.Config.RpcRegisterName.OpenImMsgName == serviceName {
		configPortList = config.Config.RpcPort.OpenImMessagePort
	}
	//4
	if config.Config.RpcRegisterName.OpenImPushName == serviceName {
		configPortList = config.Config.RpcPort.OpenImPushPort
	}
	//5
	if config.Config.RpcRegisterName.OpenImRelayName == serviceName {
		configPortList = config.Config.RpcPort.OpenImMessageGatewayPort
	}
	//6
	if config.Config.RpcRegisterName.OpenImGroupName == serviceName {
		configPortList = config.Config.RpcPort.OpenImGroupPort
	}
	//7
	if config.Config.RpcRegisterName.OpenImAuthName == serviceName {
		configPortList = config.Config.RpcPort.OpenImAuthPort
	}
	//10
	if config.Config.RpcRegisterName.OpenImOfficeName == serviceName {
		configPortList = config.Config.RpcPort.OpenImOfficePort
	}
	//11
	if config.Config.RpcRegisterName.OpenImOrganizationName == serviceName {
		configPortList = config.Config.RpcPort.OpenImOrganizationPort
	}
	//12
	if config.Config.RpcRegisterName.OpenImConversationName == serviceName {
		configPortList = config.Config.RpcPort.OpenImConversationPort
	}
	//13
	if config.Config.RpcRegisterName.OpenImCacheName == serviceName {
		configPortList = config.Config.RpcPort.OpenImCachePort
	}
	//14
	if config.Config.RpcRegisterName.OpenImRealTimeCommName == serviceName {
		configPortList = config.Config.RpcPort.OpenImRealTimeCommPort
	}
	if len(configPortList) == 0 {
		log.Error(operationID, "len(configPortList) == 0  ")
		return nil
	}
	target := rpcRegisterIP + ":" + utils.Int32ToString(int32(configPortList[0]))
	log.Info(operationID, "rpcRegisterIP ", rpcRegisterIP, " port ", configPortList, " grpc target: ", target, " serviceName: ", serviceName)
	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithUnaryInterceptor(middleware.RpcClientInterceptor))
	if err != nil {
		log.Error(operationID, "grpc.Dail failed ", err.Error())
		return nil
	}
	log.NewDebug(operationID, utils.GetSelfFuncName(), serviceName, conn)
	return conn
}

func GetDefaultConn(serviceName string, operationID string) (*grpc.ClientConn, error) {
	con := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), serviceName, operationID, config.Config.Etcd.UserName, config.Config.Etcd.Password)
	if con != nil {
		return con, nil
	}
	log.NewWarn(operationID, utils.GetSelfFuncName(), "conn is nil !!!!!", serviceName)
	con = GetConfigConn(serviceName, operationID)
	return con, nil
}
