package user

import (
	"OpenIM/internal/a2r"
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

func (u *User) getGroupClient() (user.UserClient, error) {
	conn, err := u.zk.GetConn(config.Config.RpcRegisterName.OpenImGroupName)
	if err != nil {
		return nil, err
	}
	return user.NewUserClient(conn), nil
}

func (u *User) UpdateUserInfo(c *gin.Context) {
	a2r.Call(user.UserClient.UpdateUserInfo, u.getGroupClient, c)
}

func (u *User) SetGlobalRecvMessageOpt(c *gin.Context) {
	a2r.Call(user.UserClient.SetGlobalRecvMessageOpt, u.getGroupClient, c)
}

func (u *User) GetUsersPublicInfo(c *gin.Context) {
	a2r.Call(user.UserClient.GetDesignateUsers, u.getGroupClient, c)
}

func (u *User) GetSelfUserInfo(c *gin.Context) {
	a2r.Call(user.UserClient.GetSelfUserInfo, u.getGroupClient, c)
}

func (u *User) GetUsersOnlineStatus(c *gin.Context) {
	a2r.Call(user.UserClient.GetUsersOnlineStatus, u.getGroupClient, c)
}

func (u *User) GetUsersInfoFromCache(c *gin.Context) {
	a2r.Call(user.UserClient.GetUsersInfoFromCache, u.getGroupClient, c)
}

func (u *User) GetFriendIDListFromCache(c *gin.Context) {
	a2r.Call(user.UserClient.GetFriendIDListFromCache, u.getGroupClient, c)
}

func (u *User) GetBlackIDListFromCache(c *gin.Context) {
	a2r.Call(user.UserClient.GetBlackIDListFromCache, u.getGroupClient, c)
}

//func (u *User) GetAllUsersUid(c *gin.Context) {
//	a2r.Call(user.UserClient.GetAllUsersUid, u.getGroupClient, c)
//}
//
//func (u *User) AccountCheck(c *gin.Context) {
//	a2r.Call(user.UserClient.AccountCheck, u.getGroupClient, c)
//}

func (u *User) GetUsers(c *gin.Context) {
	a2r.Call(user.UserClient.GetPaginationUsers, u.getGroupClient, c)
}
