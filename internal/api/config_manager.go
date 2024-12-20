package api

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/tools/apiresp"
)

type ConfigManager struct {
	config *config.AllConfig
}

func NewConfigManager(cfg *config.AllConfig) *ConfigManager {
	return &ConfigManager{
		config: cfg,
	}
}

func (cm *ConfigManager) GetConfig(c *gin.Context) {
	b, err := json.Marshal(cm.config)
	if err != nil {
		apiresp.GinError(c, err) // args option error
		return
	}
	apiresp.GinSuccess(c, string(b))
}
