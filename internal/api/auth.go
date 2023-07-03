package api

import (
	"github.com/gin-gonic/gin"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/a2r"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/auth"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
)

type AuthApi rpcclient.Auth

func NewAuthApi(discov discoveryregistry.SvcDiscoveryRegistry) AuthApi {
	return AuthApi(*rpcclient.NewAuth(discov))
}

func (o *AuthApi) UserToken(c *gin.Context) {
	a2r.Call(auth.AuthClient.UserToken, o.Client, c)
}

func (o *AuthApi) ParseToken(c *gin.Context) {
	a2r.Call(auth.AuthClient.ParseToken, o.Client, c)
}

func (o *AuthApi) ForceLogout(c *gin.Context) {
	a2r.Call(auth.AuthClient.ForceLogout, o.Client, c)
}
