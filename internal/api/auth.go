package api

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/a2r"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/auth"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func NewAuth(discov discoveryregistry.SvcDiscoveryRegistry) *Auth {
	conn, err := discov.GetConn(context.Background(), config.Config.RpcRegisterName.OpenImAuthName)
	if err != nil {
		panic(err)
	}
	client := auth.NewAuthClient(conn)
	return &Auth{discov: discov, conn: conn, client: client}
}

type Auth struct {
	conn   *grpc.ClientConn
	client auth.AuthClient
	discov discoveryregistry.SvcDiscoveryRegistry
}

func (o *Auth) Client() auth.AuthClient {
	return o.client
}

func (o *Auth) UserToken(c *gin.Context) {
	a2r.Call(auth.AuthClient.UserToken, o.Client, c)
}

func (o *Auth) ParseToken(c *gin.Context) {
	a2r.Call(auth.AuthClient.ParseToken, o.Client, c)
}

func (o *Auth) ForceLogout(c *gin.Context) {
	a2r.Call(auth.AuthClient.ForceLogout, o.Client, c)
}
