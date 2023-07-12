package api

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/a2r"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/third"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"strconv"
)

type ThirdApi rpcclient.Third

func NewThirdApi(discov discoveryregistry.SvcDiscoveryRegistry) ThirdApi {
	return ThirdApi(*rpcclient.NewThird(discov))
}

func (o *ThirdApi) FcmUpdateToken(c *gin.Context) {
	a2r.Call(third.ThirdClient.FcmUpdateToken, o.Client, c)
}

func (o *ThirdApi) SetAppBadge(c *gin.Context) {
	a2r.Call(third.ThirdClient.SetAppBadge, o.Client, c)
}

// #################### s3 ####################

func (o *ThirdApi) PartLimit(c *gin.Context) {
	a2r.Call(third.ThirdClient.PartLimit, o.Client, c)
}

func (o *ThirdApi) PartSize(c *gin.Context) {
	a2r.Call(third.ThirdClient.PartSize, o.Client, c)
}

func (o *ThirdApi) InitiateMultipartUpload(c *gin.Context) {
	a2r.Call(third.ThirdClient.InitiateMultipartUpload, o.Client, c)
}

func (o *ThirdApi) AuthSign(c *gin.Context) {
	a2r.Call(third.ThirdClient.AuthSign, o.Client, c)
}

func (o *ThirdApi) CompleteMultipartUpload(c *gin.Context) {
	a2r.Call(third.ThirdClient.CompleteMultipartUpload, o.Client, c)
}

func (o *ThirdApi) AccessURL(c *gin.Context) {
	a2r.Call(third.ThirdClient.AccessURL, o.Client, c)
}

func (o *ThirdApi) ObjectRedirect(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		c.String(http.StatusBadRequest, "name is empty")
		return
	}
	if name[0] == '/' {
		name = name[1:]
	}
	operationID := c.Query("operationID")
	if operationID == "" {
		operationID = strconv.Itoa(rand.Int())
	}
	ctx := mcontext.SetOperationID(c, operationID)
	resp, err := o.Client.AccessURL(ctx, &third.AccessURLReq{Name: name})
	if err != nil {
		if errs.ErrArgs.Is(err) {
			c.String(http.StatusBadRequest, err.Error())
			return
		}
		if errs.ErrRecordNotFound.Is(err) {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, resp.Url)
}
