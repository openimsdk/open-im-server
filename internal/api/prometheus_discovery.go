package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/open-im-server/v3/pkg/common/discoveryregister"
	"github.com/openimsdk/open-im-server/v3/pkg/common/prommetrics"
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/discovery/etcd"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	clientv3 "go.etcd.io/etcd/client/v3"
	"net/http"
)

type PrometheusDiscoveryApi struct {
	config *Config
	client *clientv3.Client
}

func NewPrometheusDiscoveryApi(config *Config, client discovery.SvcDiscoveryRegistry) *PrometheusDiscoveryApi {
	api := &PrometheusDiscoveryApi{
		config: config,
	}
	if config.Discovery.Enable == discoveryregister.Etcd {
		api.client = client.(*etcd.SvcDiscoveryRegistryImpl).GetClient()
	}
	return api
}

func (p *PrometheusDiscoveryApi) Enable(c *gin.Context) {
	if p.config.Discovery.Enable != discoveryregister.Etcd {
		c.JSON(http.StatusOK, []struct{}{})
		c.Abort()
	}
}

func (p *PrometheusDiscoveryApi) discovery(c *gin.Context, key string) {
	eResp, err := p.client.Get(c, prommetrics.BuildDiscoveryKey(key))
	if err != nil {
		// Log and respond with an error if preparation fails.
		apiresp.GinError(c, errs.WrapMsg(err, "etcd get err"))
		return
	}
	if len(eResp.Kvs) == 0 {
		c.JSON(http.StatusOK, []*prommetrics.Target{})
	}

	var (
		resp = &prommetrics.RespTarget{
			Targets: make([]string, 0, len(eResp.Kvs)),
		}
	)

	for i := range eResp.Kvs {
		var target prommetrics.Target
		err = json.Unmarshal(eResp.Kvs[i].Value, &target)
		if err != nil {
			log.ZError(c, "prometheus unmarshal err", errs.Wrap(err))
		}
		resp.Targets = append(resp.Targets, target.Target)
		if resp.Labels == nil {
			resp.Labels = target.Labels
		}
	}

	c.JSON(200, []*prommetrics.RespTarget{resp})
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
