package cos

import (
	"context"
	"net/http"
	"net/url"
	_ "unsafe"

	"github.com/tencentyun/cos-go-sdk-v5"
)

//go:linkname newRequest github.com/tencentyun/cos-go-sdk-v5.(*Client).newRequest
func newRequest(c *cos.Client, ctx context.Context, baseURL *url.URL, uri, method string, body interface{}, optQuery interface{}, optHeader interface{}) (req *http.Request, err error)
