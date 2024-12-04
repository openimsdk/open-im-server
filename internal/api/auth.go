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
	"github.com/gin-gonic/gin"

	"github.com/openimsdk/protocol/auth"
	"github.com/openimsdk/tools/a2r"

	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
)

type AuthApi rpcclient.Auth

func NewAuthApi(client rpcclient.Auth) AuthApi {
	return AuthApi(client)
}

func (o *AuthApi) GetAdminToken(c *gin.Context) {
	a2r.Call(auth.AuthClient.GetAdminToken, o.Client, c)
}

func (o *AuthApi) GetUserToken(c *gin.Context) {
	a2r.Call(auth.AuthClient.GetUserToken, o.Client, c)
}

func (o *AuthApi) ParseToken(c *gin.Context) {
	a2r.Call(auth.AuthClient.ParseToken, o.Client, c)
}

func (o *AuthApi) ForceLogout(c *gin.Context) {
	a2r.Call(auth.AuthClient.ForceLogout, o.Client, c)
}
