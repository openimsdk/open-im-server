package api

import (
	"OpenIM/internal/api/a2r"
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/proto/user"
	"context"
	"github.com/OpenIMSDK/openKeeper"
	"github.com/gin-gonic/gin"
)

var _ context.Context // 解决goland编辑器bug

func NewUser(zk *openKeeper.ZkClient) *User {
	return &User{zk: zk}
}

type User struct {
	zk *openKeeper.ZkClient
}

func (o *User) client() (user.UserClient, error) {
	conn, err := o.zk.GetConn(config.Config.RpcRegisterName.OpenImGroupName)
	if err != nil {
		return nil, err
	}
	return user.NewUserClient(conn), nil
}

func (o *User) UserRegister(c *gin.Context) {
	a2r.Call(user.UserClient.UserRegister, o.client, c)
}

func (o *User) UpdateUserInfo(c *gin.Context) {
	a2r.Call(user.UserClient.UpdateUserInfo, o.client, c)
}

func (o *User) SetGlobalRecvMessageOpt(c *gin.Context) {
	a2r.Call(user.UserClient.SetGlobalRecvMessageOpt, o.client, c)
}

func (o *User) GetUsersPublicInfo(c *gin.Context) {
	a2r.Call(user.UserClient.GetDesignateUsers, o.client, c)
}

func (o *User) GetAllUsersID(c *gin.Context) {
	a2r.Call(user.UserClient.GetDesignateUsers, o.client, c)
}

//
func (u *User) AccountCheck(c *gin.Context) {
	a2r.Call(user.UserClient.AccountCheck, u.client, c)
}

func (o *User) GetUsers(c *gin.Context) {
	a2r.Call(user.UserClient.GetPaginationUsers, o.client, c)
}
