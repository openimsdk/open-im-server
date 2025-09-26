package api

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/prommetrics"
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/discovery/etcd"
	"github.com/openimsdk/tools/errs"
)

type PrometheusDiscoveryApi struct {
	config *Config
	client discovery.SvcDiscoveryRegistry
}

func NewPrometheusDiscoveryApi(cfg *Config, client discovery.SvcDiscoveryRegistry) *PrometheusDiscoveryApi {
	api := &PrometheusDiscoveryApi{
		config: cfg,
	}
	if cfg.Discovery.Enable == config.ETCD {
		api.client = client.(*etcd.SvcDiscoveryRegistryImpl)
	}
	return api
}

func (p *PrometheusDiscoveryApi) Enable(c *gin.Context) {
	if p.config.Discovery.Enable != config.ETCD {
		c.JSON(http.StatusOK, []struct{}{})
		c.Abort()
	}
}

func (p *PrometheusDiscoveryApi) discovery(c *gin.Context, key string) {
	value, err := p.client.GetKeyWithPrefix(c, prommetrics.BuildDiscoveryKeyPrefix(key))
	if err != nil {
		apiresp.GinError(c, errs.WrapMsg(err, "get key value"))
		return
	}

	if len(value) == 0 {
		c.JSON(http.StatusOK, []*prommetrics.RespTarget{})
		return
	}
	var resp prommetrics.RespTarget
	for i := range value {
		var tmp prommetrics.Target
		if err = json.Unmarshal(value[i], &tmp); err != nil {
			apiresp.GinError(c, errs.WrapMsg(err, "json unmarshal err"))
			return
		}

		resp.Targets = append(resp.Targets, tmp.Target)
		resp.Labels = tmp.Labels // default label is fixed. See prommetrics.BuildDefaultTarget
	}

	c.JSON(http.StatusOK, []*prommetrics.RespTarget{&resp})
}

func (p *PrometheusDiscoveryApi) Api(c *gin.Context) {
	p.discovery(c, prommetrics.APIKeyName)
}

func (p *PrometheusDiscoveryApi) User(c *gin.Context) {
	p.discovery(c, p.config.Share.RpcRegisterName.User)
}

func (p *PrometheusDiscoveryApi) Group(c *gin.Context) {
	p.discovery(c, p.config.Share.RpcRegisterName.Group)
}

func (p *PrometheusDiscoveryApi) Msg(c *gin.Context) {
	p.discovery(c, p.config.Share.RpcRegisterName.Msg)
}

func (p *PrometheusDiscoveryApi) Friend(c *gin.Context) {
	p.discovery(c, p.config.Share.RpcRegisterName.Friend)
}

func (p *PrometheusDiscoveryApi) Conversation(c *gin.Context) {
	p.discovery(c, p.config.Share.RpcRegisterName.Conversation)
}

func (p *PrometheusDiscoveryApi) Third(c *gin.Context) {
	p.discovery(c, p.config.Share.RpcRegisterName.Third)
}

func (p *PrometheusDiscoveryApi) Auth(c *gin.Context) {
	p.discovery(c, p.config.Share.RpcRegisterName.Auth)
}

func (p *PrometheusDiscoveryApi) Push(c *gin.Context) {
	p.discovery(c, p.config.Share.RpcRegisterName.Push)
}

func (p *PrometheusDiscoveryApi) MessageGateway(c *gin.Context) {
	p.discovery(c, p.config.Share.RpcRegisterName.MessageGateway)
}

func (p *PrometheusDiscoveryApi) MessageTransfer(c *gin.Context) {
	p.discovery(c, prommetrics.MessageTransferKeyName)
}
