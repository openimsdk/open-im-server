package api

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/internal/api/a2r"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/third"
	"github.com/gin-gonic/gin"
	"net/http"
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

func (o *Third) GetURL(c *gin.Context) {
	if c.Request.Method == http.MethodPost {
		a2r.Call(third.ThirdClient.GetUrl, o.client, c)
		return
	}
	name := c.Query("name")
	if name == "" {
		c.String(http.StatusBadRequest, "name is empty")
		return
	}
	client, err := o.client()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	resp, err := client.GetUrl(c, &third.GetUrlReq{Name: name})
	if err != nil {
		if errs.ErrArgs.Is(err) {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		if errs.ErrRecordNotFound.Is(err) {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, resp.Url)
}
