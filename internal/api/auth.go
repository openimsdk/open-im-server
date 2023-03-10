package api

import (
	"OpenIM/internal/api/a2r"
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/discoveryregistry"
	auth "OpenIM/pkg/proto/auth"
	"context"
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

	a2r.Call2(auth.AuthClient.UserToken, o.client, c)
}

func (o *Auth) ParseToken(c *gin.Context) {
	a2r.Call(auth.AuthClient.ParseToken, o.client, c)
}

func (o *Auth) ForceLogout(c *gin.Context) {
	a2r.Call(auth.AuthClient.ForceLogout, o.client, c)
}
