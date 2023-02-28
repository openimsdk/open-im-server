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

func (a *Auth) getGroupClient() (auth.AuthClient, error) {
	conn, err := a.zk.GetConn(config.Config.RpcRegisterName.OpenImGroupName)
	if err != nil {
		return nil, err
	}
	return auth.NewAuthClient(conn), nil
}

func (a *Auth) UserRegister(c *gin.Context) {
	a2r.Call(auth.AuthClient.UserRegister, a.getGroupClient, c)
}

func (a *Auth) UserToken(c *gin.Context) {
	a2r.Call(auth.AuthClient.UserToken, a.getGroupClient, c)
}

func (a *Auth) ParseToken(c *gin.Context) {
	a2r.Call(auth.AuthClient.ParseToken, a.getGroupClient, c)
}

func (a *Auth) ForceLogout(c *gin.Context) {
	a2r.Call(auth.AuthClient.ForceLogout, a.getGroupClient, c)
}
