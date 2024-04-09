package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadLogConfig(t *testing.T) {
	var log Log
	err := LoadConfig("../../../config/log.yml", "IMENV_LOG", &log)
	assert.Nil(t, err)
	assert.Equal(t, "../../../../logs/", log.StorageLocation)
}

func TestLoadMinioConfig(t *testing.T) {
	var storageConfig Minio
	err := LoadConfig("../../../config/minio.yml", "IMENV_MINIO", &storageConfig)
	assert.Nil(t, err)
	assert.Equal(t, "openim", storageConfig.Bucket)
}

func TestLoadWebhooksConfig(t *testing.T) {
	var webhooks Webhooks
	err := LoadConfig("../../../config/webhooks.yml", "IMENV_WEBHOOKS", &webhooks)
	assert.Nil(t, err)
	assert.Equal(t, 5, webhooks.BeforeAddBlack.Timeout)

}

func TestLoadOpenIMRpcUserConfig(t *testing.T) {
	var user User
	err := LoadConfig("../../../config/openim-rpc-user.yml", "IMENV_OPENIM-RPC-USER", &user)
	assert.Nil(t, err)
	assert.Equal(t, "0.0.0.0", user.RPC.ListenIP)
	assert.Equal(t, []int{10110}, user.RPC.Ports)
}
