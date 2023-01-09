package utils

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/trace_log"
	"context"
	"github.com/OpenIMSDK/getcdv3"

	"google.golang.org/grpc"
	"strings"
)

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
