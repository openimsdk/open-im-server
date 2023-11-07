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

	config2 "github.com/openimsdk/open-im-server/v3/pkg/common/config"

	"github.com/gin-gonic/gin"

	"github.com/OpenIMSDK/protocol/third"
	"github.com/OpenIMSDK/tools/a2r"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/mcontext"

	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
)

type ThirdApi rpcclient.Third

func NewThirdApi(client rpcclient.Third) ThirdApi {
	return ThirdApi(client)
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
	query := make(map[string]string)
	for key, values := range c.Request.URL.Query() {
		if len(values) == 0 {
			continue
		}
		query[key] = values[0]
	}
	resp, err := o.Client.AccessURL(ctx, &third.AccessURLReq{Name: name, Query: query})
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
	c.Redirect(http.StatusFound, resp.Url)
}

// #################### logs ####################.
func (o *ThirdApi) UploadLogs(c *gin.Context) {
	a2r.Call(third.ThirdClient.UploadLogs, o.Client, c)
}

func (o *ThirdApi) DeleteLogs(c *gin.Context) {
	a2r.Call(third.ThirdClient.DeleteLogs, o.Client, c)
}

func (o *ThirdApi) SearchLogs(c *gin.Context) {
	a2r.Call(third.ThirdClient.SearchLogs, o.Client, c)
}

func GetPrometheus(c *gin.Context) {
	c.Redirect(http.StatusFound, config2.Config.Prometheus.PrometheusUrl)
}
