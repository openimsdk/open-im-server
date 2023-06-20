package api

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/a2r"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/user"
	"github.com/gin-gonic/gin"
)

func NewStatistics(discov discoveryregistry.SvcDiscoveryRegistry) *Statistics {
	return &Statistics{discov: discov}
}

type Statistics struct {
	discov discoveryregistry.SvcDiscoveryRegistry
}

func (s *Statistics) userClient(ctx context.Context) (user.UserClient, error) {
	conn, err := s.discov.GetConn(ctx, config.Config.RpcRegisterName.OpenImUserName)
	if err != nil {
		return nil, err
	}
	return user.NewUserClient(conn), nil
}

func (s *Statistics) UserRegister(c *gin.Context) {
	a2r.Call(user.UserClient.UserRegisterCount, s.userClient, c)
}
