package api

import (
	"github.com/gin-gonic/gin"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/a2r"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/user"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
)

type StatisticsApi rpcclient.User

func NewStatisticsApi(discov discoveryregistry.SvcDiscoveryRegistry) StatisticsApi {
	return StatisticsApi(*rpcclient.NewUser(discov))
}

func (s *StatisticsApi) UserRegister(c *gin.Context) {
	a2r.Call(user.UserClient.UserRegisterCount, s.Client, c)
}
