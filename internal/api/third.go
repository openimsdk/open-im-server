package api

import (
	"OpenIM/internal/api/a2r"
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/discoveryregistry"
	"OpenIM/pkg/proto/third"
	"context"

	"github.com/gin-gonic/gin"
)

var _ context.Context // 解决goland编辑器bug

func NewThird(c discoveryregistry.SvcDiscoveryRegistry) *Third {
	return &Third{c: c}
}

type Third struct {
	c discoveryregistry.SvcDiscoveryRegistry
}

func (o *Third) client() (third.ThirdClient, error) {
	conn, err := o.c.GetConn(config.Config.RpcRegisterName.OpenImThirdName)
	if err != nil {
		return nil, err
	}
	return third.NewThirdClient(conn), nil
}

func (o *Third) ApplyPut(c *gin.Context) {
	a2r.Call(third.ThirdClient.ApplyPut, o.client, c)
}

func (o *Third) GetPut(c *gin.Context) {
	a2r.Call(third.ThirdClient.GetPut, o.client, c)
}

func (o *Third) ConfirmPut(c *gin.Context) {
	a2r.Call(third.ThirdClient.ConfirmPut, o.client, c)
}

func (o *Third) GetSignalInvitationInfo(c *gin.Context) {
	a2r.Call(third.ThirdClient.GetSignalInvitationInfo, o.client, c)
}

func (o *Third) GetSignalInvitationInfoStartApp(c *gin.Context) {
	a2r.Call(third.ThirdClient.GetSignalInvitationInfoStartApp, o.client, c)
}

func (o *Third) FcmUpdateToken(c *gin.Context) {
	a2r.Call(third.ThirdClient.FcmUpdateToken, o.client, c)
}

func (o *Third) SetAppBadge(c *gin.Context) {
	a2r.Call(third.ThirdClient.SetAppBadge, o.client, c)
}
