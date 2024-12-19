package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/tools/apiresp"
	"gopkg.in/yaml.v3"
)

type ConfigManager struct {
	apiConfig *config.API
}

func NewConfigManager(api *config.API) *ConfigManager {
	return &ConfigManager{
		apiConfig: api,
	}
}

func (cm *ConfigManager) LoadApiConfig(c *gin.Context) {
	b, err := yaml.Marshal(cm.apiConfig)
	if err != nil {
		apiresp.GinError(c, err) // args option error
		return
	}
	c.Data(http.StatusOK, "text/yaml; charset=utf-8", b)
}
