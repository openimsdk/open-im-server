package api

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/internal/api/a2r"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	auth "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/auth"
	"github.com/gin-gonic/gin"
)

var _ context.Context // 解决goland编辑器bug

func NewAuth(c discoveryregistry.SvcDiscoveryRegistry) *Auth {
	return &Auth{c: c}
}

type Auth struct {
	c discoveryregistry.SvcDiscoveryRegistry
}

func (o *Auth) client() (auth.AuthClient, error) {
	conn, err := o.c.GetConn(config.Config.RpcRegisterName.OpenImAuthName)
	if err != nil {
		return nil, err
	}
	return auth.NewAuthClient(conn), nil
}

func (o *Auth) UserRegister(c *gin.Context) {
	//a2r.Call(auth.AuthClient.UserRegister, o.client, c) // todo
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
