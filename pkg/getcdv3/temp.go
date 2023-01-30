package getcdv3

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/tracelog"
	"context"
	"github.com/OpenIMSDK/getcdv3"
	"google.golang.org/grpc"
	"strings"
)

func GetDefaultConn(arg1, arg2, arg3, arg4 string) *grpc.ClientConn {
	return getcdv3.GetConn(arg1, arg2, arg3, arg4, config.Config.Etcd.UserName, config.Config.Etcd.Password)
}

func GetConn(ctx context.Context, serviceName string) (conn *grpc.ClientConn, err error) {
	defer func() {
		tracelog.SetCtxInfo(ctx, "GetConn", err, "serviceName", serviceName)
	}()
	conn = getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","),
		serviceName, tracelog.GetOperationID(ctx), config.Config.Etcd.UserName, config.Config.Etcd.Password)
	if conn == nil {
		return nil, constant.ErrInternalServer
	}
	return conn, nil
}

func GetDefaultGatewayConn4Unique(schema, addr, operationID string) []*grpc.ClientConn {
	return nil
}

func RegisterEtcd(schema, etcdAddr, myHost string, myPort int, serviceName string, ttl int, operationID string) error {
	return getcdv3.RegisterEtcd(schema, etcdAddr, myHost, myPort, serviceName, ttl, operationID)
}
