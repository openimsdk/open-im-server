package getcdv3

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/trace_log"
	"context"
	"github.com/OpenIMSDK/getcdv3"
	"google.golang.org/grpc"
	"strings"
	"sync"
)

func GetDefaultConn(arg1, arg2, arg3, arg4 string) *grpc.ClientConn {
	return nil
}

func GetConn(ctx context.Context, serviceName string) (conn *grpc.ClientConn, err error) {
	defer func() {
		trace_log.SetContextInfo(ctx, "GetConn", err, "serviceName", serviceName)
	}()
	conn = getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","),
		serviceName, trace_log.GetOperationID(ctx), config.Config.Etcd.UserName, config.Config.Etcd.Password)
	if conn == nil {
		return nil, constant.ErrRpcConn
	}
	return conn, nil
}

func GetDefaultGatewayConn4Unique(schema, addr, operationID string) []*grpc.ClientConn {
	return nil
}

func RegisterEtcd(schema, etcdAddr, myHost string, myPort int, serviceName string, ttl int, operationID string) error {
	return getcdv3.RegisterEtcd(schema, etcdAddr, myHost, myPort, serviceName, ttl, operationID)
}

var Conn4UniqueList []*grpc.ClientConn
var Conn4UniqueListMtx sync.RWMutex
var IsUpdateStart bool
var IsUpdateStartMtx sync.RWMutex
