package internal_service

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbUser "Open_IM/pkg/proto/user"
	"context"
	"strings"
)

func GetUserInfoClient(req *pbUser.GetUserInfoReq) (*pbUser.GetUserInfoResp, error) {
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImUserName)
	client := pbUser.NewUserClient(etcdConn)
	RpcResp, err := client.GetUserInfo(context.Background(), req)
	return RpcResp, err
}
