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
	// conn, err := discov.GetConn(context.Background(), config.Config.RpcRegisterName.OpenImAuthName)
	// if err != nil {
	// panic(err)
	// }
	return &Auth{discov: discov}
}

type Auth struct {
	conn   *grpc.ClientConn
	discov discoveryregistry.SvcDiscoveryRegistry
}

func (o *Auth) client(ctx context.Context) (auth.AuthClient, error) {
	c, err := o.discov.GetConn(ctx, config.Config.RpcRegisterName.OpenImAuthName)
	if err != nil {
		return nil, err
	}
	return auth.NewAuthClient(c), nil
}

func (o *Auth) UserRegister(c *gin.Context) {
	//a2r.Call(auth.AuthClient.UserRegister, o.userClient, c) // todo
}

func (o *Auth) UserToken(c *gin.Context) {
	a2r.Call(auth.AuthClient.UserToken, o.client, c)
}

func (o *Auth) ParseToken(c *gin.Context) {
	a2r.Call(auth.AuthClient.ParseToken, o.client, c)
}

func (o *Auth) ForceLogout(c *gin.Context) {
	a2r.Call(auth.AuthClient.ForceLogout, o.client, c)
}
