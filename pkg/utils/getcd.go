package utils

import (
	"Open_IM/pkg/common/config"
	"github.com/OpenIMSDK/getcdv3"
	"github.com/OpenIMSDK/open_utils/constant"
	"google.golang.org/grpc"
	"strings"
)

func GetConn(operationID, serviceName string) (*grpc.ClientConn, error) {
	conn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","),
		serviceName, operationID, config.Config.Etcd.UserName, config.Config.Etcd.Password)
	if conn == nil {
		return nil, constant.ErrGetRpcConn
	}
	return conn, nil
}
