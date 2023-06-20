package api

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/a2r"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/user"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func NewStatistics(discov discoveryregistry.SvcDiscoveryRegistry) *Statistics {
	conn, err := discov.GetConn(context.Background(), config.Config.RpcRegisterName.OpenImUserName)
	if err != nil {
		panic(err)
	}
	client := user.NewUserClient(conn)
	return &Statistics{discov: discov, client: client, conn: conn}
}

type Statistics struct {
	conn   *grpc.ClientConn
	client user.UserClient
	discov discoveryregistry.SvcDiscoveryRegistry
}

func (s *Statistics) Client() user.UserClient {
	return s.client
}

func (s *Statistics) UserRegister(c *gin.Context) {
	a2r.Call(user.UserClient.UserRegisterCount, s.Client, c)
}
