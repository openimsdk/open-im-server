package api

import (
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/protocol/crypto"
	"github.com/openimsdk/tools/a2r"
)

type CryptoApi struct {
	Client crypto.CryptoServiceClient
}

func NewCryptoApi(client crypto.CryptoServiceClient) CryptoApi {
	return CryptoApi{Client: client}
}

func (o *CryptoApi) RegisterDevice(c *gin.Context) {
	a2r.Call(c, crypto.CryptoServiceClient.RegisterDevice, o.Client)
}

func (o *CryptoApi) GetDevices(c *gin.Context) {
	a2r.Call(c, crypto.CryptoServiceClient.GetDevices, o.Client)
}

func (o *CryptoApi) RevokeDevice(c *gin.Context) {
	a2r.Call(c, crypto.CryptoServiceClient.RevokeDevice, o.Client)
}

func (o *CryptoApi) GetVirgilJWT(c *gin.Context) {
	a2r.Call(c, crypto.CryptoServiceClient.GetVirgilJWT, o.Client)
}

func (o *CryptoApi) GetGroupKeyVersion(c *gin.Context) {
	a2r.Call(c, crypto.CryptoServiceClient.GetGroupKeyVersion, o.Client)
}

func (o *CryptoApi) GetGroupKeyEvents(c *gin.Context) {
	a2r.Call(c, crypto.CryptoServiceClient.GetGroupKeyEvents, o.Client)
}

func (o *CryptoApi) SecurityPrecheck(c *gin.Context) {
	a2r.Call(c, crypto.CryptoServiceClient.SecurityPrecheck, o.Client)
}

func (o *CryptoApi) IntegrityReport(c *gin.Context) {
	a2r.Call(c, crypto.CryptoServiceClient.IntegrityReport, o.Client)
}
