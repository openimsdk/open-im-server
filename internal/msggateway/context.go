package msggateway

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"net/http"
	"strconv"
	"time"
)

type UserConnContext struct {
	RespWriter http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	RemoteAddr string
	ConnID     string
}

func (c *UserConnContext) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c *UserConnContext) Done() <-chan struct{} {
	return nil
}

func (c *UserConnContext) Err() error {
	return nil
}

func (c *UserConnContext) Value(key any) any {
	switch key {
	case constant.OpUserID:
		c.GetUserID()
	case constant.OperationID:
		c.GetOperationID()
	case constant.ConnID:
		c.GetConnID()
	case constant.OpUserPlatform:
		constant.PlatformIDToName(utils.StringToInt(c.GetPlatformID()))
	default:
		return ""
	}
	return ""
}

func newContext(respWriter http.ResponseWriter, req *http.Request) *UserConnContext {
	return &UserConnContext{
		RespWriter: respWriter,
		Req:        req,
		Path:       req.URL.Path,
		Method:     req.Method,
		RemoteAddr: req.RemoteAddr,
		ConnID:     utils.Md5(req.RemoteAddr + "_" + strconv.Itoa(int(utils.GetCurrentTimestampByMill()))),
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
	return c.ConnID
}
func (c *UserConnContext) GetUserID() string {
	return c.Req.URL.Query().Get(WsUserID)
}
func (c *UserConnContext) GetPlatformID() string {
	return c.Req.URL.Query().Get(PlatformID)
}
func (c *UserConnContext) GetOperationID() string {
	return c.Req.URL.Query().Get(OperationID)
}
