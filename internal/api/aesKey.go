package api

import (
	aesKey "github.com/OpenIMSDK/protocol/aeskey"
	"github.com/OpenIMSDK/tools/a2r"
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
)

type AesKeyApi rpcclient.AesKey

func NewAesKeyApi(client rpcclient.AesKey) AesKeyApi {
	return AesKeyApi(client)
}

func (a *AesKeyApi) GetKey(c *gin.Context) {
	a2r.Call(aesKey.AesKeyClient.AcquireAesKey, a.Client, c)
}

func (a *AesKeyApi) GetKeys(c *gin.Context) {
	a2r.Call(aesKey.AesKeyClient.AcquireAesKeys, a.Client, c)
}
