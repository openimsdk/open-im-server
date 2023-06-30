// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/a2r"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/apiresp"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/apistruct"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/tokenverify"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/user"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
	"github.com/gin-gonic/gin"
)

type UserApi rpcclient.User

func NewUserApi(discov discoveryregistry.SvcDiscoveryRegistry) UserApi {
	return UserApi(*rpcclient.NewUser(discov))
}

func (u *UserApi) UserRegister(c *gin.Context) {
	a2r.Call(user.UserClient.UserRegister, u.Client, c)
}

func (u *UserApi) UpdateUserInfo(c *gin.Context) {
	a2r.Call(user.UserClient.UpdateUserInfo, u.Client, c)
}

func (u *UserApi) SetGlobalRecvMessageOpt(c *gin.Context) {
	a2r.Call(user.UserClient.SetGlobalRecvMessageOpt, u.Client, c)
}

func (u *UserApi) GetUsersPublicInfo(c *gin.Context) {
	a2r.Call(user.UserClient.GetDesignateUsers, u.Client, c)
}

func (u *UserApi) GetAllUsersID(c *gin.Context) {
	a2r.Call(user.UserClient.GetDesignateUsers, u.Client, c)
}

func (u *UserApi) AccountCheck(c *gin.Context) {
	a2r.Call(user.UserClient.AccountCheck, u.Client, c)
}

func (u *UserApi) GetUsers(c *gin.Context) {
	a2r.Call(user.UserClient.GetPaginationUsers, u.Client, c)
}

func (u *UserApi) GetUsersOnlineStatus(c *gin.Context) {
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
