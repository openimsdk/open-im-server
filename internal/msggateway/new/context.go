package new

import (
	"OpenIM/pkg/utils"
	"net/http"
	"strconv"
)

type UserConnContext struct {
	RespWriter http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	RemoteAddr string
}

func newContext(respWriter http.ResponseWriter, req *http.Request) *UserConnContext {
	return &UserConnContext{
		RespWriter: respWriter,
		Req:        req,
		Path:       req.URL.Path,
		Method:     req.Method,
		RemoteAddr: req.RemoteAddr,
	}
}
func (c *UserConnContext) Query(key string) (string, bool) {
	var value string
	if value = c.Req.URL.Query().Get(key); value == "" {
		return value, false
	}
	return value, true
}
func (c *UserConnContext) GetHeader(key string) (string, bool) {
	var value string
	if value = c.Req.Header.Get(key); value == "" {
		return value, false
	}
	return value, true
}
func (c *UserConnContext) SetHeader(key, value string) {
	c.RespWriter.Header().Set(key, value)
}
func (c *UserConnContext) ErrReturn(error string, code int) {
	http.Error(c.RespWriter, error, code)
}
func (c *UserConnContext) GetConnID() string {
	return c.RemoteAddr + "_" + strconv.Itoa(int(utils.GetCurrentTimestampByMill()))
}
func (c *UserConnContext) GetUserID() string {
	return c.Req.URL.Query().Get(WS_USERID)
}
func (c *UserConnContext) GetPlatformID() string {
	return c.Req.URL.Query().Get(PLATFORM_ID)
}
