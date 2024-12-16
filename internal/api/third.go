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
	"context"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"google.golang.org/grpc"

	"github.com/gin-gonic/gin"
	"github.com/openimsdk/protocol/third"
	"github.com/openimsdk/tools/a2r"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/mcontext"
)

type ThirdApi struct {
	GrafanaUrl string
}

func NewThirdApi(grafanaUrl string) ThirdApi {
	return ThirdApi{GrafanaUrl: grafanaUrl}
}

func (o *ThirdApi) FcmUpdateToken(c *gin.Context) {
	a2r.CallV2(third.FcmUpdateTokenCaller.Invoke, c)
}

func (o *ThirdApi) SetAppBadge(c *gin.Context) {
	a2r.CallV2(third.SetAppBadgeCaller.Invoke, c)
}

// #################### s3 ####################

func setURLPrefixOption[A, B, C any](_ func(client C, ctx context.Context, req *A, options ...grpc.CallOption) (*B, error), fn func(*A) error) *a2r.Option[A, B] {
	return &a2r.Option[A, B]{
		BindAfter: fn,
	}
}

func setURLPrefix(c *gin.Context, urlPrefix *string) error {
	host := c.GetHeader("X-Request-Api")
	if host != "" {
		if strings.HasSuffix(host, "/") {
			*urlPrefix = host + "object/"
			return nil
		} else {
			*urlPrefix = host + "/object/"
			return nil
		}
	}
	u := url.URL{
		Scheme: "http",
		Host:   c.Request.Host,
		Path:   "/object/",
	}
	if c.Request.TLS != nil {
		u.Scheme = "https"
	}
	*urlPrefix = u.String()
	return nil
}

func (o *ThirdApi) PartLimit(c *gin.Context) {
	a2r.CallV2(third.PartLimitCaller.Invoke, c)
}

func (o *ThirdApi) PartSize(c *gin.Context) {
	a2r.CallV2(third.PartSizeCaller.Invoke, c)
}

func (o *ThirdApi) InitiateMultipartUpload(c *gin.Context) {
	opt := setURLPrefixOption(third.ThirdClient.InitiateMultipartUpload, func(req *third.InitiateMultipartUploadReq) error {
		return setURLPrefix(c, &req.UrlPrefix)
	})
	a2r.CallV2(third.InitiateMultipartUploadCaller.Invoke, c, opt)
}

func (o *ThirdApi) AuthSign(c *gin.Context) {
	a2r.CallV2(third.AuthSignCaller.Invoke, c)
}

func (o *ThirdApi) CompleteMultipartUpload(c *gin.Context) {
	opt := setURLPrefixOption(third.ThirdClient.CompleteMultipartUpload, func(req *third.CompleteMultipartUploadReq) error {
		return setURLPrefix(c, &req.UrlPrefix)
	})
	a2r.CallV2(third.CompleteMultipartUploadCaller.Invoke, c, opt)
}

func (o *ThirdApi) AccessURL(c *gin.Context) {
	a2r.CallV2(third.AccessURLCaller.Invoke, c)
}

func (o *ThirdApi) InitiateFormData(c *gin.Context) {
	a2r.CallV2(third.InitiateFormDataCaller.Invoke, c)
}

func (o *ThirdApi) CompleteFormData(c *gin.Context) {
	opt := setURLPrefixOption(third.ThirdClient.CompleteFormData, func(req *third.CompleteFormDataReq) error {
		return setURLPrefix(c, &req.UrlPrefix)
	})
	a2r.CallV2(third.CompleteFormDataCaller.Invoke, c, opt)
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
	resp, err := third.AccessURLCaller.Invoke(ctx, &third.AccessURLReq{Name: name, Query: query})
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
	a2r.CallV2(third.UploadLogsCaller.Invoke, c)
}

func (o *ThirdApi) DeleteLogs(c *gin.Context) {
	a2r.CallV2(third.DeleteLogsCaller.Invoke, c)
}

func (o *ThirdApi) SearchLogs(c *gin.Context) {
	a2r.CallV2(third.SearchLogsCaller.Invoke, c)
}

func (o *ThirdApi) GetPrometheus(c *gin.Context) {
	c.Redirect(http.StatusFound, o.GrafanaUrl)
}
