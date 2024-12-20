package api

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/openimsdk/open-im-server/v3/pkg/apistruct"
	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/version"
	"github.com/openimsdk/tools/apiresp"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/runtimeenv"
)

type ConfigManager struct {
	imAdminUserID []string
	config        *config.AllConfig
}

func NewConfigManager(IMAdminUserID []string, cfg *config.AllConfig) *ConfigManager {
	return &ConfigManager{
		imAdminUserID: IMAdminUserID,
		config:        cfg,
	}
}

func (cm *ConfigManager) CheckAdmin(c *gin.Context) {
	if err := authverify.CheckAdmin(c, cm.imAdminUserID); err != nil {
		apiresp.GinError(c, err)
		c.Abort()
	}
}

func (cm *ConfigManager) GetConfig(c *gin.Context) {
	var req apistruct.GetConfigReq
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WithDetail(err.Error()).Wrap())
		return
	}
	conf := cm.config.Name2Config(req.ConfigName)
	if conf == nil {
		apiresp.GinError(c, errs.ErrArgs.WithDetail("config name not found").Wrap())
		return
	}
	b, err := json.Marshal(cm.config)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, string(b))
}

func (cm *ConfigManager) GetConfigList(c *gin.Context) {
	var resp apistruct.GetConfigListResp
	resp.ConfigNames = cm.config.GetConfigNames()
	resp.Environment = runtimeenv.PrintRuntimeEnvironment()
	resp.Version = version.Version

	apiresp.GinSuccess(c, resp)
}

func (cm *ConfigManager) SetConfig(c *gin.Context) {
	var req apistruct.SetConfigReq
	if err := c.BindJSON(&req); err != nil {
		apiresp.GinError(c, errs.ErrArgs.WithDetail(err.Error()).Wrap())
		return
	}
	b, err := json.Marshal(cm.config)
	if err != nil {
		apiresp.GinError(c, err) // args option error
		return
	}
	apiresp.GinSuccess(c, string(b))
}
