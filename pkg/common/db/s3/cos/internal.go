package cos

import (
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
	_ "unsafe"
)

//go:linkname newRequest github.com/tencentyun/cos-go-sdk-v5.(*Client).newRequest
func newRequest(c *cos.Client, ctx context.Context, baseURL *url.URL, uri, method string, body interface{}, optQuery interface{}, optHeader interface{}) (req *http.Request, err error)
