package internal_service

import (
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbUser "Open_IM/pkg/proto/user"
	"context"
)

func GetUserInfoClient(req *pbUser.GetUserInfoReq) (*pbUser.GetUserInfoResp, error) {
	etcdConn := getcdv3.GetUserConn()
	client := pbUser.NewUserClient(etcdConn)
	RpcResp, err := client.GetUserInfo(context.Background(), req)
	return RpcResp, err
}
