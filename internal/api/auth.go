package api

import (
	"OpenIM/internal/api/a2r"
	"OpenIM/pkg/common/config"
	auth "OpenIM/pkg/proto/auth"
	"context"
	"github.com/OpenIMSDK/openKeeper"
	"github.com/gin-gonic/gin"
)

var _ context.Context // 解决goland编辑器bug

func NewAuth(zk *openKeeper.ZkClient) *Auth {
	return &Auth{zk: zk}
}

type Auth struct {
	zk *openKeeper.ZkClient
}

func (o *Auth) client() (auth.AuthClient, error) {
	conn, err := o.zk.GetDefaultConn(config.Config.RpcRegisterName.OpenImAuthName)
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
