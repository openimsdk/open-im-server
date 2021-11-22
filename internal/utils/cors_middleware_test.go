package utils

import (
	"Open_IM/pkg/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func performRequest(r http.Handler, method, origin string) *httptest.ResponseRecorder {
	return performRequestWithHeaders(r, method, origin, http.Header{})
}

func performRequestWithHeaders(r http.Handler, method, origin string, header http.Header) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, "/", nil)
	// From go/net/http/request.go:
	// For incoming requests, the Host header is promoted to the
	// Request.Host field and removed from the Header map.
	req.Host = header.Get("Host")
	header.Del("Host")
	if len(origin) > 0 {
		header.Set("Origin", origin)
	}
	req.Header = header
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func newTestRouter() *gin.Engine {
	router := gin.New()
	router.Use(utils.CorsHandler())
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "get")
	})
	router.POST("/", func(c *gin.Context) {
		c.String(http.StatusOK, "post")
	})
	router.PATCH("/", func(c *gin.Context) {
		c.String(http.StatusOK, "patch")
	})

	return router
}

func Test_CorsHandler(t *testing.T) {
	router := newTestRouter()
	// no CORS request, origin == ""
	w := performRequest(router, "GET", "")
	assert.Equal(t, "get", w.Body.String())
	assert.Equal(t, w.Header().Get("Access-Control-Allow-Origin"), "*")
	assert.Equal(t, w.Header().Get("Access-Control-Allow-Methods"), "*")
	assert.Equal(t, w.Header().Get("Access-Control-Allow-Headers"), "*")
	assert.Equal(t, w.Header().Get("Access-Control-Expose-Headers"), "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers,Cache-Control,Content-Language,Content-Type,Expires,Last-Modified,Pragma,FooBar")
	assert.Equal(t, w.Header().Get("Access-Control-Max-Age"), "172800")
	assert.Equal(t, w.Header().Get("Access-Control-Allow-Credentials"), "false")
	assert.Equal(t, w.Header().Get("content-type"), "application/json")

	w = performRequest(router, "OPTIONS", "")
	assert.Equal(t, w.Body.String(), "\"Options Request!\"")
}
