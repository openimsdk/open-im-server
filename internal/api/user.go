package api

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/apiresp"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/apistruct"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/tokenverify"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"google.golang.org/grpc"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/a2r"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/user"
	"github.com/gin-gonic/gin"
)

func NewUser(discov discoveryregistry.SvcDiscoveryRegistry) *User {
	conn, err := discov.GetConn(context.Background(), config.Config.RpcRegisterName.OpenImUserName)
	if err != nil {
		// panic(err)
	}
	log.ZInfo(context.Background(), "user rpc conn", "conn", conn)
	return &User{conn: conn, discov: discov}
}

type User struct {
	discov discoveryregistry.SvcDiscoveryRegistry
	conn   *grpc.ClientConn
}

func (u *User) client(ctx context.Context) (user.UserClient, error) {
	conn, err := u.discov.GetConn(ctx, config.Config.RpcRegisterName.OpenImUserName)
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
