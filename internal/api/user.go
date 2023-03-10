package api

import (
	"OpenIM/internal/api/a2r"
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/log"
	"OpenIM/pkg/discoveryregistry"
	"OpenIM/pkg/proto/user"
	"context"
	"github.com/gin-gonic/gin"
)

var _ context.Context // 解决goland编辑器bug

func NewUser(client discoveryregistry.SvcDiscoveryRegistry) *User {
	return &User{c: client}
}

type User struct {
	c discoveryregistry.SvcDiscoveryRegistry
}

func (u *User) client() (user.UserClient, error) {
	conn, err := u.c.GetConn(config.Config.RpcRegisterName.OpenImUserName)
	if err != nil {
		return nil, err
	}
	defer func() {
		log.NewInfo("client", conn, err)
		conns, err := u.c.GetConns(config.Config.RpcRegisterName.OpenImUserName)
		log.NewInfo("conns", conns, err)
	}()
	return user.NewUserClient(conn), nil
}

func (u *User) UserRegister(c *gin.Context) {
	a2r.Call(user.UserClient.UserRegister, u.client, c)
}

func (u *User) UpdateUserInfo(c *gin.Context) {
	a2r.Call(user.UserClient.UpdateUserInfo, u.client, c)
}

func (u *User) SetGlobalRecvMessageOpt(c *gin.Context) {
	a2r.Call(user.UserClient.SetGlobalRecvMessageOpt, u.client, c)
}

func (u *User) GetUsersPublicInfo(c *gin.Context) {
	a2r.Call(user.UserClient.GetDesignateUsers, u.client, c)
}

func (u *User) GetAllUsersID(c *gin.Context) {
	a2r.Call(user.UserClient.GetDesignateUsers, u.client, c)
}

func (u *User) AccountCheck(c *gin.Context) {
	a2r.Call(user.UserClient.AccountCheck, u.client, c)
}

func (u *User) GetUsers(c *gin.Context) {
	a2r.Call(user.UserClient.GetPaginationUsers, u.client, c)
}
