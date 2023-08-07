package api

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/aes_key"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
	"github.com/OpenIMSDK/tools/a2r"
	"github.com/gin-gonic/gin"
)

type AesKeyApi rpcclient.AesKey

func NewAesKeyApi(client rpcclient.AesKey) AesKeyApi {
	return AesKeyApi(client)
}

func (a *AesKeyApi) GetKey(c *gin.Context) {
	a2r.Call(aes_key.AesKeyClient.GetAesKey, a.Client, c)
}
func (a *AesKeyApi) GetAllKey(c *gin.Context) {
	a2r.Call(aes_key.AesKeyClient.GetAllAesKey, a.Client, c)
}
