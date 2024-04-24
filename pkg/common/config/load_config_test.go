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
	err := LoadConfig("../../../config/openim-rpc-user.yml", "IMENV_OPENIM_RPC_USER", &user)
	assert.Nil(t, err)
	//export IMENV_OPENIM_RPC_USER_RPC_LISTENIP="0.0.0.0"
	assert.Equal(t, "0.0.0.0", user.RPC.ListenIP)
	//export IMENV_OPENIM_RPC_USER_RPC_PORTS="10110,10111,10112"
	assert.Equal(t, []int{10110, 10111, 10112}, user.RPC.Ports)
}
