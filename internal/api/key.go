package api

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/a2r"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/key"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
	"github.com/gin-gonic/gin"
)

type KeyApi rpcclient.Key

func NewKeyApi(client rpcclient.Key) KeyApi {
	return KeyApi(client)
}

func (o *KeyApi) GetKey(c *gin.Context) {
	a2r.Call(key.KeyClient.GetKey, o.Client, c)
}
