package api

import (
	"context"
	"google.golang.org/grpc"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/openimsdk/protocol/third"
	"github.com/openimsdk/tools/a2r"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/mcontext"
)

type ThirdApi struct {
	GrafanaUrl string
	Client     third.ThirdClient
}

func NewThirdApi(client third.ThirdClient, grafanaUrl string) ThirdApi {
	return ThirdApi{Client: client, GrafanaUrl: grafanaUrl}
}

func (o *ThirdApi) FcmUpdateToken(c *gin.Context) {
	a2r.Call(c, third.ThirdClient.FcmUpdateToken, o.Client)
}

func (o *ThirdApi) SetAppBadge(c *gin.Context) {
	a2r.Call(c, third.ThirdClient.SetAppBadge, o.Client)
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
	a2r.Call(c, third.ThirdClient.PartLimit, o.Client)
}

func (o *ThirdApi) PartSize(c *gin.Context) {
	a2r.Call(c, third.ThirdClient.PartSize, o.Client)
}

func (o *ThirdApi) InitiateMultipartUpload(c *gin.Context) {
	opt := setURLPrefixOption(third.ThirdClient.InitiateMultipartUpload, func(req *third.InitiateMultipartUploadReq) error {
		return setURLPrefix(c, &req.UrlPrefix)
	})
	a2r.Call(c, third.ThirdClient.InitiateMultipartUpload, o.Client, opt)
}

func (o *ThirdApi) AuthSign(c *gin.Context) {
	a2r.Call(c, third.ThirdClient.AuthSign, o.Client)
}

func (o *ThirdApi) CompleteMultipartUpload(c *gin.Context) {
	opt := setURLPrefixOption(third.ThirdClient.CompleteMultipartUpload, func(req *third.CompleteMultipartUploadReq) error {
		return setURLPrefix(c, &req.UrlPrefix)
	})
	a2r.Call(c, third.ThirdClient.CompleteMultipartUpload, o.Client, opt)
}

func (o *ThirdApi) AccessURL(c *gin.Context) {
	a2r.Call(c, third.ThirdClient.AccessURL, o.Client)
}

func (o *ThirdApi) InitiateFormData(c *gin.Context) {
	a2r.Call(c, third.ThirdClient.InitiateFormData, o.Client)
}

func (o *ThirdApi) CompleteFormData(c *gin.Context) {
	opt := setURLPrefixOption(third.ThirdClient.CompleteFormData, func(req *third.CompleteFormDataReq) error {
		return setURLPrefix(c, &req.UrlPrefix)
	})
	a2r.Call(c, third.ThirdClient.CompleteFormData, o.Client, opt)
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
	a2r.Call(c, third.ThirdClient.UploadLogs, o.Client)
}

func (o *ThirdApi) DeleteLogs(c *gin.Context) {
	a2r.Call(c, third.ThirdClient.DeleteLogs, o.Client)
}

func (o *ThirdApi) SearchLogs(c *gin.Context) {
	a2r.Call(c, third.ThirdClient.SearchLogs, o.Client)
}

func (o *ThirdApi) GetPrometheus(c *gin.Context) {
	c.Redirect(http.StatusFound, o.GrafanaUrl)
}
