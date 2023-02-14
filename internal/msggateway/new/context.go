package new

import "net/http"

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
func (c *UserConnContext) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}
func (c *UserConnContext) GetHeader(key string) string {
	return c.Req.Header.Get(key)
}
