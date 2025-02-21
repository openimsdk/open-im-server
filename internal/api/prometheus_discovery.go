package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	conf "github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/prommetrics"
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/discovery/etcd"
	"github.com/openimsdk/tools/errs"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type PrometheusDiscoveryApi struct {
	config *Config
	client *clientv3.Client
	kv     discovery.KeyValue
}

func NewPrometheusDiscoveryApi(config *Config, client discovery.Conn) *PrometheusDiscoveryApi {
	api := &PrometheusDiscoveryApi{
		config: config,
	}
	if config.Discovery.Enable == conf.ETCD {
		api.client = client.(*etcd.SvcDiscoveryRegistryImpl).GetClient()
	}
	return api
}

func (p *PrometheusDiscoveryApi) discovery(c *gin.Context, key string) {
	value, err := p.kv.GetKey(c, prommetrics.BuildDiscoveryKey(key))
	if err != nil {
		if errors.Is(err, discovery.ErrNotSupportedKeyValue) {
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
	if err := json.Unmarshal(value, &resp); err != nil {
		apiresp.GinError(c, errs.WrapMsg(err, "json unmarshal err"))
		return
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
