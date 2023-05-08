package api

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/apiresp"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/apistruct"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/tokenverify"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/a2r"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/user"
	"github.com/gin-gonic/gin"
)

var _ context.Context // 解决goland编辑器bug

func NewUser(client discoveryregistry.SvcDiscoveryRegistry) *User {
	return &User{c: client}
}

type User struct {
	c discoveryregistry.SvcDiscoveryRegistry
}

func (u *User) client(ctx context.Context) (user.UserClient, error) {
	conn, err := u.c.GetConn(ctx, config.Config.RpcRegisterName.OpenImUserName)
	if err != nil {
		return nil, err
	}
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
func (u *User) GetUsersOnlineStatus(c *gin.Context) {
	params := apistruct.ManagementSendMsgReq{}
	if err := c.BindJSON(&params); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WithDetail(err.Error()).Wrap())
		return
	}
	if !tokenverify.IsAppManagerUid(c) {
		apiresp.GinError(c, errs.ErrNoPermission.Wrap("only app manager can send message"))
		return
	}

}
