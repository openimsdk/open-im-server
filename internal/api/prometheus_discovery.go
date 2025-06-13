package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openimsdk/open-im-server/v3/pkg/common/prommetrics"
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/errs"
)

type PrometheusDiscoveryApi struct {
	config *Config
	kv     discovery.KeyValue
}

func NewPrometheusDiscoveryApi(config *Config, client discovery.SvcDiscoveryRegistry) *PrometheusDiscoveryApi {
	api := &PrometheusDiscoveryApi{
		config: config,
		kv:     client,
	}
	return api
}

func (p *PrometheusDiscoveryApi) discovery(c *gin.Context, key string) {
	value, err := p.kv.GetKeyWithPrefix(c, prommetrics.BuildDiscoveryKeyPrefix(key))
	if err != nil {
		if errors.Is(err, discovery.ErrNotSupported) {
			c.JSON(http.StatusOK, []struct{}{})
			return
		}
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
	p.discovery(c, p.config.Discovery.RpcService.User)
}

func (p *PrometheusDiscoveryApi) Group(c *gin.Context) {
	p.discovery(c, p.config.Discovery.RpcService.Group)
}

func (p *PrometheusDiscoveryApi) Msg(c *gin.Context) {
	p.discovery(c, p.config.Discovery.RpcService.Msg)
}

func (p *PrometheusDiscoveryApi) Friend(c *gin.Context) {
	p.discovery(c, p.config.Discovery.RpcService.Friend)
}

func (p *PrometheusDiscoveryApi) Conversation(c *gin.Context) {
	p.discovery(c, p.config.Discovery.RpcService.Conversation)
}

func (p *PrometheusDiscoveryApi) Third(c *gin.Context) {
	p.discovery(c, p.config.Discovery.RpcService.Third)
}

func (p *PrometheusDiscoveryApi) Auth(c *gin.Context) {
	p.discovery(c, p.config.Discovery.RpcService.Auth)
}

func (p *PrometheusDiscoveryApi) Push(c *gin.Context) {
	p.discovery(c, p.config.Discovery.RpcService.Push)
}

func (p *PrometheusDiscoveryApi) MessageGateway(c *gin.Context) {
	p.discovery(c, p.config.Discovery.RpcService.MessageGateway)
}

func (p *PrometheusDiscoveryApi) MessageTransfer(c *gin.Context) {
	p.discovery(c, prommetrics.MessageTransferKeyName)
}
