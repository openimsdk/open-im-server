package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadLogConfig(t *testing.T) {
	var log Log
	err := LoadConfig("../../../config/log.yml", "OPENIM_LOG", &log)
	assert.Nil(t, err)
	assert.Equal(t, "/data/workspaces/open-im-server/_output/logs/", log.StorageLocation)
}

func TestLoadMinioConfig(t *testing.T) {
	var storageConfig Minio
	err := LoadConfig("../../../config/minio.yml", "OPENIM_MINIO", &storageConfig)
	assert.Nil(t, err)
	assert.Equal(t, "openim", storageConfig.Bucket)
}

func TestLoadWebhooksConfig(t *testing.T) {
	var webhooks Webhooks
	err := LoadConfig("../../../config/webhooks.yml", "OPENIM_WEBHOOKS", &webhooks)
	assert.Nil(t, err)
	assert.Equal(t, 5, webhooks.AddBlackBefore.Timeout)
}
