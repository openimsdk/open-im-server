package api

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/tools/apiresp"
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
	b, err := json.Marshal(cm.apiConfig)
	if err != nil {
		apiresp.GinError(c, err) // args option error
		return
	}
	apiresp.GinSuccess(c, string(b))
}
