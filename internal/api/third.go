// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"math/rand"
	"net/http"
	"strconv"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/a2r"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/third"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
	"github.com/gin-gonic/gin"
)

type ThirdApi rpcclient.Third

func NewThirdApi(discov discoveryregistry.SvcDiscoveryRegistry) ThirdApi {
	return ThirdApi(*rpcclient.NewThird(discov))
}

func (o *ThirdApi) ApplyPut(c *gin.Context) {
	a2r.Call(third.ThirdClient.ApplyPut, o.Client, c)
}

func (o *ThirdApi) GetPut(c *gin.Context) {
	a2r.Call(third.ThirdClient.GetPut, o.Client, c)
}

func (o *ThirdApi) ConfirmPut(c *gin.Context) {
	a2r.Call(third.ThirdClient.ConfirmPut, o.Client, c)
}

func (o *ThirdApi) GetHash(c *gin.Context) {
	a2r.Call(third.ThirdClient.GetHashInfo, o.Client, c)
}

func (o *ThirdApi) FcmUpdateToken(c *gin.Context) {
	a2r.Call(third.ThirdClient.FcmUpdateToken, o.Client, c)
}

func (o *ThirdApi) SetAppBadge(c *gin.Context) {
	a2r.Call(third.ThirdClient.SetAppBadge, o.Client, c)
}

func (o *ThirdApi) GetURL(c *gin.Context) {
	if c.Request.Method == http.MethodPost {
		a2r.Call(third.ThirdClient.GetUrl, o.Client, c)
		return
	}
	name := c.Query("name")
	if name == "" {
		c.String(http.StatusBadRequest, "name is empty")
		return
	}
	operationID := c.Query("operationID")
	if operationID == "" {
		operationID = "auto_" + strconv.Itoa(rand.Int())
	}
	expires, _ := strconv.ParseInt(c.Query("expires"), 10, 64)
	if expires <= 0 {
		expires = 3600 * 1000
	}
	attachment, _ := strconv.ParseBool(c.Query("attachment"))
	c.Set(constant.OperationID, operationID)
	resp, err := o.Client.GetUrl(mcontext.SetOperationID(c, operationID), &third.GetUrlReq{Name: name, Expires: expires, Attachment: attachment})
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
